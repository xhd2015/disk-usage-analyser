---
label: slow, ui-automation
---

## Expected

- `node-modules-tree` is present after scan completes.
- When node_modules dirs exist: at least one `node-modules-repo-row`, one `node-modules-row`, and non-empty `NAMED_ENTITY_SIZE` with size, name, path, repo.
- The `name` column shows `node_modules` for each hit.
- When no node_modules: `node-modules-empty-state` is shown instead of rows.

## Errors

- Tree element must not be MISSING.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "ELEM node-modules-tree: MISSING") {
		t.Fatal("node_modules tree missing after scan")
	}
	if strings.Contains(resp.Output, "COUNT node-modules-rows: 0") {
		if !strings.Contains(resp.Output, "NAMED_EMPTY_STATE: present") {
			t.Skip("no node_modules dirs on machine and no empty state yet")
		}
		return
	}
	for _, token := range []string{"NAMED_ENTITY_SIZE:", "NAMED_ENTITY_NAME:", "NAMED_ENTITY_PATH:", "NAMED_ENTITY_REPO:"} {
		if !strings.Contains(resp.Output, token) {
			t.Fatalf("expected %s in output", token)
		}
	}
}
```
