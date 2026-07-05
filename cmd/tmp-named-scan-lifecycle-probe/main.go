// Command tmp-named-scan-lifecycle-probe checks that tmpfiles.Scan does not return
// before async emitNamedDirSized goroutines call OnNamedHit.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"disk-usage-analyser/tmpfiles"

	"github.com/xhd2015/dot-pkgs/go-pkgs/git/scan_repo"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <homeDir>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	homeDir := os.Args[1]

	var scanDone atomic.Bool
	var postScanHits atomic.Int64
	var totalHits atomic.Int64

	ctx := context.Background()
	scanStarted := time.Now()
	summary, err := tmpfiles.Scan(ctx, tmpfiles.ScanOptions{
		Names: []string{"node_modules"},
		Roots: []string{filepath.Join(homeDir, "Projects")},
		OnNamedPreview: func(hit tmpfiles.NamedHit, repo scan_repo.Repo) error {
			return nil
		},
		OnNamedHit: func(hit tmpfiles.NamedHit, repo scan_repo.Repo) error {
			totalHits.Add(1)
			if scanDone.Load() {
				postScanHits.Add(1)
			}
			return nil
		},
	}, homeDir, nil)
	scanDuration := time.Since(scanStarted)
	hitsAtReturn := totalHits.Load()
	scanDone.Store(true)
	time.Sleep(500 * time.Millisecond)

	if err != nil {
		fmt.Fprintf(os.Stderr, "scan error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("scan=%s atReturn=%d total=%d postScan=%d repos=%d\n",
		scanDuration, hitsAtReturn, totalHits.Load(), postScanHits.Load(), summary.Repos)

	const minHits int64 = 2
	if totalHits.Load() < minHits {
		fmt.Fprintf(os.Stderr, "total hits %d < %d\n", totalHits.Load(), minHits)
		os.Exit(1)
	}
	if hitsAtReturn < totalHits.Load() || postScanHits.Load() > 0 {
		os.Exit(2)
	}
}