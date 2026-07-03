# Scenario

**Leaf**: nested `node_modules` directories produce separate hits with correct size exclusion

## Preconditions

- A git repository contains `app/node_modules/` (with a 6-byte file) and `app/node_modules/pkg/node_modules/` (with a 7-byte file).
- The outer hit must NOT include the inner hit's size.

## Steps

1. Create `~/Projects/app/.git`.
2. Create `app/node_modules/outer.txt` (6 bytes).
3. Create `app/node_modules/pkg/node_modules/inner/lib.txt` (7 bytes).
4. Run `scan --name node_modules`.

## Context

- `computeDirSize` skips subdirectories whose basename is in `names`, preventing double-counting.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/app")
	mkdir(t, app, "node_modules/pkg/node_modules/inner")
	writeData(t, app, "node_modules/outer.txt", []byte("outer\n"))
	writeData(t, app, "node_modules/pkg/node_modules/inner/lib.txt", []byte("nested\n"))
	req.Args = []string{"scan", "--name", "node_modules"}
	return nil
}
```
