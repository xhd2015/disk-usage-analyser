# Scenario

**Leaf**: two repositories produce hits with distinct `repoPath` metadata

## Steps

1. Create `~/Projects/alpha/.git` with Mach-O binary.
2. Create `~/Work/beta/.git` with ELF binary.
3. Run `binaries-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	alpha := repo(t, req.HomeDir, "Projects/alpha")
	writeMachO(t, alpha, "bin/alpha")
	beta := repo(t, req.HomeDir, "Work/beta")
	writeELF(t, beta, "bin/beta")
	req.Op = "binaries-scan"
	return nil
}
```
