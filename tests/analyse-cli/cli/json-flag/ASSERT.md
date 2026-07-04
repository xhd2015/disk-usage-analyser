## Expected

- Exit code 0.
- Stdout is one JSON object with `root`, `rows`, and `summary`.
- `summary.size` is 2048; `rows` has one entry for `sub`.
- `summary` and row include `pnpmSharedSize` (0), `pnpmSharedHuman`, `bunSharedSize` (0), and `bunSharedHuman`.

## Exit Code

- 0

```go
import (
	"encoding/json"
	"path/filepath"
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
	var payload map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(resp.Stdout)), &payload); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, resp.Stdout)
	}
	for _, key := range []string{"root", "rows", "summary"} {
		if payload[key] == nil {
			t.Fatalf("json missing %q key: %s", key, resp.Stdout)
		}
	}
	var root string
	if err := json.Unmarshal(payload["root"], &root); err != nil || root != req.FixtureDir {
		t.Fatalf("json root: expected %q, got %q err=%v", req.FixtureDir, root, err)
	}
	var rows []map[string]any
	if err := json.Unmarshal(payload["rows"], &rows); err != nil {
		t.Fatalf("json rows: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected one row, got %d: %#v", len(rows), rows)
	}
	wantSub := filepath.Join(req.FixtureDir, "sub")
	if got, _ := rows[0]["path"].(string); got != wantSub {
		t.Fatalf("row path: expected %q, got %q", wantSub, got)
	}
	var summary map[string]any
	if err := json.Unmarshal(payload["summary"], &summary); err != nil {
		t.Fatalf("json summary: %v", err)
	}
	size, ok := summary["size"].(float64)
	if !ok || int64(size) != 2048 {
		t.Fatalf("summary size: expected 2048, got %#v", summary["size"])
	}
	pnpmShared, ok := summary["pnpmSharedSize"].(float64)
	if !ok || int64(pnpmShared) != 0 {
		t.Fatalf("summary pnpmSharedSize: expected 0, got %#v", summary["pnpmSharedSize"])
	}
	if human, _ := summary["pnpmSharedHuman"].(string); human == "" {
		t.Fatalf("summary pnpmSharedHuman: expected non-empty human size, got %#v", summary["pnpmSharedHuman"])
	}
	bunShared, ok := summary["bunSharedSize"].(float64)
	if !ok || int64(bunShared) != 0 {
		t.Fatalf("summary bunSharedSize: expected 0, got %#v", summary["bunSharedSize"])
	}
	if bunHuman, _ := summary["bunSharedHuman"].(string); bunHuman == "" {
		t.Fatalf("summary bunSharedHuman: expected non-empty human size, got %#v", summary["bunSharedHuman"])
	}
	rowPnpm, ok := rows[0]["pnpmSharedSize"].(float64)
	if !ok || int64(rowPnpm) != 0 {
		t.Fatalf("row pnpmSharedSize: expected 0, got %#v", rows[0]["pnpmSharedSize"])
	}
	if rowHuman, _ := rows[0]["pnpmSharedHuman"].(string); rowHuman == "" {
		t.Fatalf("row pnpmSharedHuman: expected non-empty human size, got %#v", rows[0]["pnpmSharedHuman"])
	}
	rowBun, ok := rows[0]["bunSharedSize"].(float64)
	if !ok || int64(rowBun) != 0 {
		t.Fatalf("row bunSharedSize: expected 0, got %#v", rows[0]["bunSharedSize"])
	}
	if rowBunHuman, _ := rows[0]["bunSharedHuman"].(string); rowBunHuman == "" {
		t.Fatalf("row bunSharedHuman: expected non-empty human size, got %#v", rows[0]["bunSharedHuman"])
	}
}
```