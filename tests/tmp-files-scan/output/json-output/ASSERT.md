## Expected

- First line is a valid JSON object for the hit.
- JSON fields include path, size, sizeHuman, kind, typeDesc, repoPath, and repoName.
- Final line is the text summary.

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
	exact := map[string]string{
		"path": "~/Projects/json-app/bin/json-app",
		"kind": "macho",
		"repoPath": "~/Projects/json-app",
		"repoName": "json-app",
	}
	for key, want := range exact {
		if got, _ := hit[key].(string); got != want {
			t.Fatalf("json %s: expected %q, got %#v in %#v", key, want, hit[key], hit)
		}
	}
	if hit["size"] == nil || hit["sizeHuman"] == "" || hit["typeDesc"] == "" {
		t.Fatalf("json hit missing size fields or typeDesc: %#v", hit)
	}
	if !strings.HasPrefix(lines[1], "Found 1 binaries, total ") {
		t.Fatalf("expected text summary as last line, got %q", lines[1])
	}
}
```
