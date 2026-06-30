# Scenario

**Leaf**: SSE includes `summary` event with correct binary and repo counts

## Steps

1. Create `~/Projects/sum-a/.git` and `~/Projects/sum-b/.git` each with one Mach-O binary.
2. Run `binaries-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	a := repo(t, req.HomeDir, "Projects/sum-a")
	writeMachO(t, a, "bin/a")
	b := repo(t, req.HomeDir, "Projects/sum-b")
	writeMachO(t, b, "bin/b")
	req.Op = "binaries-scan"
	return nil
}
```
