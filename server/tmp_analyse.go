package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type TmpLocation struct {
	Path       string `json:"path"`
	Label      string `json:"label"`
	Category   string `json:"category"`
	Size       int64  `json:"size"`
	FileCount  int64  `json:"fileCount"`
	RebootSafe bool   `json:"rebootSafe"`
}

type TmpSummary struct {
	Locations       []TmpLocation `json:"locations"`
	TotalSize       int64         `json:"totalSize"`
	ReclaimableSize int64         `json:"reclaimableSize"`
}

func DiscoverLocations(homeDir string) []TmpLocation {
	return []TmpLocation{
		{Path: filepath.Join(homeDir, ".Trash"), Label: "User Trash", Category: "trash", RebootSafe: true},
		{Path: filepath.Join(homeDir, "Library", "Caches"), Label: "User Caches", Category: "cache", RebootSafe: true},
		{Path: filepath.Join(homeDir, "Library", "Logs"), Label: "User Logs", Category: "log", RebootSafe: true},
		{Path: os.TempDir(), Label: "System Temp", Category: "temp", RebootSafe: false},
		{Path: "/tmp", Label: "System Tmp", Category: "temp", RebootSafe: false},
	}
}

func CalculateSize(fsys fs.FS, root string) (int64, int64, error) {
	var size int64
	var count int64
	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		size += info.Size()
		count++
		return nil
	})
	return size, count, err
}

func ScanWithProgress(fsys fs.FS, root string, onProgress func(size int64, count int64)) (int64, int64, error) {
	var size int64
	var count int64
	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		size += info.Size()
		count++
		onProgress(size, count)
		return nil
	})
	return size, count, err
}

func BuildSummary(locations []TmpLocation) TmpSummary {
	var total, reclaimable int64
	for _, loc := range locations {
		total += loc.Size
		if loc.RebootSafe {
			reclaimable += loc.Size
		}
	}
	return TmpSummary{
		Locations:       locations,
		TotalSize:       total,
		ReclaimableSize: reclaimable,
	}
}

func BuildProgressPayload(label string, curSize, curCount int64, accumulatedSize, accumulatedReclaimable int64, rebootSafe bool) map[string]interface{} {
	totalSize := accumulatedSize + curSize
	reclaimableSize := accumulatedReclaimable
	if rebootSafe {
		reclaimableSize += curSize
	}
	return map[string]interface{}{
		"label":           label,
		"size":            curSize,
		"fileCount":       curCount,
		"totalSize":       totalSize,
		"reclaimableSize": reclaimableSize,
	}
}

func HandleTmpAnalyse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		sendSSEEvent(w, "server_error", map[string]string{"error": err.Error()})
		return
	}

	locations := DiscoverLocations(homeDir)

	var accumulatedSize, accumulatedReclaimable int64

	for i := range locations {
		fsys := os.DirFS(locations[i].Path)
		label := locations[i].Label
		rebootSafe := locations[i].RebootSafe

		size, count, err := ScanWithProgress(fsys, ".", func(curSize int64, curCount int64) {
			progress := BuildProgressPayload(label, curSize, curCount, accumulatedSize, accumulatedReclaimable, rebootSafe)
			sendSSEEvent(w, "progress", progress)
			flusher.Flush()
		})
		if err != nil {
			log.Printf("Error scanning %s: %v (accumulated: %d bytes, %d files)", locations[i].Path, err, size, count)
		}
		locations[i].Size = size
		locations[i].FileCount = count

		accumulatedSize += size
		if rebootSafe {
			accumulatedReclaimable += size
		}

		if err := sendSSEEvent(w, "location", locations[i]); err != nil {
			return
		}
		flusher.Flush()
	}

	summary := BuildSummary(locations)
	if err := sendSSEEvent(w, "summary", summary); err != nil {
		return
	}
	flusher.Flush()

	sendSSEEvent(w, "done", map[string]string{"status": "complete"})
	flusher.Flush()
}

func sendSSEEvent(w http.ResponseWriter, event string, data interface{}) error {
	jsonData, _ := json.Marshal(data)
	_, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, jsonData)
	if err != nil {
		log.Printf("Error sending SSE event %s: %v", event, err)
		return err
	}
	return nil
}
