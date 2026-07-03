## Expected

- First stdout line is valid JSON with `"type":"named"`, `"name":"node_modules"`, `"size":6`, `"sizeHuman":"6 B"`.
- Required fields: `path`, `name`, `size`, `sizeHuman`, `repoPath`, `repoName`.
- Second line is the text summary: `Found 0 binaries, 1 named entries, total 6 B`.

## Side Effects

- None outside the temporary fixture tree.

## Errors

- No error is returned.

## Exit Code

- 0

```go
import (
	"encoding/json"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("unexpected scan error: %v", resp.Err)
	}
	lines := strings.Split(strings.TrimSpace(resp.Stdout), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected JSON hit and summary, got %d lines:\n%s", len(lines), resp.Stdout)
	}
	var hit map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &hit); err != nil {
		t.Fatalf("first line is not valid JSON: %v\n%s", err, lines[0])
	}
	if got, _ := hit["type"].(string); got != "named" {
		t.Fatalf("expected type=named, got %#v", hit["type"])
	}
	if got, _ := hit["name"].(string); got != "node_modules" {
		t.Fatalf("expected name=node_modules, got %#v", hit["name"])
	}
	if got, _ := hit["size"].(float64); got != 6 {
		t.Fatalf("expected size=6, got %v", hit["size"])
	}
	if got, _ := hit["sizeHuman"].(string); got != "6 B" {
		t.Fatalf("expected sizeHuman='6 B', got %#v", hit["sizeHuman"])
	}
	if got, _ := hit["path"].(string); got == "" {
		t.Fatal("expected non-empty path")
	}
	if got, _ := hit["repoPath"].(string); got == "" {
		t.Fatal("expected non-empty repoPath")
	}
	if got, _ := hit["repoName"].(string); got != "app" {
		t.Fatalf("expected repoName=app, got %#v", hit["repoName"])
	}
	if !strings.HasPrefix(lines[1], "Found 0 binaries, 1 named entries, total 6 B") {
		t.Fatalf("bad summary line: %q", lines[1])
	}
}

```
