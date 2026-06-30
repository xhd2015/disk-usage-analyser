# Tmp Binaries Scan

Backend SSE tests for `GET /api/tmp-binaries-scan`: binary discovery in git repos,
Go/Mach-O/ELF classification, ignored directories, streaming order, summary, and tilde paths.

## Version

0.0.2

# DSN (Domain Specific Notion)

The **tmp-binaries-scan handler** wraps `tmpfiles` scan logic under the user home
directory. For each discovered git repository it walks regular files, skips ignored
basenames, classifies binaries via `detect.DetectFileType` with Go buildinfo precedence,
and streams **binary** SSE events carrying **BinaryHit** records immediately. A **summary**
event reports repo and binary counts plus total size; **done** marks completion. Scan root
defaults to `~` with unlimited depth in v1.

## Decision Tree

```
tmp-binaries-scan/
├── streaming/
│   └── binary-before-done/
├── classification/
│   ├── go-binary/
│   ├── macho-binary/
│   └── elf-binary/
├── discovery/
│   ├── multi-repos/
│   └── respect-ignore-dirs/
└── output/
    ├── summary-event/
    └── tilde-paths/
```

## Test Index

| Leaf | Op |
|------|-----|
| streaming/binary-before-done | binaries-sse-order |
| classification/go-binary | binaries-scan |
| classification/macho-binary | binaries-scan |
| classification/elf-binary | binaries-scan |
| discovery/multi-repos | binaries-scan |
| discovery/respect-ignore-dirs | binaries-scan |
| output/summary-event | binaries-scan |
| output/tilde-paths | binaries-scan |

## How to Run

```sh
doctest vet ./tests/tmp-binaries-scan
doctest test ./tests/tmp-binaries-scan
```

```go
import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"disk-usage-analyser/server"
)

var _ = os.Setenv("GO111MODULE", "off")

type Request struct {
	Op      string
	HomeDir string
}

type Response struct {
	SSEOutput        string
	EventTypes       []string
	Binaries         []server.BinaryHit
	Summary          *server.BinaryScanSummary
	BinaryBeforeDone bool
}

func runBinariesScan(req *Request) (*Response, error) {
	handler := http.HandlerFunc(server.HandleTmpBinariesScan)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	httpReq, err := http.NewRequest("GET", srv.URL+"/api/tmp-binaries-scan", nil)
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
	parsed, err := server.ParseBinariesSSE(string(body))
	if err != nil {
		return nil, err
	}
	return &Response{
		SSEOutput:  parsed.SSEOutput,
		EventTypes: parsed.EventTypes,
		Binaries:   parsed.Binaries,
		Summary:    parsed.Summary,
	}, nil
}

func runBinariesSSEOrder(req *Request) (*Response, error) {
	parsed, err := runBinariesScan(req)
	if err != nil {
		return nil, err
	}
	parsed.BinaryBeforeDone = server.EventBefore(parsed.EventTypes, "binary", "done")
	return parsed, nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Setenv("HOME", req.HomeDir)

	switch req.Op {
	case "binaries-scan":
		return runBinariesScan(req)
	case "binaries-sse-order":
		return runBinariesSSEOrder(req)
	default:
		t.Fatalf("unknown test op: %q", req.Op)
		return nil, nil
	}
}

func tilde(home, path string) string {
	if path == home {
		return "~"
	}
	if strings.HasPrefix(path, home+string(os.PathSeparator)) {
		return "~" + path[len(home):]
	}
	return path
}
```