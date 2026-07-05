# Scenario

**Feature**: long node_modules paths use prefix truncation, tooltip full path, and copy button

```
path cell -> truncatePathKeepSuffix display text -> data-full-path attr -> Tooltip (full path + Copy) on hover
```

## Preconditions

- `data-testid="node-modules-path-{rowKey}"` span exposes `data-full-path` with the complete path.
- Visible span shows prefix-truncated display text (suffix fully visible) when path exceeds column width.
- Tooltip contains monospace full path and `data-testid="node-modules-path-copy-{rowKey}"` copy button.

```go
func Setup(t *testing.T, req *Request) error {
	if req.ScriptFile == "" {
		req.ScriptFile = "path-truncation-tooltip.js"
	}
	return nil
}
```