package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"disk-usage-analyser/analyse"
	"disk-usage-analyser/tmpfiles"

	"github.com/xhd2015/dot-pkgs/go-pkgs/git/scan_repo"
)

const namedEnrichWorkerCount = 6

type NamedHit struct {
	Path            string `json:"path"`
	Name            string `json:"name"`
	Size            int64  `json:"size"`
	SizeHuman       string `json:"sizeHuman"`
	RepoPath        string `json:"repoPath"`
	RepoName        string `json:"repoName"`
	PackageManager  string `json:"packageManager"`
	HasPackageJSON  bool   `json:"hasPackageJson"`
	GitTracked      bool   `json:"gitTracked"`
	PnpmSharedSize  int64  `json:"pnpmSharedSize"`
	PnpmSharedHuman string `json:"pnpmSharedHuman"`
	BunSharedSize   int64  `json:"bunSharedSize"`
	BunSharedHuman  string `json:"bunSharedHuman"`
	SharedSize      int64  `json:"sharedSize"`
	SharedHuman     string `json:"sharedHuman"`
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

type namedEnrichmentJob struct {
	DisplayPath    string
	AbsPath        string
	PackageManager string
	HasPackageJSON bool
	GitTracked     bool
}

func applyPackageManager(hit *NamedHit, absPath, scanName string) {
	if scanName == "node_modules" {
		hit.PackageManager = DetectPackageManager(absPath)
	} else {
		hit.PackageManager = "unknown"
	}
}

func applyHasPackageJSON(hit *NamedHit, absPath, scanName string) {
	if scanName == "node_modules" {
		hit.HasPackageJSON = DetectHasPackageJSON(absPath)
	} else {
		hit.HasPackageJSON = false
	}
}

func applyGitTracked(hit *NamedHit, absPath, scanName string) {
	if scanName == "node_modules" {
		hit.GitTracked = DetectGitTrackedPackageJSON(absPath)
	} else {
		hit.GitTracked = false
	}
}

func applySharedMetrics(hit *NamedHit, absPath string) {
	applySharedMetricsLogged(hit, absPath, nil)
}

func applySharedMetricsLogged(hit *NamedHit, absPath string, log func(format string, args ...any)) {
	zeroSharedFields(hit)

	t0 := time.Now()
	result, err := analyse.Analyse(absPath)
	elapsed := time.Since(t0)
	if log != nil {
		if err != nil {
			log("analyse.Analyse duration=%s err=%v", elapsed, err)
		} else {
			summary := result.Summary
			log("analyse.Analyse duration=%s pnpm_shared=%s (%d B) bun_shared=%s (%d B)",
				elapsed, summary.PnpmSharedHuman, summary.PnpmSharedSize,
				summary.BunSharedHuman, summary.BunSharedSize)
		}
	}
	if err != nil {
		return
	}

	summary := result.Summary
	hit.PnpmSharedSize = summary.PnpmSharedSize
	hit.PnpmSharedHuman = summary.PnpmSharedHuman
	hit.BunSharedSize = summary.BunSharedSize
	hit.BunSharedHuman = summary.BunSharedHuman
	hit.SharedSize = summary.PnpmSharedSize + summary.BunSharedSize
	hit.SharedHuman = formatAnalyseHumanSize(hit.SharedSize)
	if log != nil {
		log("shared_total=%s (%d B)", hit.SharedHuman, hit.SharedSize)
	}
}

// NamedEnrichedTraceFields returns the named_enriched SSE fields used for PM/pkgjson
// without running analyse (for CLI diagnostics).
func NamedEnrichedTraceFields(displayPath, packageManager, absPath string, hasPackageJSON bool) NamedHit {
	hit := NamedHit{
		Path:           displayPath,
		PackageManager: packageManager,
		HasPackageJSON: hasPackageJSON,
	}
	zeroSharedFields(&hit)
	return hit
}

func buildNamedEnriched(displayPath, packageManager, absPath string, hasPackageJSON, gitTracked bool) NamedHit {
	hit := NamedHit{
		Path:           displayPath,
		PackageManager: packageManager,
		HasPackageJSON: hasPackageJSON,
		GitTracked:     gitTracked,
	}
	applySharedMetrics(&hit, absPath)
	return hit
}

func zeroSharedFields(hit *NamedHit) {
	hit.PnpmSharedSize = 0
	hit.PnpmSharedHuman = "0 B"
	hit.BunSharedSize = 0
	hit.BunSharedHuman = "0 B"
	hit.SharedSize = 0
	hit.SharedHuman = "0 B"
}

// formatAnalyseHumanSize mirrors analyse.formatHumanSize for SSE shared totals.
func formatAnalyseHumanSize(size int64) string {
	if size <= 1024 {
		return fmt.Sprintf("%d B", size)
	}
	units := []string{"K", "M", "G", "T", "P"}
	value := float64(size)
	for _, unit := range units {
		value /= 1024
		if value < 1024 {
			if value == float64(int64(value)) {
				return fmt.Sprintf("%.0f%s", value, unit)
			}
			return fmt.Sprintf("%.1f%s", value, unit)
		}
	}
	return fmt.Sprintf("%.1fE", value/1024)
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
		namedMutex sync.Mutex
		namedCount int
		totalSize  int64
		repoCount  int
	)

	var writeMu sync.Mutex
	var jobCh chan namedEnrichmentJob
	var enrichWG sync.WaitGroup

	emit := func(event string, payload interface{}) error {
		writeMu.Lock()
		defer writeMu.Unlock()
		if err := sendSSEEvent(w, event, payload); err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}

	scanStarted := time.Now()
	enrichedCount := 0
	var enrichedCountMu sync.Mutex

	if name == "node_modules" {
		jobCh = make(chan namedEnrichmentJob, 256)
		for i := 0; i < namedEnrichWorkerCount; i++ {
			enrichWG.Add(1)
			go func() {
				defer enrichWG.Done()
				for job := range jobCh {
					if ctx.Err() != nil {
						return
					}
					log.Printf("[named-enrich] start path=%s", job.DisplayPath)
					started := time.Now()
					enriched := buildNamedEnriched(job.DisplayPath, job.PackageManager, job.AbsPath, job.HasPackageJSON, job.GitTracked)
					duration := time.Since(started)
					log.Printf("[named-enrich] done path=%s duration=%s shared=%s", job.DisplayPath, duration, enriched.SharedHuman)
					if ctx.Err() != nil {
						return
					}
					if err := emit("named_enriched", enriched); err != nil {
						return
					}
					enrichedCountMu.Lock()
					enrichedCount++
					enrichedCountMu.Unlock()
				}
			}()
		}
	}

	buildPass1Hit := func(hit tmpfiles.NamedHit, repo scan_repo.Repo) (NamedHit, string) {
		absPath := resolveTildePath(hit.Path, homeDir)
		serverHit := NamedHit{
			Path:      hit.Path,
			Name:      hit.Name,
			Size:      hit.Size,
			SizeHuman: hit.SizeHuman,
			RepoPath:  hit.RepoPath,
			RepoName:  hit.RepoName,
		}
		applyPackageManager(&serverHit, absPath, name)
		applyHasPackageJSON(&serverHit, absPath, name)
		applyGitTracked(&serverHit, absPath, name)
		zeroSharedFields(&serverHit)
		return serverHit, absPath
	}

	scanRoots := []string{homeDir}
	if name == "node_modules" {
		scanRoots = tmpfiles.PrioritizedScanRoots(homeDir)
	}

	summary, err := tmpfiles.Scan(ctx, tmpfiles.ScanOptions{
		Names:  []string{name},
		Roots:  scanRoots,
		OnNamedPreview: func(hit tmpfiles.NamedHit, repo scan_repo.Repo) error {
			if name != "node_modules" {
				return nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			absPath := resolveTildePath(hit.Path, homeDir)
			serverHit := NamedHit{
				Path:      hit.Path,
				Name:      hit.Name,
				Size:      hit.Size,
				SizeHuman: hit.SizeHuman,
				RepoPath:  hit.RepoPath,
				RepoName:  hit.RepoName,
			}
			applyPackageManager(&serverHit, absPath, name)
			applyHasPackageJSON(&serverHit, absPath, name)
			applyGitTracked(&serverHit, absPath, name)
			zeroSharedFields(&serverHit)
			log.Printf("[named-scan] pass1 preview path=%s", serverHit.Path)
			return emit("named", serverHit)
		},
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

			serverHit, absPath := buildPass1Hit(hit, repo)
			storeNamedSession(homeDir, serverHit, repo, absPath)

			if name == "node_modules" {
				log.Printf("[named-scan] pass1 hit path=%s size=%d", serverHit.Path, serverHit.Size)
				if serverHit.Size > 0 {
					if err := emit("named_size", map[string]interface{}{
						"path":           serverHit.Path,
						"size":           serverHit.Size,
						"sizeHuman":      serverHit.SizeHuman,
						"hasPackageJson": serverHit.HasPackageJSON,
						"gitTracked":     serverHit.GitTracked,
					}); err != nil {
						return err
					}
				}
				job := namedEnrichmentJob{
					DisplayPath:    serverHit.Path,
					AbsPath:        absPath,
					PackageManager: serverHit.PackageManager,
					HasPackageJSON: serverHit.HasPackageJSON,
					GitTracked:     serverHit.GitTracked,
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case jobCh <- job:
				default:
					select {
					case <-ctx.Done():
						return ctx.Err()
					case jobCh <- job:
					}
				}
				return nil
			}

			if err := sendSSEEvent(w, "named", serverHit); err != nil {
				return err
			}
			flusher.Flush()
			return nil
		},
	}, homeDir, nil)
	if err != nil {
		if name == "node_modules" {
			close(jobCh)
			enrichWG.Wait()
		}
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

	if name != "node_modules" {
		if err := sendSSEEvent(w, "summary", scanSummary); err != nil {
			return
		}
		flusher.Flush()

		if err := sendSSEEvent(w, "scan_complete", map[string]string{"status": "scan_complete"}); err != nil {
			return
		}
		flusher.Flush()

		sendSSEEvent(w, "done", map[string]string{"status": "complete"})
		flusher.Flush()
		return
	}

	close(jobCh)

	log.Printf("[named-scan] scan_complete hits=%d", namedCount)

	if err := emit("summary", scanSummary); err != nil {
		return
	}
	if err := emit("scan_complete", map[string]string{"status": "scan_complete"}); err != nil {
		return
	}

	enrichWG.Wait()

	if ctx.Err() != nil {
		return
	}
	enrichedCountMu.Lock()
	finalEnriched := enrichedCount
	enrichedCountMu.Unlock()
	log.Printf("[named-scan] done enriched=%d duration=%s", finalEnriched, time.Since(scanStarted))
	emit("done", map[string]string{"status": "complete"})
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
