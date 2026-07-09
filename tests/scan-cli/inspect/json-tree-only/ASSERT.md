## Expected

- Exit code 0.
- JSON is ViewResult-shaped: `min`, `maxDepth` (1), `sourceFile`, `tree`.
- Field is `min` not `threshold`.
- Tree has 3 depth-1 children; `big` has no nested children emitted at maxDepth 1.
- Trailing blank line.

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
		t.Fatalf("expected exit 0, got %d (err=%v)\n%s", resp.ExitCode, resp.Err, resp.Stdout)
	}
	line := firstJSONObjectLine(t, resp.Stdout)
	var payload map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, resp.Stdout)
	}
	if payload["threshold"] != nil {
		t.Fatal("must not emit threshold; use min")
	}
	if payload["min"] == nil {
		t.Fatal("ViewResult/TreeResult must include min")
	}
	var min int64
	if err := json.Unmarshal(payload["min"], &min); err != nil || min != 0 {
		t.Fatalf("inspect default min: expected 0, got %d err=%v", min, err)
	}
	var maxDepth int
	if payload["maxDepth"] == nil {
		t.Fatal("missing maxDepth")
	}
	if err := json.Unmarshal(payload["maxDepth"], &maxDepth); err != nil || maxDepth != 1 {
		t.Fatalf("inspect default maxDepth: expected 1, got %d err=%v", maxDepth, err)
	}
	if payload["sourceFile"] == nil {
		t.Fatal("inspect JSON view should include sourceFile")
	}
	var sourceFile string
	_ = json.Unmarshal(payload["sourceFile"], &sourceFile)
	if !strings.Contains(sourceFile, "tree.json") {
		t.Fatalf("sourceFile should reference tree.json, got %q", sourceFile)
	}
	if payload["tree"] == nil {
		t.Fatal("missing tree")
	}
	var tree map[string]any
	if err := json.Unmarshal(payload["tree"], &tree); err != nil {
		t.Fatalf("tree: %v", err)
	}
	children, _ := tree["children"].([]any)
	if len(children) != 3 {
		t.Fatalf("depth-1 tree should have 3 children, got %#v", tree["children"])
	}
	for _, c := range children {
		m, _ := c.(map[string]any)
		if name, _ := m["name"].(string); name == "big" {
			nested, _ := m["children"].([]any)
			if len(nested) > 0 {
				t.Fatalf("maxDepth 1: big should not expose nested children in tree, got %#v", nested)
			}
		}
	}
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
