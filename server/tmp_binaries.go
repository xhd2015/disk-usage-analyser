package server

import (
	"net/http"
	"os"
	"sync"

	"disk-usage-analyser/tmpfiles"

	"github.com/xhd2015/dot-pkgs/go-pkgs/git/scan_repo"
)

type BinaryHit struct {
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	SizeHuman string `json:"sizeHuman"`
	Kind      string `json:"kind"`
	TypeDesc  string `json:"typeDesc"`
	RepoPath  string `json:"repoPath"`
	RepoName  string `json:"repoName"`
}

type BinaryScanSummary struct {
	Repos      int    `json:"repos"`
	Binaries   int    `json:"binaries"`
	TotalSize  int64  `json:"totalSize"`
	TotalHuman string `json:"totalHuman"`
}

type binaryScanEntry struct {
	DisplayPath string
	AbsPath     string
	Size        int64
	Repo        scan_repo.Repo
}

var (
	binarySessionMu sync.Mutex
	binarySessions  = make(map[string]map[string]binaryScanEntry)
)

func clearBinarySession(homeDir string) {
	binarySessionMu.Lock()
	defer binarySessionMu.Unlock()
	binarySessions[homeDir] = make(map[string]binaryScanEntry)
}

func storeBinarySession(homeDir string, hit BinaryHit, repo scan_repo.Repo, absPath string) {
	binarySessionMu.Lock()
	defer binarySessionMu.Unlock()
	session, ok := binarySessions[homeDir]
	if !ok {
		session = make(map[string]binaryScanEntry)
		binarySessions[homeDir] = session
	}
	session[hit.Path] = binaryScanEntry{
		DisplayPath: hit.Path,
		AbsPath:     absPath,
		Size:        hit.Size,
		Repo:        repo,
	}
}

func lookupBinarySession(homeDir, displayPath string) (binaryScanEntry, bool) {
	binarySessionMu.Lock()
	defer binarySessionMu.Unlock()
	session, ok := binarySessions[homeDir]
	if !ok {
		return binaryScanEntry{}, false
	}
	entry, ok := session[displayPath]
	return entry, ok
}

func HandleTmpBinariesScan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		sendSSEEvent(w, "server_error", map[string]string{"error": err.Error()})
		flusher.Flush()
		return
	}

	clearBinarySession(homeDir)

	summary, err := tmpfiles.Scan(ctx, tmpfiles.ScanOptions{
		Roots:    []string{homeDir},
		MaxDepth: 0,
	}, homeDir, func(hit tmpfiles.BinaryHit, repo scan_repo.Repo) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		absPath := resolveTildePath(hit.Path, homeDir)
		serverHit := BinaryHit{
			Path:      hit.Path,
			Size:      hit.Size,
			SizeHuman: hit.SizeHuman,
			Kind:      hit.Kind,
			TypeDesc:  hit.TypeDesc,
			RepoPath:  hit.RepoPath,
			RepoName:  hit.RepoName,
		}
		storeBinarySession(homeDir, serverHit, repo, absPath)

		if err := sendSSEEvent(w, "binary", serverHit); err != nil {
			return err
		}
		flusher.Flush()
		return nil
	})
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		sendSSEEvent(w, "server_error", map[string]string{"error": err.Error()})
		flusher.Flush()
		return
	}

	scanSummary := BinaryScanSummary{
		Repos:      summary.Repos,
		Binaries:   summary.Binaries,
		TotalSize:  summary.TotalSize,
		TotalHuman: summary.TotalHuman,
	}
	if err := sendSSEEvent(w, "summary", scanSummary); err != nil {
		return
	}
	flusher.Flush()

	sendSSEEvent(w, "done", map[string]string{"status": "complete"})
	flusher.Flush()
}