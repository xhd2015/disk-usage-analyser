# Scenario

**Feature**: longer path display budget before ellipsis truncation

```
TruncatedPath -> PATH_VISIBLE_CHAR_LIMIT -> truncatePathKeepSuffix(path, limit)
```

## Preconditions

- `PATH_VISIBLE_CHAR_LIMIT` is exported from `disk-usage-analyser-react/src/pathDisplay.ts`.
- Limit is **56** (up from 40) so more path suffix is visible in the node_modules table.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Op == "" {
		req.Op = "path-visible-limit"
	}
	return nil
}
```