# Scenario

**Leaf**: repository beyond `--max-depth` is not discovered

## Preconditions

- One repository is shallow enough for max depth.
- Another repository is nested deeper than max depth.

## Steps

1. Create `~/shallow/.git` with a Mach-O binary.
2. Create `~/a/b/deep/.git` with a Mach-O binary.
3. Run `scan --max-depth 1`.

## Context

- `0` means unlimited; non-zero values are passed through to `scan_repo.Scan`.

```go
func Setup(t *testing.T, req *Request) error {
	shallow := repo(t, req.HomeDir, "shallow")
	writeMachO(t, shallow, "bin/shallow")
	deep := repo(t, req.HomeDir, "a/b/deep")
	writeMachO(t, deep, "bin/deep")
	req.Args = []string{"scan", "--max-depth", "1"}
	return nil
}
```
