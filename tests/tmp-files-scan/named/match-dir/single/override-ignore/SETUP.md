# Scenario

**Leaf**: `--name=node_modules` reports the directory even though it is in `ignoredDirBasenames`

## Preconditions

- A git repository exists with a `node_modules` directory (normally ignored by `ignoredDirBasenames`).
- `--name=node_modules` should override the ignore and report it as a named hit.

## Steps

1. Create `~/Projects/app/.git`.
2. Create `app/node_modules/pkg/file.txt` with "data\n" (5 bytes).
3. Create `app/src/main.go` (a regular file, not a named or binary hit) to confirm other ignored dirs are still skipped.
4. Run `scan --name node_modules`.

## Context

- The name-match logic takes precedence over `ignoredDirBasenames` for entries named in `--name`.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/app")
	mkdir(t, app, "node_modules/pkg")
	writeData(t, app, "node_modules/pkg/file.txt", []byte("data\n"))
	writeText(t, app, "src/main.go", "package "+"main\n")
	req.Args = []string{"scan", "--name", "node_modules"}
	return nil
}
```
