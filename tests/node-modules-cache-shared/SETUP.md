# Scenario

**Feature**: node-modules-cache-shared CLI harness

## Preconditions

- Tests call `nmcacheshared.RunCLI` with buffered stdout/stderr.
- Inventory fixture lives at `testdata/inventory.json` (discovered by walking up from cwd).

## Steps

1. Leaf Setup sets `req.Args` (include `inventoryPath(t)` where an input file is required).
2. Run harness invokes `RunCLI`.

```go
func Setup(t *testing.T, req *Request) error {
	_ = inventoryPath(t)
	return nil
}
```