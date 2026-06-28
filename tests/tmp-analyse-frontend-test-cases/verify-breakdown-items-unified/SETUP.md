# Scenario

**Feature**: breakdown-items-unified UI verification

```
# React page renders tmp-analyse cards and scan UI
User -> /tmp-analyse -> playwright-debug checks data-testid elements
```

## Preconditions
- Backend uses `breakdownItems: [{path, size, fileCount}]` array, removing the primary/extra distinction
- For 1-item locations: shows `card-path` centered (unchanged)
- For 2+-item locations: shows unified `breakdown-items` list with all items and sizes
- The summarized totals (`size`, `fileCount`) are still returned at the top level

## Steps
1. Set req.ScriptFile to "breakdown-items-unified.js"
2. The script navigates to /tmp-analyse and checks the unified breakdown structure
3. Verifies multi-item cards show all entries in breakdown
4. Verifies single-item cards still use centered card-path

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "breakdown-items-unified.js"
	return nil
}
```
