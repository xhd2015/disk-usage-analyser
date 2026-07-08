## Expected

- Exit code 0.
- JSON object includes `path`, `totalSize`, `threshold`, `maxDepth`, and nested `tree`.
- No `items` key.
- `threshold` is 1048576 (1M); `maxDepth` is 24 (JSON default).
- Stdout ends with a trailing blank line.

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
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d (err=%v)", resp.ExitCode, resp.Err)
	}
	content := strings.TrimRight(resp.Stdout, "\n")
	lines := strings.Split(content, "\n")
	if len(lines) != 1 {
		t.Fatalf("expected one JSON line, got %d lines:\n%s", len(lines), resp.Stdout)
	}
	var payload map[string]json.RawMessage
	if err := json.Unmarshal([]byte(lines[0]), &payload); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, resp.Stdout)
	}
	for _, key := range []string{"path", "totalSize", "threshold", "maxDepth", "tree"} {
		if payload[key] == nil {
			t.Fatalf("json missing %q key: %s", key, lines[0])
		}
	}
	if payload["items"] != nil {
		t.Fatal("json must not include items key")
	}
	var threshold int64
	if err := json.Unmarshal(payload["threshold"], &threshold); err != nil || threshold != 1<<20 {
		t.Fatalf("threshold: expected %d, got %d err=%v", 1<<20, threshold, err)
	}
	var maxDepth int
	if err := json.Unmarshal(payload["maxDepth"], &maxDepth); err != nil || maxDepth != 24 {
		t.Fatalf("maxDepth: expected 24, got %d err=%v", maxDepth, err)
	}
	var tree map[string]any
	if err := json.Unmarshal(payload["tree"], &tree); err != nil {
		t.Fatalf("json tree: %v", err)
	}
	if tree["name"] != "." {
		t.Fatalf("tree root name: expected %q, got %#v", ".", tree["name"])
	}
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```