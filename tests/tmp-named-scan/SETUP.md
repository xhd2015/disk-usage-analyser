# Scenario

**Feature**: tmp-named-scan shared fixture harness

```
GET /api/tmp-named-scan?name=node_modules -> pass1 named (zero shared) -> summary -> scan_complete -> pass2 named_enriched -> done
```

## Preconditions

- Tests must never scan the real user home directory.
- Each leaf creates a temporary fake home and sets `HOME` before httptest.
- Git repositories are directories containing `.git/`.
- `node_modules` fixtures are directories with at least one file for non-zero size.
- Shared-size darwin leaves use `cp -c` (APFS clone) with isolated pnpm store and bun cache env vars.

## Steps

1. Create a fresh fixture home for each leaf.
2. Place `node_modules` trees inside git repos under the fake home.
3. Set `req.Op` to `named-scan` and run the httptest SSE harness.

## Context

- Raw JSON helpers parse `named` and `named_enriched` event payloads.
- Event-order helpers use `EventTypes` indices for `scan_complete` / `named_enriched` ordering.

```go
import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const file4K = 4096

func Setup(t *testing.T, req *Request) error {
	req.HomeDir = filepath.Join(t.TempDir(), "home")
	if err := os.MkdirAll(req.HomeDir, 0755); err != nil {
		return err
	}
	if req.Op == "" {
		req.Op = "named-scan"
	}
	if req.Name == "" {
		req.Name = "node_modules"
	}
	return nil
}

func repo(t *testing.T, home string, rel string) string {
	t.Helper()
	dir := filepath.Join(home, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0755); err != nil {
		t.Fatalf("create repo %s: %v", rel, err)
	}
	return dir
}

func writeFile(t *testing.T, base string, rel string, data []byte, mode os.FileMode) string {
	t.Helper()
	path := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir parent %s: %v", rel, err)
	}
	if err := os.WriteFile(path, data, mode); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
	return path
}

func writeSizedFile(t *testing.T, base string, rel string, size int64) string {
	t.Helper()
	data := make([]byte, size)
	for i := range data {
		data[i] = 'x'
	}
	return writeFile(t, base, rel, data, 0644)
}

func nodeModulesRepo(t *testing.T, home string, rel string) string {
	t.Helper()
	app := repo(t, home, rel)
	writeFile(t, app, "node_modules/pkg/dep.txt", []byte("fixture"), 0644)
	return app
}

func vendorRepo(t *testing.T, home string, rel string) string {
	t.Helper()
	app := repo(t, home, rel)
	writeFile(t, app, "vendor/pkg/dep.txt", []byte("fixture"), 0644)
	return app
}

func parseNamedEventsFromSSE(body string) []namedEventJSON {
	var hits []namedEventJSON
	lines := strings.Split(body, "\n")
	var currentEvent string
	for _, line := range lines {
		if strings.HasPrefix(line, "event: ") {
			currentEvent = strings.TrimPrefix(line, "event: ")
		}
		if strings.HasPrefix(line, "data: ") && currentEvent == "named" {
			data := strings.TrimPrefix(line, "data: ")
			var hit namedEventJSON
			if err := json.Unmarshal([]byte(data), &hit); err == nil {
				hits = append(hits, hit)
			}
		}
	}
	return hits
}

func namedEventDataObjects(t *testing.T, body string) []map[string]interface{} {
	t.Helper()
	return sseEventDataObjects(t, body, "named")
}

func parseNamedEnrichedEventsFromSSE(body string) []namedEventJSON {
	var hits []namedEventJSON
	lines := strings.Split(body, "\n")
	var currentEvent string
	for _, line := range lines {
		if strings.HasPrefix(line, "event: ") {
			currentEvent = strings.TrimPrefix(line, "event: ")
		}
		if strings.HasPrefix(line, "data: ") && currentEvent == "named_enriched" {
			data := strings.TrimPrefix(line, "data: ")
			var hit namedEventJSON
			if err := json.Unmarshal([]byte(data), &hit); err == nil {
				hits = append(hits, hit)
			}
		}
	}
	return hits
}

func namedEnrichedEventDataObjects(t *testing.T, body string) []map[string]interface{} {
	t.Helper()
	return sseEventDataObjects(t, body, "named_enriched")
}

func namedSizeEventDataObjects(t *testing.T, body string) []map[string]interface{} {
	t.Helper()
	return sseEventDataObjects(t, body, "named_size")
}

func sseEventDataObjects(t *testing.T, body string, eventName string) []map[string]interface{} {
	t.Helper()
	var objs []map[string]interface{}
	lines := strings.Split(body, "\n")
	var currentEvent string
	for _, line := range lines {
		if strings.HasPrefix(line, "event: ") {
			currentEvent = strings.TrimPrefix(line, "event: ")
		}
		if strings.HasPrefix(line, "data: ") && currentEvent == eventName {
			data := strings.TrimPrefix(line, "data: ")
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(data), &obj); err != nil {
				t.Fatalf("unmarshal %s event: %v\n%s", eventName, err, data)
			}
			objs = append(objs, obj)
		}
	}
	return objs
}

func firstEventIndex(events []string, event string) int {
	for i, e := range events {
		if e == event {
			return i
		}
	}
	return -1
}

func lastEventIndex(events []string, event string) int {
	last := -1
	for i, e := range events {
		if e == event {
			last = i
		}
	}
	return last
}

func containsEvent(events []string, event string) bool {
	return firstEventIndex(events, event) >= 0
}

func eventFirstBefore(events []string, before, after string) bool {
	beforeIdx := firstEventIndex(events, before)
	afterIdx := firstEventIndex(events, after)
	return beforeIdx >= 0 && afterIdx >= 0 && beforeIdx < afterIdx
}

func eventLastBefore(events []string, before, after string) bool {
	beforeIdx := lastEventIndex(events, before)
	afterIdx := firstEventIndex(events, after)
	return beforeIdx >= 0 && afterIdx >= 0 && beforeIdx < afterIdx
}

func cpClone(t *testing.T, src string, dest string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		t.Fatalf("mkdir clone parent: %v", err)
	}
	if err := exec.Command("cp", "-c", src, dest).Run(); err != nil {
		t.Fatalf("cp -c %s -> %s: %v", src, dest, err)
	}
}

const pnpmHash = "abc123deadbeef"
const bunPkgDir = "pkg@1.0.0@@@1"

func pnpmStoreFilesDir(tempRoot string) string {
	return filepath.Join(tempRoot, "store", "files")
}

func bunCacheRoot(tempRoot string) string {
	return filepath.Join(tempRoot, "cache")
}

func writePnpmStoreFile(t *testing.T, storeFilesDir string, size int64) string {
	t.Helper()
	rel := filepath.Join("aa", pnpmHash)
	return writeSizedFile(t, storeFilesDir, rel, size)
}

func writeBunCacheFile(t *testing.T, cacheRoot string, size int64) string {
	t.Helper()
	rel := filepath.Join(bunPkgDir, "file")
	return writeSizedFile(t, cacheRoot, rel, size)
}
```