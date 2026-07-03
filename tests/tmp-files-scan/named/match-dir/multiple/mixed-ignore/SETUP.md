# Scenario

**Leaf**: `--name=node_modules` reports `node_modules` but `vendor` (not in --name) is still skipped

## Preconditions

- A git repository contains `node_modules/` (named) and `vendor/` (not named, still ignored).
- `vendor/` contains a Mach-O binary that should NOT be reported.

## Steps

1. Create `~/Projects/app/.git`.
2. Create `app/node_modules/pkg/file.txt` (6 bytes).
3. Create `app/vendor/tool` as a Mach-O binary (104 bytes).
4. Run `scan --name node_modules`.

## Context

- Only basenames in `--name` override `ignoredDirBasenames`; other ignored dirs remain skipped for binary scanning.
- The Mach-O in `vendor` must not appear in binary hits.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/app")
	mkdir(t, app, "node_modules/pkg")
	writeData(t, app, "node_modules/pkg/file.txt", []byte("hello\n"))
	writeMachO(t, app, "vendor/tool")
	req.Args = []string{"scan", "--name", "node_modules"}
	return nil
}
```
