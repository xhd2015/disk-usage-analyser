# Scenario

**Leaf**: two repositories with binaries produce two hits and correct repo metadata

## Preconditions

- The fake home contains two separate git repositories.

## Steps

1. Create `~/Projects/alpha/.git` with a Mach-O binary.
2. Create `~/Work/beta/.git` with an ELF binary.
3. Run `scan`.

## Context

- Result ordering is not significant; metadata correctness is.

```go
func Setup(t *testing.T, req *Request) error {
	alpha := repo(t, req.HomeDir, "Projects/alpha")
	writeMachO(t, alpha, "bin/alpha")
	beta := repo(t, req.HomeDir, "Work/beta")
	writeELF(t, beta, "bin/beta")
	req.Args = []string{"scan"}
	return nil
}
```
