# Scenario

**Feature**: output shape and immediate-child row ordering

```
analyse -> ReadDir immediate children -> one row per child (subtree total) -> summary
```

## Preconditions

- Multi-level directory trees with regular files only.

## Context

- Only immediate children of the analysed root produce TSV/JSON rows (no grandchild rows).
- Each child row reports subtree totals; the final summary row repeats root totals with the root path.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "analyse"
	return nil
}
```