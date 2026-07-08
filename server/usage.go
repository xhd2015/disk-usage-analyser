package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"disk-usage-analyser/usagescan"
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
			fmt.Fprintf(w, "event: server_error\ndata: {\"error\": \"Internal Server Error: %v\"}\n\n", r)
		}
	}()

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

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	if err := sendEvent(w, "path", map[string]string{"path": dirPath}); err != nil {
		return
	}
	flusher.Flush()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	result, err := usagescan.ScanFlat(dirPath, usagescan.Options{
		Context: ctx,
		OnItem: func(item usagescan.Item) {
			status := "done"
			if item.IsDir {
				status = "pending"
			}
			if sendErr := sendEvent(w, "item", FileInfo{
				Name:   item.Name,
				Size:   item.Size,
				IsDir:  item.IsDir,
				Status: status,
			}); sendErr != nil {
				log.Printf("Client disconnected, stopping scan")
				cancel()
			} else {
				flusher.Flush()
			}
		},
	})
	if err != nil {
		log.Printf("Error scanning directory %s: %v", dirPath, err)
		sendEvent(w, "server_error", map[string]string{"error": err.Error()})
		return
	}

	for _, item := range result.Items {
		if !item.IsDir {
			continue
		}
		if err := sendEvent(w, "item", FileInfo{
			Name:   item.Name,
			Size:   item.Size,
			IsDir:  true,
			Status: "done",
		}); err != nil {
			log.Printf("Client disconnected, stopping scan")
			return
		}
		flusher.Flush()
	}

	sendEvent(w, "done", nil)
	flusher.Flush()
}

func handleMoveToTrash(w http.ResponseWriter, r *http.Request) {
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
		escapedPath := strings.ReplaceAll(path, "\"", "\\\"")
		script := fmt.Sprintf(`tell application "Finder" to move POSIX file "%s" to trash`, escapedPath)
		cmd := exec.Command("osascript", "-e", script)
		out, err := cmd.CombinedOutput()
		if err != nil {
			http.Error(w, fmt.Sprintf("Trash failed: %v, %s", err, string(out)), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Move to trash not supported on this OS", http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleRefresh(w http.ResponseWriter, r *http.Request) {
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