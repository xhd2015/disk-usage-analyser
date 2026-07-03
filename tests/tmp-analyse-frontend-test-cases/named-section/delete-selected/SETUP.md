# Scenario

**Leaf**: Delete Selected opens confirmation modal, then removes deleted paths from tree

## Steps

1. Set req.ScriptFile to named-delete-selected.js.
2. Scan, select one node_modules hit, confirm delete modal, verify row removed.

## Context

- Confirmation modal is always shown before delete.
- Delete calls `POST /api/tmp-named-delete` with selected paths.
- Named hits are directories; delete uses `os.RemoveAll`.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "named-delete-selected.js"
	return nil
}
```
