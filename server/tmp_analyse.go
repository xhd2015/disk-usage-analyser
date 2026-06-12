package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

type TmpLocation struct {
	Path       string   `json:"path"`
	Label      string   `json:"label"`
	Category   string   `json:"category"`
	Size       int64    `json:"size"`
	FileCount  int64    `json:"fileCount"`
	RebootSafe bool     `json:"rebootSafe"`
	Detected   bool     `json:"detected"`
	ExtraPaths []string `json:"extraPaths,omitempty"`
	ExtraSizes []int64  `json:"extraSizes,omitempty"`
	ExtraCounts []int64 `json:"extraCounts,omitempty"`
}

type TmpSummary struct {
	Locations       []TmpLocation `json:"locations"`
	TotalSize       int64         `json:"totalSize"`
	ReclaimableSize int64         `json:"reclaimableSize"`
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func DiscoverLocations(homeDir string) []TmpLocation {
	coreLocations := []TmpLocation{
		{Path: filepath.Join(homeDir, ".Trash"), Label: "User Trash", Category: "trash", RebootSafe: true, Detected: true},
		{Path: filepath.Join(homeDir, "Library", "Caches"), Label: "User Caches", Category: "cache", RebootSafe: true, Detected: true},
		{Path: filepath.Join(homeDir, "Library", "Logs"), Label: "User Logs", Category: "log", RebootSafe: true, Detected: true},
		{Path: os.TempDir(), Label: "System Temp", Category: "temp", RebootSafe: false, Detected: true},
		{Path: "/tmp", Label: "System Tmp", Category: "temp", RebootSafe: false, Detected: true},
	}

	softwareLocations := []TmpLocation{
		{Path: filepath.Join(homeDir, "go", "pkg", "mod"), Label: "Go", Category: "go", RebootSafe: true, ExtraPaths: []string{filepath.Join(homeDir, "Library", "Caches", "go-build")}},
		{Path: filepath.Join(homeDir, ".npm"), Label: "npm", Category: "npm", RebootSafe: true},
		{Path: filepath.Join(homeDir, ".bun", "install", "cache"), Label: "Bun", Category: "bun", RebootSafe: true},
		{Path: filepath.Join(homeDir, "Library", "Caches", "Yarn"), Label: "Yarn", Category: "yarn", RebootSafe: true},
		{Path: filepath.Join(homeDir, "Library", "pnpm", "store"), Label: "pnpm", Category: "pnpm", RebootSafe: true},
		{Path: filepath.Join(homeDir, "Library", "Caches", "pip"), Label: "pip", Category: "pip", RebootSafe: true},
		{Path: filepath.Join(homeDir, ".cargo", "registry", "cache"), Label: "Cargo", Category: "cargo", RebootSafe: true},
		{Path: filepath.Join(homeDir, ".gem"), Label: "Ruby Gems", Category: "ruby", RebootSafe: true},
		{Path: filepath.Join(homeDir, "Library", "Containers", "com.docker.docker"), Label: "Docker", Category: "docker", RebootSafe: true},
		{Path: filepath.Join(homeDir, ".local", "share", "containers"), Label: "Podman", Category: "podman", RebootSafe: true},
		{Path: "/usr/local/var/log/nginx", Label: "Nginx", Category: "nginx", RebootSafe: true},
		{Path: filepath.Join(homeDir, ".gradle", "caches"), Label: "Gradle", Category: "gradle", RebootSafe: true},
		{Path: filepath.Join(homeDir, ".m2", "repository"), Label: "Maven", Category: "maven", RebootSafe: true},
		{Path: filepath.Join(homeDir, "Library", "Android", "sdk"), Label: "Android", Category: "android", RebootSafe: true},
		{Path: filepath.Join(homeDir, "Library", "Caches", "Homebrew"), Label: "Homebrew", Category: "brew", RebootSafe: true},
		{Path: filepath.Join(homeDir, "Library", "Developer", "Xcode", "DerivedData"), Label: "Xcode", Category: "xcode", RebootSafe: true, ExtraPaths: []string{filepath.Join(homeDir, "Library", "Developer", "CoreSimulator", "Devices")}},
		{Path: filepath.Join(homeDir, ".composer", "cache"), Label: "Composer", Category: "composer", RebootSafe: true},
	}

	for i := range softwareLocations {
		softwareLocations[i].Detected = pathExists(softwareLocations[i].Path)
	}

	return append(coreLocations, softwareLocations...)
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

func HandleTmpAnalyseLocations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	locations := DiscoverLocations(homeDir)
	json.NewEncoder(w).Encode(locations)
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

	if runtime.GOOS != "darwin" {
		sendSSEEvent(w, "unsupported_platform", map[string]string{"os": runtime.GOOS})
		flusher.Flush()
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		sendSSEEvent(w, "server_error", map[string]string{"error": err.Error()})
		return
	}

	locations := DiscoverLocations(homeDir)

	if err := sendSSEEvent(w, "locations", locations); err != nil {
		return
	}
	flusher.Flush()

	var accumulatedSize, accumulatedReclaimable int64

	for i := range locations {
		if !locations[i].Detected {
			continue
		}

		primaryPath := locations[i].Path
		label := locations[i].Label
		rebootSafe := locations[i].RebootSafe
		extraPaths := locations[i].ExtraPaths

		allPaths := []string{primaryPath}
		allPaths = append(allPaths, extraPaths...)

		var totalSize, totalCount int64
		var extraSizes []int64
		var extraCounts []int64

		for pi, scanPath := range allPaths {
			fsys := os.DirFS(scanPath)
			size, count, err := ScanWithProgress(fsys, ".", func(curSize int64, curCount int64) {
				progress := BuildProgressPayload(label, curSize, curCount, accumulatedSize, accumulatedReclaimable, rebootSafe)
				sendSSEEvent(w, "progress", progress)
				flusher.Flush()
			})
			if err != nil {
				log.Printf("Error scanning %s (%s): %v (got %d bytes, %d files)", label, scanPath, err, size, count)
				size = 0
				count = 0
			}
			if pi == 0 {
				totalSize = size
				totalCount = count
			} else {
				totalSize += size
				totalCount += count
				extraSizes = append(extraSizes, size)
				extraCounts = append(extraCounts, count)
			}
		}

		locations[i].Size = totalSize
		locations[i].FileCount = totalCount
		if len(extraSizes) > 0 {
			locations[i].ExtraSizes = extraSizes
		}
		if len(extraCounts) > 0 {
			locations[i].ExtraCounts = extraCounts
		}

		accumulatedSize += totalSize
		if rebootSafe {
			accumulatedReclaimable += totalSize
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
