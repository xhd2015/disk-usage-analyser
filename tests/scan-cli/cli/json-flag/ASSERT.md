## Expected

- Exit code 0.
- Stdout is one JSON object with `path`, `totalSize`, `min`, and nested `tree`.
- `min` is 1; `totalSize` is 500; `tree.children` include `big.txt` (400) and `small.txt` (100).
- Children sorted by size descending.
- No `threshold` key.
- Stdout ends with a trailing blank line.

## Exit Code

- 0

```go
import (
	"encoding/json"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d (err=%v)", resp.ExitCode, resp.Err)
	}
	line := firstJSONObjectLine(t, resp.Stdout)
	var payload map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, resp.Stdout)
	}
	for _, key := range []string{"path", "totalSize", "min", "tree"} {
		if payload[key] == nil {
			t.Fatalf("json missing %q key: %s", key, line)
		}
	}
	if payload["items"] != nil {
		t.Fatalf("json must not include items key: %s", line)
	}
	if payload["threshold"] != nil {
		t.Fatalf("json must not include threshold key (renamed to min): %s", line)
	}
	var path string
	if err := json.Unmarshal(payload["path"], &path); err != nil || path != req.FixtureDir {
		t.Fatalf("json path: expected %q, got %q err=%v", req.FixtureDir, path, err)
	}
	var min int64
	if err := json.Unmarshal(payload["min"], &min); err != nil || min != 1 {
		t.Fatalf("json min: expected 1, got %d err=%v", min, err)
	}
	var totalSize int64
	if err := json.Unmarshal(payload["totalSize"], &totalSize); err != nil || totalSize != 500 {
		t.Fatalf("json totalSize: expected 500, got %d err=%v", totalSize, err)
	}
	var tree map[string]any
	if err := json.Unmarshal(payload["tree"], &tree); err != nil {
		t.Fatalf("json tree: %v", err)
	}
	if name, _ := tree["name"].(string); name != "." {
		t.Fatalf("tree root name: expected %q, got %q", ".", name)
	}
	children, ok := tree["children"].([]any)
	if !ok || len(children) != 2 {
		t.Fatalf("tree children: expected 2, got %#v", tree["children"])
	}
	first, ok := children[0].(map[string]any)
	if !ok {
		t.Fatalf("first child: %#v", children[0])
	}
	if name, _ := first["name"].(string); name != "big.txt" {
		t.Fatalf("first child should be big.txt (largest), got %q", name)
	}
	if size, ok := first["size"].(float64); !ok || int64(size) != 400 {
		t.Fatalf("big.txt size: expected 400, got %#v", first["size"])
	}
	if isDir, ok := first["isDir"].(bool); !ok || isDir {
		t.Fatalf("big.txt isDir: expected false, got %#v", first["isDir"])
	}
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
