# Tmp Named Scan

Backend SSE tests for `GET /api/tmp-named-scan`: two-pass scan (fast `named`
discovery with zero shared fields, then async `named_enriched` from `analyse`),
package manager detection, shared size metrics, `scan_complete` / `done` event
order, and vendor scans that skip enrichment.

## Version

0.0.2

# DSN (Domain Specific Notion)

The **tmp-named-scan handler** walks git repositories under the user home directory,
finds directories whose basename matches the `name` query parameter (typically
`node_modules`), and streams results on one SSE connection in two passes.

**Pass 1** emits **named** events immediately per hit with path/size/repo fields,
fast **DetectPackageManager** (lockfiles, `packageManager` field, default npm), and zeroed shared columns. After
the walk it emits **summary**, then **scan_complete** so the UI can show rows
without waiting for analyse.

**Pass 2** (node_modules only) starts a background enrichment worker when the
scan begins. Each **named** hit is queued immediately; the worker runs
**analyse** and emits **named_enriched** as results arrive — these may
interleave with later **named** events and may appear before **scan_complete**.
After the walk: **summary**, **scan_complete**, wait for the worker, then
**done**. **Vendor** scans skip pass 2: **scan_complete** then **done** with no
**named_enriched** events.

For **node_modules** hits, **gitTracked** reports whether a sibling
`package.json` (same parent directory as the hit) is tracked in git via
`git ls-files --error-unmatch package.json`. The field appears on **named**,
**named_size**, and **named_enriched** SSE events.

Tests use a fake **HOME**, httptest SSE, **ParseNamedSSE** (extended for
`scan_complete` / `named_enriched`), and raw JSON helpers.

## Decision Tree

```
tmp-named-scan/
├── two-pass/
│   ├── fast-named-zero-shared/
│   ├── event-order/
│   ├── pipelined-enrichment-starts-early/
│   ├── concurrent-sse-writes/
│   ├── enrichment-populates-shared/
│   ├── async-sizing-after-scan-complete/
│   ├── client-disconnect-after-scan-complete/
│   └── vendor-skips-enrichment/
├── detect-package-manager/
│   ├── bun-lockb/
│   ├── pnpm-lock/
│   ├── npm-lock/
│   ├── yarn-lock/
│   ├── unknown/
│   ├── has-package-json/
│   ├── no-package-json/
│   ├── package-json-default-npm/
│   ├── package-json-package-manager-pnpm/
│   ├── package-json-package-manager-yarn/
│   ├── lockfile-wins-over-field/
│   └── nested-pnpm-node-modules/
├── shared-size/
│   ├── darwin-clone/
│   └── non-darwin/
├── git-tracked/
│   ├── tracked-package-json/
│   ├── untracked-package-json/
│   ├── no-sibling-package-json/
│   └── nested-pnpm-no-sibling/
└── sse-fields/
    └── present/
```

## Test Index

| Leaf | Op |
|------|-----|
| two-pass/fast-named-zero-shared | named-scan |
| two-pass/event-order | named-scan |
| two-pass/pipelined-enrichment-starts-early | named-scan |
| two-pass/concurrent-sse-writes | named-scan |
| two-pass/enrichment-populates-shared | named-scan |
| two-pass/async-sizing-after-scan-complete | named-scan-lifecycle |
| two-pass/client-disconnect-after-scan-complete | named-scan-disconnect |
| two-pass/vendor-skips-enrichment | named-scan |
| detect-package-manager/bun-lockb | named-scan |
| detect-package-manager/pnpm-lock | named-scan |
| detect-package-manager/npm-lock | named-scan |
| detect-package-manager/yarn-lock | named-scan |
| detect-package-manager/unknown | named-scan |
| detect-package-manager/has-package-json | named-scan |
| detect-package-manager/no-package-json | named-scan |
| detect-package-manager/package-json-default-npm | named-scan |
| detect-package-manager/package-json-package-manager-pnpm | named-scan |
| detect-package-manager/package-json-package-manager-yarn | named-scan |
| detect-package-manager/lockfile-wins-over-field | named-scan |
| detect-package-manager/nested-pnpm-node-modules | named-scan |
| shared-size/darwin-clone | named-scan |
| shared-size/non-darwin | named-scan |
| git-tracked/tracked-package-json | named-scan |
| git-tracked/untracked-package-json | named-scan |
| git-tracked/no-sibling-package-json | named-scan |
| git-tracked/nested-pnpm-no-sibling | named-scan |
| sse-fields/present | named-scan |

## How to Run

```sh
doctest vet ./tests/tmp-named-scan
doctest test ./tests/tmp-named-scan
```

