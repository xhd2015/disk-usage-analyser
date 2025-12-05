package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// Ensure absolute path
	if !filepath.IsAbs(dirPath) {
		absPath, err := filepath.Abs(dirPath)
		if err != nil {
			http.Error(w, "Invalid path: "+err.Error(), http.StatusBadRequest)
			return
		}
		dirPath = absPath
	}

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

	// Start workers for directories
	for _, dir := range subDirs {
		wg.Add(1)
		go func(d fs.DirEntry) {
			defer wg.Done()
			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }()

			onProgress := func(currentSize int64) {
				resultChan <- FileInfo{
					Name:   d.Name(),
					Size:   currentSize,
					IsDir:  true,
					Status: "pending",
				}
			}

			size := getDirSizeWithProgress(filepath.Join(dirPath, d.Name()), onProgress)
			resultChan <- FileInfo{
				Name:   d.Name(),
				Size:   size,
				IsDir:  true,
				Status: "done",
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
		sendEvent(w, "item", item)
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

func sendEvent(w http.ResponseWriter, event string, data interface{}) {
	jsonData, _ := json.Marshal(data)
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, jsonData)
}

func getDirSizeWithProgress(path string, onProgress func(int64)) int64 {
	var size int64
	var count int
	lastUpdate := time.Now()

	filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
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
