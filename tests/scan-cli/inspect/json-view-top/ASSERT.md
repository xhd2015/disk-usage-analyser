## Expected

- Exit code 0.
- ViewResult JSON: `tree` present, `matches` length 2, `min` field, no `threshold`.
- Matches ordered by size desc: `huge.bin` then `mid.bin` (root skipped).
- Trailing blank line.

## Exit Code

- 0

```go
import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d (err=%v)\n%s", resp.ExitCode, resp.Err, resp.Stdout)
	}
	line := firstJSONObjectLine(t, resp.Stdout)
	var payload map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, resp.Stdout)
	}
	if payload["threshold"] != nil {
		t.Fatal("must not emit threshold")
	}
	if payload["tree"] == nil || payload["matches"] == nil {
		t.Fatalf("ViewResult needs tree and matches: %s", line)
	}
	var matches []map[string]any
	if err := json.Unmarshal(payload["matches"], &matches); err != nil {
		t.Fatalf("matches: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d: %#v", len(matches), matches)
	}
	if name, _ := matches[0]["name"].(string); name != "huge.bin" {
		t.Fatalf("first match should be huge.bin, got %#v", matches[0])
	}
	if name, _ := matches[1]["name"].(string); name != "mid.bin" {
		t.Fatalf("second match should be mid.bin, got %#v", matches[1])
	}
	hugePath := filepath.Join(req.FixtureDir, "huge.bin")
	if path, _ := matches[0]["path"].(string); path != hugePath {
		t.Fatalf("huge path: expected %q, got %q", hugePath, path)
	}
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
