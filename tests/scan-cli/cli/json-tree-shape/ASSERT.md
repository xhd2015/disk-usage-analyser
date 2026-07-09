## Expected

- Exit code 0.
- JSON object includes `path`, `totalSize`, `min`, `maxDepth`, and nested `tree`.
- No `items` key; no `threshold` key.
- `min` is 1048576 (1M); `maxDepth` is 24 (JSON capture default).
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
	for _, key := range []string{"path", "totalSize", "min", "maxDepth", "tree"} {
		if payload[key] == nil {
			t.Fatalf("json missing %q key: %s", key, line)
		}
	}
	if payload["items"] != nil {
		t.Fatal("json must not include items key")
	}
	if payload["threshold"] != nil {
		t.Fatal("json must not include threshold key (renamed to min)")
	}
	var min int64
	if err := json.Unmarshal(payload["min"], &min); err != nil || min != 1<<20 {
		t.Fatalf("min: expected %d, got %d err=%v", 1<<20, min, err)
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
