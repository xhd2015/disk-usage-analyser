# Scenario

**Feature**: TreeResult API with threshold and max-depth options

```
Scan(path, ScanOptions{threshold, maxDepth}) -> count all bytes -> filter display -> TreeResult
```

## Preconditions

- Leaves set `req.ScanOpts` when non-default threshold or maxDepth is required.
- Unset `req.ScanOpts` uses text CLI defaults: threshold `1M`, maxDepth `3`.
- `totalSize` always reflects full recursive bytes at the scan root.

## Context

- Display filtering omits nodes below threshold; counting includes them in ancestor and root sizes.
- `MaxDepth` limits visible branch depth; `0` means unlimited.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "scan"
	return nil
}
```