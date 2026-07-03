# Scenario

**Leaf**: root with no git repositories; `--name` has nothing to scan

## Preconditions

- The fake home has no `.git` directories.

## Steps

1. Create `~/notes/readme.txt` (a plain file outside any repo).
2. Run `scan --name node_modules`.

## Context

- Without a git repo, the walk never starts, so named hits are always zero.

```go
func Setup(t *testing.T, req *Request) error {
	writeText(t, req.HomeDir, "notes/readme.txt", "hello\n")
	req.Args = []string{"scan", "--name", "node_modules"}
	return nil
}
```