```go
import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"disk-usage-analyser/server"
	"disk-usage-analyser/tmpfiles"

	"github.com/xhd2015/dot-pkgs/go-pkgs/git/scan_repo"
)

var _ = os.Setenv("GO111MODULE", "off")

type namedEventJSON struct {
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

type Request struct {
	Op      string
	Name    string
	HomeDir string
}

type Response struct {
	SSEOutput  string
	EventTypes []string
	NamedHits  []server.NamedHit
	Summary    *server.NamedScanSummary
	NamedJSON  []namedEventJSON
	NamedEnrichedJSON []namedEventJSON
	HasScanComplete               bool
	HasNamedEnriched              bool
	HasServerError                bool
	LastNamedBeforeScanComplete   bool
	ScanCompleteBeforeNamedEnriched bool
	NamedEnrichedBeforeDone       bool
	ScanCompleteBeforeDone        bool
	FirstNamedEnrichedIndex       int
	FirstScanCompleteIndex        int
	PostScanNamedHits             int64
	TotalNamedHits                int64
	HitsAtScanReturn              int64
	ScanDuration                  time.Duration
	ServerLog                     string
	DisconnectSawScanComplete     bool
}

func runNamedScan(req *Request) (*Response, error) {
	handler := http.HandlerFunc(server.HandleTmpNamedScan)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	name := req.Name
	if name == "" {
		name = "node_modules"
	}
	httpReq, err := http.NewRequest("GET", srv.URL+"/api/tmp-named-scan?name="+name, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	parsed, err := server.ParseNamedSSE(string(body))
	if err != nil {
		return nil, err
	}
	enriched := parseNamedEnrichedEventsFromSSE(string(body))
	events := parsed.EventTypes
	return &Response{
		SSEOutput:  parsed.SSEOutput,
		EventTypes: events,
		NamedHits:  parsed.NamedHits,
		Summary:    parsed.Summary,
		NamedJSON:  parseNamedEventsFromSSE(string(body)),
		NamedEnrichedJSON: enriched,
		HasScanComplete:               containsEvent(events, "scan_complete"),
		HasNamedEnriched:              len(enriched) > 0,
		HasServerError:                containsEvent(events, "server_error"),
		LastNamedBeforeScanComplete:   eventLastBefore(events, "named", "scan_complete"),
		ScanCompleteBeforeNamedEnriched: eventFirstBefore(events, "scan_complete", "named_enriched"),
		NamedEnrichedBeforeDone:       eventFirstBefore(events, "named_enriched", "done"),
		ScanCompleteBeforeDone:        eventFirstBefore(events, "scan_complete", "done"),
		FirstNamedEnrichedIndex:       firstEventIndex(events, "named_enriched"),
		FirstScanCompleteIndex:        firstEventIndex(events, "scan_complete"),
	}, nil
}

func runNamedScanLifecycleProbe(t *testing.T, req *Request) error {
	t.Helper()
	t.Setenv("HOME", req.HomeDir)
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	moduleRoot := filepath.Clean(filepath.Join(cwd, "..", "..", "..", ".."))
	cmd := exec.Command("go", "run", "./cmd/tmp-named-scan-lifecycle-probe", req.HomeDir)
	cmd.Dir = moduleRoot
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	out, err := cmd.CombinedOutput()
	t.Logf("lifecycle probe:\n%s", string(out))
	if err != nil {
		outStr := string(out)
		if strings.Contains(outStr, "atReturn=0") && strings.Contains(outStr, "postScan=") {
			return fmt.Errorf("async sizing finished after Scan returned:\n%s", outStr)
		}
		return fmt.Errorf("lifecycle probe failed: %v\n%s", err, outStr)
	}
	return nil
}

func runNamedScanLifecycle(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	t.Setenv("HOME", req.HomeDir)

	var scanDone atomic.Bool
	var postScanHits atomic.Int64
	var totalHits atomic.Int64

	name := req.Name
	if name == "" {
		name = "node_modules"
	}

	ctx := context.Background()
	scanStarted := time.Now()
	roots := []string{filepath.Join(req.HomeDir, "Projects")}
	summary, err := tmpfiles.Scan(ctx, tmpfiles.ScanOptions{
		Names: []string{name},
		Roots: roots,
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
	}, req.HomeDir, nil)
	scanDuration := time.Since(scanStarted)
	hitsAtReturn := totalHits.Load()
	scanDone.Store(true)
	time.Sleep(500 * time.Millisecond)

	t.Logf("scan duration: %s", scanDuration)

	return &Response{
		Summary: &server.NamedScanSummary{
			Repos:     summary.Repos,
			NamedHits: summary.NamedHits,
		},
		PostScanNamedHits: postScanHits.Load(),
		TotalNamedHits:    totalHits.Load(),
		HitsAtScanReturn:  hitsAtReturn,
		ScanDuration:      scanDuration,
	}, err
}

func runNamedScanDisconnectAfterScanComplete(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	t.Setenv("HOME", req.HomeDir)

	var logBuf bytes.Buffer
	prevLog := log.Writer()
	log.SetOutput(&logBuf)
	defer log.SetOutput(prevLog)

	handler := http.HandlerFunc(server.HandleTmpNamedScan)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	name := req.Name
	if name == "" {
		name = "node_modules"
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	httpReq, err := http.NewRequestWithContext(ctx, "GET", srv.URL+"/api/tmp-named-scan?name="+name, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	sawScanComplete := false
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "event: scan_complete") {
			sawScanComplete = true
			break
		}
	}
	cancel()
	resp.Body.Close()
	time.Sleep(5 * time.Second)

	return &Response{
		ServerLog:                 logBuf.String(),
		DisconnectSawScanComplete: sawScanComplete,
	}, nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Setenv("HOME", req.HomeDir)

	switch req.Op {
	case "named-scan":
		return runNamedScan(req)
	case "named-scan-lifecycle":
		if err := runNamedScanLifecycleProbe(t, req); err != nil {
			return nil, err
		}
		return &Response{}, nil
	case "named-scan-disconnect":
		return runNamedScanDisconnectAfterScanComplete(t, req)
	default:
		t.Fatalf("unknown test op: %q", req.Op)
		return nil, nil
	}
}
```