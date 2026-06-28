# Tmp Analyse Test Cases

Backend tests for the tmp-analyse feature: location discovery, filesystem scanning,
SSE streaming, runtime stats (Docker/Podman), and live breakdown progress.

## Version

0.0.2

# DSN (Domain Specific Notion)

The **tmp-analyse handler** discovers cache and log locations on macOS, then scans
each detected path while streaming **SSE events** to the frontend. **DiscoverLocations**
builds the location catalog with breakdown rows for multi-path tools. **CalculateSize**
and **ScanWithProgress** walk filesystems and emit byte totals. **BuildProgressPayload**
(and breakdown extensions) shape per-card progress events; **BuildSummary** aggregates
global totals. For Docker and Podman cards, **CollectRuntimeStats** invokes
`docker|podman system df --format json`, **ParseSystemDFJSON** and **ParseHumanSize**
normalize CLI output into **TmpRuntimeItem** records filtered to Images and Build Cache,
attached on the final **location** event. The frontend renders cards, breakdown rows,
and optional runtime sections from these events.

## Test Tree

```
tmp-analyse-test-cases/
├── verify-discover-locations … verify-npm-dynamic-breakdown
├── runtime-stats/DOCTEST.md          # nested root — Docker/Podman runtime stats
├── breakdown-progress/DOCTEST.md     # nested root — live breakdown progress
├── podman-vm-internal/DOCTEST.md     # nested root — Podman VM internal storage
```

Nested trees (self-contained roots):

```
runtime-stats/
│   ├── parse-system-df/
│   │   ├── docker-ndjson/
│   │   ├── podman-ndjson/
│   │   └── malformed-line-rejected/
│   ├── parse-human-size/
│   │   ├── megabytes-decimal/
│   │   ├── gigabytes-compact/
│   │   ├── bytes-and-zero/
│   │   └── invalid-input/
│   ├── filter-runtime-items/
│   │   └── keeps-images-and-build-cache/
│   ├── collect-runtime-stats/
│   │   ├── docker-mock-success/
│   │   ├── podman-mock-success/
│   │   ├── binary-missing/
│   │   └── command-fails/
│   └── sse-runtime-items/
│       └── docker-location-event/

```

## Test Index

| Leaf | Op |
|------|-----|
| verify-discover-locations | discover-locations |
| verify-software-locations | software-locations |
| verify-detected-by-existence | detected-by-existence |
| verify-multi-path-locations | multi-path-locations |
| verify-initial-locations-event | initial-locations-event |
| verify-unsupported-platform | unsupported-platform |
| verify-extra-path-scan | extra-path-scan |
| verify-locations-rest-endpoint | locations-rest-endpoint |
| verify-calculate-size | calculate-size |
| verify-empty-dir | empty-dir |
| verify-nested-dirs | nested-dirs |
| verify-summary-totals | summary-totals |
| verify-sse-format | sse-format |
| verify-progress-stream | progress-stream |
| verify-scan-with-partial-error | scan-with-partial-error |
| verify-totals-in-progress | totals-in-progress |
| verify-tilde-paths | tilde-paths |
| verify-swap-location | swap-location |
| verify-swap-not-reclaimable | swap-not-reclaimable |
| verify-npm-dynamic-breakdown | npm-dynamic-breakdown |

See `runtime-stats/DOCTEST.md`, `breakdown-progress/DOCTEST.md`, and `podman-vm-internal/DOCTEST.md` for nested test indices.

## How to Run

```sh
doctest vet ./tests/tmp-analyse-test-cases
doctest test ./tests/tmp-analyse-test-cases
doctest test -v ./tests/tmp-analyse-test-cases/runtime-stats
doctest test -v ./tests/tmp-analyse-test-cases/breakdown-progress
doctest test -v ./tests/tmp-analyse-test-cases/podman-vm-internal
```

