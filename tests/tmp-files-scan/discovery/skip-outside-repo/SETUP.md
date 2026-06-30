# Scenario

**Leaf**: binary outside any git repository is not reported

## Preconditions

- The scan root contains a Mach-O file but no ancestor `.git` directory for that file.

## Steps

1. Write `~/Downloads/tool` as a Mach-O fixture outside any repo.
2. Run `scan`.

## Context

- Binary discovery is scoped to repositories found by `scan_repo.Scan`.

```go
func Setup(t *testing.T, req *Request) error {
	writeMachO(t, req.HomeDir, "Downloads/tool")
	req.Args = []string{"scan"}
	return nil
}
```
