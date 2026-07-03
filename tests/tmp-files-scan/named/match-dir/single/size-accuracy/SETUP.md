# Scenario

**Leaf**: recursive size of a `node_modules` directory with known file contents matches expected value

## Preconditions

- A git repository contains a `node_modules` directory with three files of known sizes (100, 200, 300 bytes).

## Steps

1. Create `~/Projects/app/.git`.
2. Create files under `node_modules/`: `a.bin` (100 bytes), `sub/b.bin` (200 bytes), `sub/c.bin` (300 bytes).
3. Run `scan --name node_modules`.
4. Verify the named hit size is exactly 600 bytes.

## Context

- Tests the accuracy of `computeDirSize` across nested subdirectories.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/app")
	mkdir(t, app, "node_modules/sub")
	writeData(t, app, "node_modules/a.bin", make([]byte, 100))
	writeData(t, app, "node_modules/sub/b.bin", make([]byte, 200))
	writeData(t, app, "node_modules/sub/c.bin", make([]byte, 300))
	req.Args = []string{"scan", "--name", "node_modules"}
	return nil
}
```
