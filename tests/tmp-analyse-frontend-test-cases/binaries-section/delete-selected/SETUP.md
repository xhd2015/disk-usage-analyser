# Scenario

**Leaf**: Delete Selected opens confirmation modal, then removes deleted paths from tree

## Steps

1. Set req.ScriptFile to binaries-delete-selected.js.
2. Scan, select one binary, confirm delete modal, verify row removed.

## Context

- Confirmation modal is always shown before delete (confirmed default).

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "binaries-delete-selected.js"
	return nil
}
```