```go
import (
	"bufio"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"disk-usage-analyser/server"
)

type Request struct {
	Op       string
	HomeDir  string
	FS       interface{}
}

type Response struct {
	Locations        []server.TmpLocation
	Size             int64
	FileCount        int64
	TotalSize        int64
	ReclaimableSize  int64
	SSEOutput        string
	CategoryCount    map[string]int
	DetectedCount    int
	NotDetectedCount int
	ExtraPaths       []string
	ExtraSizes       []int64
	ExtraCounts      []int64
}

func Run(t *testing.T, req *Request) (*Response, error) {
	switch req.Op {
	case "discover-locations":
		return runDiscoverLocations(req)
	case "software-locations":
		return runSoftwareLocations(req)
	case "detected-by-existence":
		return runDetectedByExistence(req)
	case "multi-path-locations":
		return runMultiPathLocations(req)
	case "initial-locations-event":
		return runInitialLocationsEvent(req)
	case "unsupported-platform":
		return runUnsupportedPlatform(req)
	case "extra-path-scan":
		return runExtraPathScan(req)
	case "locations-rest-endpoint":
		return runLocationsRestEndpoint(req)
	case "calculate-size":
		return runCalculateSize(req)
	case "empty-dir":
		return runEmptyDir(req)
	case "nested-dirs":
		return runNestedDirs(req)
	case "summary-totals":
		return runSummaryTotals(req)
	case "sse-format":
		return runSSEFormat(req)
	case "progress-stream":
		return runProgressStream(req)
	case "scan-with-partial-error":
		return runScanWithPartialError(t, req)
	case "totals-in-progress":
		return runTotalsInProgress(req)
	case "tilde-paths":
		return runTildePaths(req)
	case "swap-location":
		return runSwapLocation(req)
	case "swap-not-reclaimable":
		return runSwapNotReclaimable(req)
	case "npm-dynamic-breakdown":
		return runNpmDynamicBreakdown(req)
	default:
		t.Fatalf("unknown test op: %q", req.Op)
		return nil, nil
	}
}

func runDiscoverLocations(req *Request) (*Response, error) {
	locations := server.DiscoverLocations(req.HomeDir)
	categoryCount := make(map[string]int)
	detectedCount := 0
	notDetectedCount := 0
	for _, loc := range locations {
		categoryCount[loc.Category]++
		if loc.Detected {
			detectedCount++
		} else {
			notDetectedCount++
		}
	}
	return &Response{
		Locations:        locations,
		CategoryCount:    categoryCount,
		DetectedCount:    detectedCount,
		NotDetectedCount: notDetectedCount,
	}, nil
}

func runSoftwareLocations(req *Request) (*Response, error) {
	locations := server.DiscoverLocations(req.HomeDir)
	coreCategories := map[string]bool{"trash": true, "temp": true, "cache": true, "log": true, "swap": true}
	var softwareLocs []server.TmpLocation
	for _, loc := range locations {
		if !coreCategories[loc.Category] {
			softwareLocs = append(softwareLocs, loc)
		}
	}
	return &Response{Locations: softwareLocs}, nil
}

func runDetectedByExistence(req *Request) (*Response, error) {
	locations := server.DiscoverLocations(req.HomeDir)
	categoryCount := make(map[string]int)
	detectedCount := 0
	notDetectedCount := 0
	for _, loc := range locations {
		categoryCount[loc.Category]++
		if loc.Detected {
			detectedCount++
		} else {
			notDetectedCount++
		}
	}
	return &Response{
		Locations:        locations,
		DetectedCount:    detectedCount,
		NotDetectedCount: notDetectedCount,
		CategoryCount:    categoryCount,
	}, nil
}

func runMultiPathLocations(req *Request) (*Response, error) {
	locations := server.DiscoverLocations(req.HomeDir)
	coreCategories := map[string]bool{"trash": true, "temp": true, "cache": true, "log": true, "swap": true}
	var softwareLocs []server.TmpLocation
	var goLoc, xcodeLoc *server.TmpLocation
	for i, loc := range locations {
		if !coreCategories[loc.Category] {
			softwareLocs = append(softwareLocs, loc)
			if loc.Category == "go" {
				goLoc = &locations[i]
			}
			if loc.Category == "xcode" {
				xcodeLoc = &locations[i]
			}
		}
	}
	extraPaths := []string{}
	if goLoc != nil {
		extraPaths = goLoc.ExtraPaths
	}
	if xcodeLoc != nil {
		extraPaths = append(extraPaths, xcodeLoc.ExtraPaths...)
	}
	return &Response{Locations: softwareLocs, ExtraPaths: extraPaths}, nil
}

func runInitialLocationsEvent(req *Request) (*Response, error) {
	handler := http.HandlerFunc(server.HandleTmpAnalyse)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	httpReq, err := http.NewRequest("GET", srv.URL+"/api/tmp-analyse", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	eventCount := 0
	var firstEventType, firstEventData string

	for scanner.Scan() && eventCount < 1 {
		line := scanner.Text()
		if strings.HasPrefix(line, "event: ") {
			firstEventType = strings.TrimPrefix(line, "event: ")
		}
		if strings.HasPrefix(line, "data: ") {
			firstEventData = strings.TrimPrefix(line, "data: ")
			eventCount++
		}
	}

	var locations []map[string]interface{}
	if err := json.Unmarshal([]byte(firstEventData), &locations); err != nil {
		return &Response{SSEOutput: firstEventType + " | " + firstEventData}, nil
	}
	detectedCount := 0
	notDetectedCount := 0
	for _, loc := range locations {
		if detected, ok := loc["detected"].(bool); ok && detected {
			detectedCount++
		} else {
			notDetectedCount++
		}
	}
	return &Response{
		SSEOutput:        firstEventType,
		DetectedCount:    detectedCount,
		NotDetectedCount: notDetectedCount,
	}, nil
}

func runUnsupportedPlatform(req *Request) (*Response, error) {
	handler := http.HandlerFunc(server.HandleTmpAnalyse)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	httpReq, err := http.NewRequest("GET", srv.URL+"/api/tmp-analyse", nil)
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
	return &Response{SSEOutput: string(body)}, nil
}

func runExtraPathScan(req *Request) (*Response, error) {
	primaryFS := req.FS.(fs.FS)
	extraFS := fstest.MapFS{
		"x.txt": &fstest.MapFile{Data: make([]byte, 500)},
		"y.txt": &fstest.MapFile{Data: make([]byte, 300)},
	}
	primarySize, primaryCount, err := server.CalculateSize(primaryFS, ".")
	if err != nil {
		return nil, err
	}
	extraSize, extraCount, err := server.CalculateSize(extraFS, ".")
	if err != nil {
		return nil, err
	}
	return &Response{
		Size:        primarySize,
		FileCount:   primaryCount,
		ExtraSizes:  []int64{extraSize},
		ExtraCounts: []int64{extraCount},
	}, nil
}

func runLocationsRestEndpoint(req *Request) (*Response, error) {
	handler := http.HandlerFunc(server.HandleTmpAnalyseLocations)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/tmp-analyse-locations")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	contentType := resp.Header.Get("Content-Type")

	var locations []server.TmpLocation
	if err := json.Unmarshal(body, &locations); err != nil {
		return nil, err
	}
	categoryCount := make(map[string]int)
	detectedCount := 0
	notDetectedCount := 0
	for _, loc := range locations {
		categoryCount[loc.Category]++
		if loc.Detected {
			detectedCount++
		} else {
			notDetectedCount++
		}
	}
	return &Response{
		Locations:        locations,
		CategoryCount:    categoryCount,
		DetectedCount:    detectedCount,
		NotDetectedCount: notDetectedCount,
		SSEOutput:        contentType,
	}, nil
}

func runCalculateSize(req *Request) (*Response, error) {
	fsys := req.FS.(fs.FS)
	size, count, err := server.CalculateSize(fsys, ".")
	if err != nil {
		return nil, err
	}
	return &Response{Size: size, FileCount: count}, nil
}

func runEmptyDir(req *Request) (*Response, error) {
	return runCalculateSize(req)
}

func runNestedDirs(req *Request) (*Response, error) {
	return runCalculateSize(req)
}

func runSummaryTotals(req *Request) (*Response, error) {
	locations := []server.TmpLocation{
		{Path: "/Users/x/.Trash", Label: "User Trash", Category: "trash", Size: 1000, FileCount: 5, RebootSafe: true, Reclaimable: true},
		{Path: "/Users/x/Library/Caches", Label: "User Caches", Category: "cache", Size: 500, FileCount: 10, RebootSafe: true, Reclaimable: true},
		{Path: "/tmp", Label: "System Temp", Category: "temp", Size: 2000, FileCount: 3, RebootSafe: false, Reclaimable: false},
		{Path: "/Users/x/Library/Logs", Label: "User Logs", Category: "log", Size: 300, FileCount: 8, RebootSafe: true, Reclaimable: true},
	}
	summary := server.BuildSummary(locations)
	return &Response{TotalSize: summary.TotalSize, ReclaimableSize: summary.ReclaimableSize}, nil
}

func runSSEFormat(req *Request) (*Response, error) {
	handler := http.HandlerFunc(server.HandleTmpAnalyse)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	httpReq, err := http.NewRequest("GET", srv.URL+"/api/tmp-analyse", nil)
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
	return &Response{SSEOutput: string(body)}, nil
}

func runProgressStream(req *Request) (*Response, error) {
	fsys := req.FS.(fs.FS)
	var progressSizes []int64
	size, count, err := server.ScanWithProgress(fsys, ".", func(s int64, c int64) {
		progressSizes = append(progressSizes, s)
	})
	if err != nil {
		return nil, err
	}
	resp := &Response{Size: size, FileCount: count, TotalSize: int64(len(progressSizes))}
	if len(progressSizes) > 0 {
		resp.ReclaimableSize = progressSizes[len(progressSizes)-1]
	}
	return resp, nil
}

func runScanWithPartialError(t *testing.T, req *Request) (*Response, error) {
	dir, err := os.MkdirTemp("", "scan-err-test")
	if err != nil {
		return nil, err
	}
	defer func() {
		os.Chmod(filepath.Join(dir, "locked"), 0700)
		os.RemoveAll(dir)
	}()

	os.WriteFile(filepath.Join(dir, "a.txt"), make([]byte, 100), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), make([]byte, 200), 0644)
	locked := filepath.Join(dir, "locked")
	os.Mkdir(locked, 0700)
	os.WriteFile(filepath.Join(locked, "c.txt"), make([]byte, 50), 0644)
	os.Chmod(locked, 0000)

	fsys := os.DirFS(dir)
	var progressSizes []int64
	size, count, scanErr := server.ScanWithProgress(fsys, ".", func(s int64, c int64) {
		progressSizes = append(progressSizes, s)
	})
	hadErr := int64(0)
	if scanErr != nil {
		hadErr = 1
	}
	return &Response{
		Size:            size,
		FileCount:       count,
		TotalSize:       int64(len(progressSizes)),
		ReclaimableSize: hadErr,
	}, nil
}

func runTotalsInProgress(req *Request) (*Response, error) {
	p1 := server.BuildProgressPayload("Loc B", 200, 10, 1000, 1000, true)
	p2 := server.BuildProgressPayload("Loc C", 500, 20, 1200, 1200, false)
	return &Response{
		Size:            p1["totalSize"].(int64),
		FileCount:       p1["reclaimableSize"].(int64),
		TotalSize:       p2["totalSize"].(int64),
		ReclaimableSize: p2["reclaimableSize"].(int64),
	}, nil
}

func runTildePaths(req *Request) (*Response, error) {
	return &Response{Locations: server.DiscoverLocations(req.HomeDir)}, nil
}

func runSwapLocation(req *Request) (*Response, error) {
	locations := server.DiscoverLocations(req.HomeDir)
	resp := &Response{Locations: locations}
	for _, loc := range locations {
		if loc.Category == "swap" {
			resp.Size = loc.Size
			resp.FileCount = loc.FileCount
			resp.SSEOutput = loc.Label
			break
		}
	}
	return resp, nil
}

func runSwapNotReclaimable(req *Request) (*Response, error) {
	locations := []server.TmpLocation{
		{Path: "/Users/x/.Trash", Label: "User Trash", Category: "trash", Size: 1000, FileCount: 5, RebootSafe: true, Reclaimable: true},
		{Path: "/private/var/vm/", Label: "Swap", Category: "swap", Size: 2048, FileCount: 3, RebootSafe: true, Reclaimable: false},
		{Path: "/tmp", Label: "System Temp", Category: "temp", Size: 500, FileCount: 2, RebootSafe: false, Reclaimable: false},
	}
	summary := server.BuildSummary(locations)
	return &Response{TotalSize: summary.TotalSize, ReclaimableSize: summary.ReclaimableSize}, nil
}

func runNpmDynamicBreakdown(req *Request) (*Response, error) {
	npmDir := filepath.Join(req.HomeDir, ".npm")
	entries, err := os.ReadDir(npmDir)
	if err != nil {
		return nil, err
	}
	items := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			items = append(items, filepath.Join(npmDir, e.Name()))
		}
	}
	return &Response{ExtraPaths: items, DetectedCount: len(items)}, nil
}
```