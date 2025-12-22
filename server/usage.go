package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var InitialDir string

type FileInfo struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	IsDir  bool   `json:"isDir"`
	Status string `json:"status"` // "pending", "done"
}

func handleUsage(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in handleUsage: %v\nStack: %s", r, debug.Stack())
			// Try to send error event if possible
			fmt.Fprintf(w, "event: server_error\ndata: {\"error\": \"Internal Server Error: %v\"}\n\n", r)
		}
	}()

	// CORS for dev
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		return
	}

	dirPath := r.URL.Query().Get("path")
	if dirPath == "" {
		if InitialDir != "" {
			dirPath = InitialDir
		} else {
			var err error
			dirPath, err = os.Getwd()
			if err != nil {
				log.Printf("Error getting current working directory: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// Ensure absolute path
	if !filepath.IsAbs(dirPath) {
		absPath, err := filepath.Abs(dirPath)
		if err != nil {
			log.Printf("Error resolving absolute path for %s: %v", dirPath, err)
			http.Error(w, "Invalid path: "+err.Error(), http.StatusBadRequest)
			return
		}
		dirPath = absPath
	}

	log.Printf("Starting usage scan for path: %s", dirPath)

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Send path info event
	sendEvent(w, "path", map[string]string{"path": dirPath})
	flusher.Flush()

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Printf("Error reading directory %s: %v", dirPath, err)
		sendEvent(w, "server_error", map[string]string{"error": err.Error()})
		return
	}

	// Identify subdirectories and files
	var subDirs []fs.DirEntry
	var files []fs.DirEntry

	for _, entry := range entries {
		if entry.IsDir() {
			subDirs = append(subDirs, entry)
		} else {
			files = append(files, entry)
		}
	}

	// Send all files immediately
	for _, entry := range files {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		sendEvent(w, "item", FileInfo{
			Name:   entry.Name(),
			Size:   info.Size(),
			IsDir:  false,
			Status: "done",
		})
	}
	// Send all directories immediately with pending status
	for _, entry := range subDirs {
		sendEvent(w, "item", FileInfo{
			Name:   entry.Name(),
			Size:   0,
			IsDir:  true,
			Status: "pending",
		})
	}
	flusher.Flush()

	// Channel to collect results from workers
	resultChan := make(chan FileInfo)
	var wg sync.WaitGroup
	// Limit concurrency
	sem := make(chan struct{}, 20)

	// Create cancellable context
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Start workers for directories
	for _, dir := range subDirs {
		wg.Add(1)
		go func(d fs.DirEntry) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic in worker for %s: %v", d.Name(), r)
				}
			}()

			// Acquire semaphore
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()

			onProgress := func(currentSize int64) {
				select {
				case resultChan <- FileInfo{
					Name:   d.Name(),
					Size:   currentSize,
					IsDir:  true,
					Status: "pending",
				}:
				case <-ctx.Done():
				}
			}

			size := getDirSizeWithProgress(ctx, filepath.Join(dirPath, d.Name()), onProgress)

			select {
			case resultChan <- FileInfo{
				Name:   d.Name(),
				Size:   size,
				IsDir:  true,
				Status: "done",
			}:
			case <-ctx.Done():
			}
		}(dir)
	}

	// Closer routine
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Stream results as they arrive
	for item := range resultChan {
		if err := sendEvent(w, "item", item); err != nil {
			log.Printf("Client disconnected, stopping scan")
			return
		}
		flusher.Flush()
	}

	sendEvent(w, "done", nil)
	flusher.Flush()
}

func handleMoveToTrash(w http.ResponseWriter, r *http.Request) {
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path required", http.StatusBadRequest)
		return
	}

	if runtime.GOOS == "darwin" {
		// Use AppleScript to move to trash via Finder
		// Escape double quotes in path
		escapedPath := strings.ReplaceAll(path, "\"", "\\\"")
		script := fmt.Sprintf(`tell application "Finder" to move POSIX file "%s" to trash`, escapedPath)
		cmd := exec.Command("osascript", "-e", script)
		out, err := cmd.CombinedOutput()
		if err != nil {
			http.Error(w, fmt.Sprintf("Trash failed: %v, %s", err, string(out)), http.StatusInternalServerError)
			return
		}
	} else {
		// Fallback or error?
		// User requested safer delete.
		http.Error(w, "Move to trash not supported on this OS", http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func sendEvent(w http.ResponseWriter, event string, data interface{}) error {
	jsonData, _ := json.Marshal(data)
	_, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, jsonData)
	if err != nil {
		log.Printf("Error sending event %s: %v", event, err)
		return err
	}
	return nil
}

func getDirSizeWithProgress(ctx context.Context, path string, onProgress func(int64)) int64 {
	var size int64
	var count int
	lastUpdate := time.Now()

	filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			// Don't log every permission denied, might be too noisy, but for now user requested logs
			// Filter out common harmless errors if needed, but "incomplete chunked encoding" suggests something worse
			log.Printf("Error walking %s: %v", p, err)
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err == nil {
				size += info.Size()
				count++
			}
			if count%1000 == 0 {
				if time.Since(lastUpdate) > 200*time.Millisecond {
					onProgress(size)
					lastUpdate = time.Now()
				}
			}
		}
		return nil
	})
	return size
}
