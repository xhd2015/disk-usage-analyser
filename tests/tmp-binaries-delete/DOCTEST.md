# Tmp Binaries Delete

Backend POST tests for `POST /api/tmp-binaries-delete`: single and batch delete,
validation (non-binary, directory, not-in-scan-results), already-deleted paths,
and partial batch success.

## Version

0.0.2

# DSN (Domain Specific Notion)

The **tmp-binaries-delete handler** accepts a **DeleteBinariesRequest** with absolute
or `~/` paths from the **current scan session**. Each path is re-validated as a regular
binary file (same `classifyFile` logic as scan) before unlink. Paths not present in the
active scan results are rejected. The response is **DeleteBinariesResult** with per-path
`deleted` and `failed` entries plus `freedSize`. Partial delete is allowed. Remote-backed
filesystem paths are rejected via the same `remotefs` guard as scan.

## Decision Tree

```
tmp-binaries-delete/
├── delete/
│   ├── single-binary/
│   └── multiple-binaries/
├── validation/
│   ├── reject-non-binary/
│   ├── reject-directory/
│   ├── reject-not-in-scan-results/
│   └── already-deleted/
└── batch/
    └── partial-success/
```

## Test Index

| Leaf | Op |
|------|-----|
| delete/single-binary | delete-binaries |
| delete/multiple-binaries | delete-binaries |
| validation/reject-non-binary | delete-binaries |
| validation/reject-directory | delete-binaries |
| validation/reject-not-in-scan-results | delete-binaries |
| validation/already-deleted | delete-binaries |
| batch/partial-success | delete-binaries |

## How to Run

```sh
doctest vet ./tests/tmp-binaries-delete
doctest test ./tests/tmp-binaries-delete
```

```go
import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"disk-usage-analyser/server"
)

type Request struct {
	Op                    string
	HomeDir               string
	Paths                 []string
	ScanFirst             bool
	SkipScan              bool
	ExtraPaths            []string
	PreDelete             []string
	OverwriteBeforeDelete map[string]string
	AddAfterScan          map[string][]byte
	RemoveAfterScan       []string
}

type Response struct {
	StatusCode   int
	Result       *server.DeleteBinariesResult
	ScanPaths    []string
	FileExists   map[string]bool
	Body         string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Setenv("HOME", req.HomeDir)

	var scanPaths []string
	if !req.SkipScan {
		paths, err := runBinariesScanCollect(t, req.HomeDir)
		if err != nil {
			return nil, err
		}
		scanPaths = paths
	}

	for rel, data := range req.AddAfterScan {
		abs := filepath.Join(req.HomeDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(abs, data, 0755); err != nil {
			return nil, err
		}
	}
	for rel, content := range req.OverwriteBeforeDelete {
		abs := filepath.Join(req.HomeDir, filepath.FromSlash(rel))
		if err := os.WriteFile(abs, []byte(content), 0755); err != nil {
			return nil, err
		}
	}
	for _, p := range req.RemoveAfterScan {
		abs := resolvePath(req.HomeDir, p)
		_ = os.Remove(abs)
	}

	deletePaths := req.Paths
	if req.ScanFirst && len(deletePaths) == 0 {
		deletePaths = scanPaths
	}
	deletePaths = append(deletePaths, req.ExtraPaths...)

	result, status, body, err := postDelete(t, deletePaths)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		StatusCode: status,
		Result:     result,
		ScanPaths:  scanPaths,
		Body:       body,
		FileExists: map[string]bool{},
	}
	for _, p := range deletePaths {
		abs := resolvePath(req.HomeDir, p)
		_, statErr := os.Stat(abs)
		resp.FileExists[p] = statErr == nil
	}
	for _, p := range req.PreDelete {
		abs := resolvePath(req.HomeDir, p)
		_, statErr := os.Stat(abs)
		key := p
		if !strings.HasPrefix(key, "~") {
			key = tilde(req.HomeDir, abs)
		}
		resp.FileExists[key] = statErr == nil
	}
	return resp, nil
}

func runBinariesScanCollect(t *testing.T, homeDir string) ([]string, error) {
	t.Setenv("HOME", homeDir)
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

	var paths []string
	lines := strings.Split(string(body), "\n")
	var currentEvent string
	for _, line := range lines {
		if strings.HasPrefix(line, "event: ") {
			currentEvent = strings.TrimPrefix(line, "event: ")
		}
		if strings.HasPrefix(line, "data: ") && currentEvent == "binary" {
			data := strings.TrimPrefix(line, "data: ")
			var hit server.BinaryHit
			if err := json.Unmarshal([]byte(data), &hit); err == nil {
				paths = append(paths, hit.Path)
			}
		}
	}
	return paths, nil
}

func postDelete(t *testing.T, paths []string) (*server.DeleteBinariesResult, int, string, error) {
	payload, err := json.Marshal(server.DeleteBinariesRequest{Paths: paths})
	if err != nil {
		return nil, 0, "", err
	}
	handler := http.HandlerFunc(server.HandleTmpBinariesDelete)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	httpReq, err := http.NewRequest("POST", srv.URL+"/api/tmp-binaries-delete", bytes.NewReader(payload))
	if err != nil {
		return nil, 0, "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, 0, "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, "", err
	}
	var result server.DeleteBinariesResult
	if resp.StatusCode == http.StatusOK {
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, resp.StatusCode, string(body), err
		}
	}
	return &result, resp.StatusCode, string(body), nil
}

func resolvePath(home, p string) string {
	if strings.HasPrefix(p, "~/") {
		return filepath.Join(home, p[2:])
	}
	return p
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