package server

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"

	"disk-usage-analyser/tmpfiles"

	"github.com/xhd2015/dot-pkgs/go-pkgs/git/scan_repo"
)

type NamedHit struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	SizeHuman string `json:"sizeHuman"`
	RepoPath  string `json:"repoPath"`
	RepoName  string `json:"repoName"`
}

type NamedScanSummary struct {
	Repos      int    `json:"repos"`
	NamedHits  int    `json:"namedHits"`
	TotalSize  int64  `json:"totalSize"`
	TotalHuman string `json:"totalHuman"`
}

type namedScanEntry struct {
	DisplayPath string
	AbsPath     string
	Size        int64
	Repo        scan_repo.Repo
}

var (
	namedSessionMu sync.Mutex
	namedSessions  = make(map[string]map[string]namedScanEntry)
)

func clearNamedSession(homeDir string) {
	namedSessionMu.Lock()
	defer namedSessionMu.Unlock()
	namedSessions[homeDir] = make(map[string]namedScanEntry)
}

func storeNamedSession(homeDir string, hit NamedHit, repo scan_repo.Repo, absPath string) {
	namedSessionMu.Lock()
	defer namedSessionMu.Unlock()
	session, ok := namedSessions[homeDir]
	if !ok {
		session = make(map[string]namedScanEntry)
		namedSessions[homeDir] = session
	}
	session[hit.Path] = namedScanEntry{
		DisplayPath: hit.Path,
		AbsPath:     absPath,
		Size:        hit.Size,
		Repo:        repo,
	}
}

func lookupNamedSession(homeDir, displayPath string) (namedScanEntry, bool) {
	namedSessionMu.Lock()
	defer namedSessionMu.Unlock()
	session, ok := namedSessions[homeDir]
	if !ok {
		return namedScanEntry{}, false
	}
	entry, ok := session[displayPath]
	return entry, ok
}

func HandleTmpNamedScan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		sendSSEEvent(w, "server_error", map[string]string{"error": "name query parameter required"})
		flusher.Flush()
		return
	}

	ctx := r.Context()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		sendSSEEvent(w, "server_error", map[string]string{"error": err.Error()})
		flusher.Flush()
		return
	}

	clearNamedSession(homeDir)

	var (
		namedMutex   sync.Mutex
		namedCount   int
		totalSize    int64
		repoCount    int
	)

	summary, err := tmpfiles.ScanBinaries(ctx, tmpfiles.ScanOptions{
		Names:  []string{name},
		Roots:  []string{homeDir},
		OnRepo: func(repo scan_repo.Repo) error {
			namedMutex.Lock()
			repoCount++
			namedMutex.Unlock()
			return nil
		},
		OnNamedHit: func(hit tmpfiles.NamedHit, repo scan_repo.Repo) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			namedMutex.Lock()
			namedCount++
			totalSize += hit.Size
			namedMutex.Unlock()

			absPath := resolveTildePath(hit.Path, homeDir)

			serverHit := NamedHit{
				Path:      hit.Path,
				Name:      hit.Name,
				Size:      hit.Size,
				SizeHuman: hit.SizeHuman,
				RepoPath:  hit.RepoPath,
				RepoName:  hit.RepoName,
			}
			storeNamedSession(homeDir, serverHit, repo, absPath)

			if err := sendSSEEvent(w, "named", serverHit); err != nil {
				return err
			}
			flusher.Flush()
			return nil
		},
	}, homeDir, nil)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		sendSSEEvent(w, "server_error", map[string]string{"error": err.Error()})
		flusher.Flush()
		return
	}

	scanSummary := NamedScanSummary{
		Repos:      repoCount,
		NamedHits:  namedCount,
		TotalSize:  totalSize,
		TotalHuman: summary.TotalHuman,
	}
	if err := sendSSEEvent(w, "summary", scanSummary); err != nil {
		return
	}
	flusher.Flush()

	sendSSEEvent(w, "done", map[string]string{"status": "complete"})
	flusher.Flush()
}

type DeleteNamedRequest struct {
	Paths []string `json:"paths"`
}

type DeleteNamedResult struct {
	Deleted    []string        `json:"deleted"`
	Failed     []DeleteFailure `json:"failed"`
	FreedSize  int64           `json:"freedSize"`
	FreedHuman string          `json:"freedHuman"`
}

func HandleTmpNamedDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, `{"error":"home directory unavailable"}`, http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"failed to read request body"}`, http.StatusBadRequest)
		return
	}

	var req DeleteNamedRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	result := DeleteNamedResult{
		Deleted: []string{},
		Failed:  []DeleteFailure{},
	}

	for _, displayPath := range req.Paths {
		entry, ok := lookupNamedSession(homeDir, displayPath)
		if !ok {
			result.Failed = append(result.Failed, DeleteFailure{
				Path:  displayPath,
				Error: "path not in current scan results",
			})
			continue
		}

		absPath := entry.AbsPath

		info, err := os.Stat(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				result.Failed = append(result.Failed, DeleteFailure{
					Path:  displayPath,
					Error: "directory not found",
				})
				continue
			}
			result.Failed = append(result.Failed, DeleteFailure{
				Path:  displayPath,
				Error: err.Error(),
			})
			continue
		}

		if !info.IsDir() {
			result.Failed = append(result.Failed, DeleteFailure{
				Path:  displayPath,
				Error: "not a directory",
			})
			continue
		}

		if err := os.RemoveAll(absPath); err != nil {
			result.Failed = append(result.Failed, DeleteFailure{
				Path:  displayPath,
				Error: err.Error(),
			})
			continue
		}

		result.Deleted = append(result.Deleted, displayPath)
		result.FreedSize += entry.Size
	}

	result.FreedHuman = tmpfiles.FormatHumanSize(result.FreedSize)
	json.NewEncoder(w).Encode(result)
}
