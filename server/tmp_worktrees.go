package server

import (
	"context"
	"io/fs"
	"net/http"
	"os"
	"sort"
	"strings"

	"disk-usage-analyser/tmpfiles"

	"github.com/xhd2015/dot-pkgs/go-pkgs/git/scan_repo"
)

type WorktreeRepoRow struct {
	RepoPath  string `json:"repoPath"`
	RepoName  string `json:"repoName"`
	Size      int64  `json:"size"`
	SizeHuman string `json:"sizeHuman"`
	FileCount int64  `json:"fileCount"`
}

type WorktreeHit struct {
	RepoPath  string `json:"repoPath"`
	RepoName  string `json:"repoName"`
	Path      string `json:"path"`
	Head      string `json:"head"`
	IsMain    bool   `json:"isMain"`
	Size      int64  `json:"size"`
	SizeHuman string `json:"sizeHuman"`
	FileCount int64  `json:"fileCount"`
}

type WorktreeScanSummary struct {
	Repos      int    `json:"repos"`
	Worktrees  int    `json:"worktrees"`
	TotalSize  int64  `json:"totalSize"`
	TotalHuman string `json:"totalHuman"`
}

func HandleTmpWorktreesScan(w http.ResponseWriter, r *http.Request) {
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

	var repoCount, worktreeCount int
	var totalSize int64

	_, err = scan_repo.Scan(ctx, scan_repo.Options{
		Roots:         []string{homeDir},
		MaxDepth:      0,
		ListWorktrees: true,
		OnRepo: func(repo scan_repo.Repo) error {
			if repo.RepoType != scan_repo.RepoTypeMain || len(repo.Worktrees) == 0 {
				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			repoCount++
			displayRepoPath := tmpfiles.DisplayPath(repo.Path, homeDir)

			worktrees := sortWorktreesMainFirst(repo.Worktrees)
			var mainSize, mainCount int64
			for _, wt := range worktrees {
				if !wt.IsMain {
					continue
				}
				size, count, err := sizeWorktreeCheckout(ctx, wt.Path)
				if err != nil {
					break
				}
				mainSize, mainCount = size, count
				break
			}

			repoRow := WorktreeRepoRow{
				RepoPath:  displayRepoPath,
				RepoName:  repo.Name,
				Size:      mainSize,
				SizeHuman: tmpfiles.FormatHumanSize(mainSize),
				FileCount: mainCount,
			}
			if err := sendSSEEvent(w, "repo", repoRow); err != nil {
				return err
			}
			flusher.Flush()
			totalSize += mainSize

			for _, wt := range worktrees {
				if wt.IsMain {
					continue
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				size, count, err := sizeWorktreeCheckout(ctx, wt.Path)
				if err != nil {
					continue
				}

				hit := WorktreeHit{
					RepoPath:  displayRepoPath,
					RepoName:  repo.Name,
					Path:      tmpfiles.DisplayPath(wt.Path, homeDir),
					Head:      normalizeWorktreeHead(wt.Head),
					IsMain:    false,
					Size:      size,
					SizeHuman: tmpfiles.FormatHumanSize(size),
					FileCount: count,
				}

				if err := sendSSEEvent(w, "worktree", hit); err != nil {
					return err
				}
				flusher.Flush()

				worktreeCount++
				totalSize += size
			}
			return nil
		},
	})
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		sendSSEEvent(w, "server_error", map[string]string{"error": err.Error()})
		flusher.Flush()
		return
	}

	summary := WorktreeScanSummary{
		Repos:      repoCount,
		Worktrees:  worktreeCount,
		TotalSize:  totalSize,
		TotalHuman: tmpfiles.FormatHumanSize(totalSize),
	}
	if err := sendSSEEvent(w, "summary", summary); err != nil {
		return
	}
	flusher.Flush()

	sendSSEEvent(w, "done", map[string]string{"status": "complete"})
	flusher.Flush()
}

func sortWorktreesMainFirst(worktrees []scan_repo.Worktree) []scan_repo.Worktree {
	sorted := append([]scan_repo.Worktree(nil), worktrees...)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].IsMain != sorted[j].IsMain {
			return sorted[i].IsMain
		}
		return sorted[i].Path < sorted[j].Path
	})
	return sorted
}

func normalizeWorktreeHead(head string) string {
	head = strings.TrimSpace(head)
	if strings.HasPrefix(head, "refs/heads/") {
		return strings.TrimPrefix(head, "refs/heads/")
	}
	return head
}

func sizeWorktreeCheckout(ctx context.Context, checkoutPath string) (int64, int64, error) {
	select {
	case <-ctx.Done():
		return 0, 0, ctx.Err()
	default:
	}

	fsys := os.DirFS(checkoutPath)
	return calculateDirSize(ctx, fsys, ".")
}

func calculateDirSize(ctx context.Context, fsys fs.FS, root string) (int64, int64, error) {
	var size int64
	var count int64
	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return nil
			}
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if d.IsDir() {
			if path != root && d.Name() == ".git" {
				return fs.SkipDir
			}
			return nil
		}
		info, err := d.Info()
		if err != nil {
			if os.IsPermission(err) {
				return nil
			}
			return err
		}
		size += info.Size()
		count++
		return nil
	})
	return size, count, err
}