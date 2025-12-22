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

// Semaphore to limit concurrent ReadDir operations
var scanSem = make(chan struct{}, 20)

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
	if err := sendEvent(w, "path", map[string]string{"path": dirPath}); err != nil {
		return
	}
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
	// Limit concurrency for top level response handling
	// Note: scanDirRecursive now handles its own concurrency,
	// but we still want to limit how many `getDirSizeWithCache` we invoke concurrently from here
	// to avoid overwhelming the system if a folder has 10k subfolders.
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

			fullPath := filepath.Join(dirPath, d.Name())

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

			// Use the smart cache-aware scanner
			size := getDirSizeWithCache(ctx, fullPath, onProgress)

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

func handleRefresh(w http.ResponseWriter, r *http.Request) {
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

	// Ensure absolute path
	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			http.Error(w, "Invalid path: "+err.Error(), http.StatusBadRequest)
			return
		}
		path = absPath
	}

	log.Printf("Invalidating cache for path: %s", path)
	GlobalCache.Invalidate(path)

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

// getDirSizeWithCache checks the cache first. If scanning is needed, it performs it.
// If scanning is already in progress (by another request/worker), it subscribes to it.
func getDirSizeWithCache(ctx context.Context, path string, onProgress func(int64)) int64 {
	entry, exists := GlobalCache.GetOrCreateEntry(path)

	if !exists {
		// We own it. Start scanning in background.
		go scanDirRecursive(ctx, path, entry)
	}

	// Subscribe to progress updates
	unsubscribe := entry.Subscribe(func(s int64) {
		onProgress(s)
	})
	defer unsubscribe()

	// Wait until done or context cancelled
	select {
	case <-entry.doneCh:
		return entry.Size
	case <-ctx.Done():
		return entry.Size
	}
}

// scanDirRecursive implements a recursive scan to correctly handle cache population
// It updates the entry in real-time as subdirectories are scanned.
func scanDirRecursive(ctx context.Context, dirPath string, entry *CacheEntry) {
	defer entry.MarkDone()

	// Acquire semaphore for IO (ReadDir)
	select {
	case scanSem <- struct{}{}:
	case <-ctx.Done():
		return
	}

	entries, err := os.ReadDir(dirPath)
	// Release semaphore
	<-scanSem

	if err != nil {
		log.Printf("Error reading %s: %v", dirPath, err)
		return
	}

	var (
		mu          sync.Mutex
		filesSize   int64
		subDirSizes = make(map[string]int64)
		dirty       bool
		wg          sync.WaitGroup
	)

	// Ticker to push updates to entry
	ticker := time.NewTicker(200 * time.Millisecond)
	doneCh := make(chan struct{})

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-doneCh:
				return
			case <-ticker.C:
				mu.Lock()
				if dirty {
					total := filesSize
					for _, s := range subDirSizes {
						total += s
					}
					entry.UpdateSize(total)
					dirty = false
				}
				mu.Unlock()
			}
		}
	}()

	updateLocal := func(name string, size int64) {
		mu.Lock()
		subDirSizes[name] = size
		dirty = true
		mu.Unlock()
	}

	for _, e := range entries {
		if ctx.Err() != nil {
			break
		}

		if !e.IsDir() {
			info, err := e.Info()
			if err == nil {
				mu.Lock()
				filesSize += info.Size()
				dirty = true
				mu.Unlock()
			}
		} else {
			subPath := filepath.Join(dirPath, e.Name())
			subName := e.Name()

			wg.Add(1)

			// Handle subdirectories
			subEntry, exists := GlobalCache.GetOrCreateEntry(subPath)

			if !exists {
				// We start it
				go scanDirRecursive(ctx, subPath, subEntry)
			}

			// Subscribe to changes
			unsub := subEntry.Subscribe(func(size int64) {
				updateLocal(subName, size)
			})

			// Wait for done to decrement WG
			go func() {
				defer wg.Done()
				defer unsub() // Unsubscribe when done waiting
				subEntry.Wait()
			}()
		}
	}

	// Wait for all children to complete
	wg.Wait()
	close(doneCh)

	// Final update
	mu.Lock()
	total := filesSize
	for _, s := range subDirSizes {
		total += s
	}
	entry.UpdateSize(total)
	mu.Unlock()
}
