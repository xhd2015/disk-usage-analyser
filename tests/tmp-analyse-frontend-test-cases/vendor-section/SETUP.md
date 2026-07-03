# Scenario

**Feature**: Vendor directory scans with independent scan control

```
vendor-scan-btn -> SSE named events -> repo-grouped tree
```

## Preconditions

- Vendor hits grouped by `repoPath` with size badges.
- Vendor scan is independent from node_modules scan (separate SSE + session).
- Tree content is left-aligned within the card.

## Steps

1. Run vendor scan.
2. Verify tree renders with size, name, path, repo columns.

```go
func Setup(t *testing.T, req *Request) error {
	_ = req.ScriptFile
	return nil
}
```
