# Scenario

**Leaf**: a regular file named `node_modules` (not a directory) is reported as a named hit

## Preconditions

- A git repository contains a plain file (not a directory) whose basename matches `--name`.
- The file is not a binary (it's plain text).

## Steps

1. Create `~/Projects/app/.git`.
2. Create a regular file `app/node_modules` with content "data\n" (5 bytes).
3. Run `scan --name node_modules`.

## Context

- File matches report the direct file size and do not block sibling traversal.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/app")
	writeData(t, app, "node_modules", []byte("data\n"))
	req.Args = []string{"scan", "--name", "node_modules"}
	return nil
}
```
