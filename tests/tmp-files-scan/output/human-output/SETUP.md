# Scenario

**Leaf**: default human output prints hit lines plus summary

## Preconditions

- Human output is the default when `--json` is omitted.

## Steps

1. Create `~/Projects/human-app/.git`.
2. Add one Mach-O binary.
3. Run `scan`.

## Context

- The expected line shape is `size  kind  path  (repo: repoPath)`.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/human-app")
	writeMachO(t, app, "bin/human-app")
	req.Args = []string{"scan"}
	return nil
}
```
