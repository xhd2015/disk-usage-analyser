package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

type TmpLocation struct {
	Path           string             `json:"path"`
	Label          string             `json:"label"`
	Category       string             `json:"category"`
	Size           int64              `json:"size"`
	FileCount      int64              `json:"fileCount"`
	RebootSafe     bool               `json:"rebootSafe"`
	Reclaimable    bool               `json:"reclaimable"`
	Detected       bool               `json:"detected"`
	ExtraPaths     []string           `json:"-"`
	BreakdownItems []TmpBreakdownItem `json:"breakdownItems,omitempty"`
	RuntimeItems   []TmpRuntimeItem   `json:"runtimeItems,omitempty"`
	VmInternal     *TmpVmInternal     `json:"vmInternal,omitempty"`
}

type TmpBreakdownItem struct {
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	FileCount int64  `json:"fileCount"`
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

func tildePath(homeDir string, path string) string {
	if path == homeDir {
		return "~"
	}
	if strings.HasPrefix(path, homeDir+string(filepath.Separator)) {
		return "~" + path[len(homeDir):]
	}
	return path
}

func resolveTildePath(path string, homeDir string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

func DiscoverLocations(homeDir string) []TmpLocation {
	coreLocations := []TmpLocation{
		{Path: filepath.Join(homeDir, ".Trash"), Label: "User Trash", Category: "trash", RebootSafe: true, Detected: true},
		{Path: filepath.Join(homeDir, "Library", "Caches"), Label: "User Caches", Category: "cache", RebootSafe: true, Detected: true},
		{Path: filepath.Join(homeDir, "Library", "Logs"), Label: "User Logs", Category: "log", RebootSafe: true, Detected: true},
		{Path: "/private/var/vm/", Label: "Swap", Category: "swap", RebootSafe: true, Reclaimable: false, Detected: true},
		{Path: "/tmp", Label: "System Temp", Category: "temp", RebootSafe: false, Detected: true},
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
		{Path: filepath.Join(homeDir, ".local", "share", "opencode", "snapshot"), Label: "OpenCode", Category: "opencode", RebootSafe: true,
			ExtraPaths: []string{
				filepath.Join(homeDir, ".local", "share", "opencode", "project"),
				filepath.Join(homeDir, ".local", "share", "opencode", "tool-output"),
				filepath.Join(homeDir, ".local", "share", "opencode", "storage"),
				filepath.Join(homeDir, ".local", "share", "opencode", "log"),
				filepath.Join(homeDir, ".cache", "opencode"),
				filepath.Join(homeDir, ".local", "state", "opencode"),
			}},
		{Path: filepath.Join(homeDir, ".claude", "plugins"), Label: "Claude Code", Category: "claude", RebootSafe: true,
			ExtraPaths: []string{
				filepath.Join(homeDir, ".claude", "telemetry"),
				filepath.Join(homeDir, ".claude", "todos"),
				filepath.Join(homeDir, ".claude", "cache"),
				filepath.Join(homeDir, ".claude", "backups"),
			}},
		{Path: filepath.Join(homeDir, ".codex"), Label: "Codex (OpenAI)", Category: "codex", RebootSafe: true,
			ExtraPaths: []string{
				filepath.Join(homeDir, "Library", "Application Support", "codex"),
			}},
		{Path: filepath.Join(homeDir, "Library", "Application Support", "Cursor"), Label: "Cursor", Category: "cursor", RebootSafe: true,
			ExtraPaths: []string{
				filepath.Join(homeDir, "Library", "Application Support", "Caches", "cursor-updater"),
				filepath.Join(homeDir, "Library", "Caches", "cursor-compile-cache"),
			}},
	}

	for i := range softwareLocations {
		softwareLocations[i].Detected = pathExists(softwareLocations[i].Path)
	}

	all := append(coreLocations, softwareLocations...)
	for i := range all {
		all[i].Reclaimable = all[i].RebootSafe
		if all[i].Category == "swap" {
			all[i].Reclaimable = false
		}
	}
	for i := range all {
		all[i].Path = tildePath(homeDir, all[i].Path)
		for j := range all[i].ExtraPaths {
			all[i].ExtraPaths[j] = tildePath(homeDir, all[i].ExtraPaths[j])
		}
		items := []TmpBreakdownItem{{Path: all[i].Path}}
		for j := range all[i].ExtraPaths {
			items = append(items, TmpBreakdownItem{Path: all[i].ExtraPaths[j]})
		}
		all[i].BreakdownItems = items
	}
	return all
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
		if loc.Reclaimable {
			reclaimable += loc.Size
		}
	}
	return TmpSummary{
		Locations:       locations,
		TotalSize:       total,
		ReclaimableSize: reclaimable,
	}
}

func BuildProgressPayload(label string, curSize, curCount int64, accumulatedSize, accumulatedReclaimable int64, reclaimable bool) map[string]interface{} {
	totalSize := accumulatedSize + curSize
	reclaimableSize := accumulatedReclaimable
	if reclaimable {
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

func BuildBreakdownProgressPayload(
	label string,
	completedSizes []int64,
	completedCounts []int64,
	activeIndex int,
	activeSize int64,
	activeCount int64,
	accumulatedSize int64,
	accumulatedReclaimable int64,
	reclaimable bool,
	breakdownPath string,
) map[string]interface{} {
	var cardSize int64
	var cardFileCount int64
	for _, s := range completedSizes {
		cardSize += s
	}
	for _, c := range completedCounts {
		cardFileCount += c
	}
	cardSize += activeSize
	cardFileCount += activeCount

	totalSize := accumulatedSize + cardSize
	reclaimableSize := accumulatedReclaimable
	if reclaimable {
		reclaimableSize += cardSize
	}

	payload := map[string]interface{}{
		"label":              label,
		"size":               cardSize,
		"fileCount":          cardFileCount,
		"breakdownIndex":     activeIndex,
		"breakdownSize":      activeSize,
		"breakdownFileCount": activeCount,
		"totalSize":          totalSize,
		"reclaimableSize":    reclaimableSize,
	}
	if breakdownPath != "" {
		payload["breakdownPath"] = breakdownPath
	}
	return payload
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

		label := locations[i].Label
		reclaimable := locations[i].Reclaimable

		var totalSize, totalCount int64

		if label == "npm" {
			npmResolvedPath := resolveTildePath(locations[i].Path, homeDir)
			entries, err := os.ReadDir(npmResolvedPath)
			if err == nil && len(entries) > 0 {
				var subItems []TmpBreakdownItem
				var completedSizes []int64
				var completedCounts []int64
				for _, e := range entries {
					if !e.IsDir() {
						continue
					}
					subPath := filepath.Join(npmResolvedPath, e.Name())
					tildeSubPath := tildePath(homeDir, subPath)
					fsys := os.DirFS(subPath)
					size, count, err := ScanWithProgress(fsys, ".", func(curSize int64, curCount int64) {
						progress := BuildBreakdownProgressPayload(
							label,
							completedSizes,
							completedCounts,
							len(completedSizes),
							curSize,
							curCount,
							accumulatedSize,
							accumulatedReclaimable,
							reclaimable,
							tildeSubPath,
						)
						sendSSEEvent(w, "progress", progress)
						flusher.Flush()
					})
					if err != nil {
						log.Printf("Error scanning npm subdir %s (%s): %v", label, subPath, err)
						size = 0
						count = 0
					}
					totalSize += size
					totalCount += count
					completedSizes = append(completedSizes, size)
					completedCounts = append(completedCounts, count)
					subItems = append(subItems, TmpBreakdownItem{
						Path:      tildeSubPath,
						Size:      size,
						FileCount: count,
					})
				}
				locations[i].BreakdownItems = subItems
			} else {
				scanPath := resolveTildePath(locations[i].BreakdownItems[0].Path, homeDir)
				fsys := os.DirFS(scanPath)
				size, count, err := ScanWithProgress(fsys, ".", func(curSize int64, curCount int64) {
					progress := BuildProgressPayload(label, curSize, curCount, accumulatedSize, accumulatedReclaimable, reclaimable)
					sendSSEEvent(w, "progress", progress)
					flusher.Flush()
				})
				if err != nil {
					log.Printf("Error scanning %s (%s): %v (got %d bytes, %d files)", label, scanPath, err, size, count)
					size = 0
					count = 0
				}
				totalSize = size
				totalCount = count
				locations[i].BreakdownItems[0].Size = size
				locations[i].BreakdownItems[0].FileCount = count
			}
		} else {
			useBreakdown := len(locations[i].BreakdownItems) >= 2
			var completedSizes []int64
			var completedCounts []int64
			for bi := range locations[i].BreakdownItems {
				scanPath := resolveTildePath(locations[i].BreakdownItems[bi].Path, homeDir)

				fsys := os.DirFS(scanPath)
				size, count, err := ScanWithProgress(fsys, ".", func(curSize int64, curCount int64) {
					var progress map[string]interface{}
					if useBreakdown {
						progress = BuildBreakdownProgressPayload(
							label,
							completedSizes,
							completedCounts,
							bi,
							curSize,
							curCount,
							accumulatedSize,
							accumulatedReclaimable,
							reclaimable,
							"",
						)
					} else {
						progress = BuildProgressPayload(label, curSize, curCount, accumulatedSize, accumulatedReclaimable, reclaimable)
					}
					sendSSEEvent(w, "progress", progress)
					flusher.Flush()
				})
				if err != nil {
					log.Printf("Error scanning %s (%s): %v (got %d bytes, %d files)", label, scanPath, err, size, count)
					size = 0
					count = 0
				}

				locations[i].BreakdownItems[bi].Size = size
				locations[i].BreakdownItems[bi].FileCount = count
				totalSize += size
				totalCount += count
				if useBreakdown {
					completedSizes = append(completedSizes, size)
					completedCounts = append(completedCounts, count)
				}
			}
		}

		locations[i].Size = totalSize
		locations[i].FileCount = totalCount

		switch locations[i].Category {
		case "podman":
			if vm, _ := CollectPodmanVmInternal(); vm != nil {
				locations[i].VmInternal = vm
			}
			if items, _ := CollectPodmanRuntimeViaSSH(); len(items) > 0 {
				locations[i].RuntimeItems = items
			}
		case "docker":
			if items, _ := CollectRuntimeStats(locations[i].Category); len(items) > 0 {
				locations[i].RuntimeItems = items
			}
		}

		accumulatedSize += totalSize
		if reclaimable {
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

func isSSEClientDisconnect(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET) {
		return true
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		if errors.Is(opErr.Err, syscall.EPIPE) || errors.Is(opErr.Err, syscall.ECONNRESET) {
			return true
		}
	}
	var syscallErr *os.SyscallError
	if errors.As(err, &syscallErr) {
		if errors.Is(syscallErr.Err, syscall.EPIPE) || errors.Is(syscallErr.Err, syscall.ECONNRESET) {
			return true
		}
	}
	msg := err.Error()
	return strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "connection reset by peer") ||
		strings.Contains(msg, "write on closed")
}

func sendSSEEvent(w http.ResponseWriter, event string, data interface{}) error {
	jsonData, _ := json.Marshal(data)
	_, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, jsonData)
	if err != nil {
		if !isSSEClientDisconnect(err) {
			log.Printf("Error sending SSE event %s: %v", event, err)
		}
		return err
	}
	return nil
}
