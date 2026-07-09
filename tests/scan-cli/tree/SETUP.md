# Scenario

**Feature**: TreeResult API with min and max-depth options

```
Scan(path, ScanOptions{min, maxDepth}) -> count all bytes -> filter display -> TreeResult
```

## Preconditions

- Leaves set `req.ScanOpts` when non-default min or maxDepth is required.
- Unset `req.ScanOpts` uses text CLI defaults: min `1M`, maxDepth `3`.
- `totalSize` always reflects full recursive bytes at the scan root.

## Context

- Display filtering omits nodes below min; counting includes them in ancestor and root sizes.
- `MaxDepth` limits visible branch depth; `0` means unlimited.
- Result field is `Min` (JSON `min`), not `Threshold`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "scan"
	return nil
}
```
