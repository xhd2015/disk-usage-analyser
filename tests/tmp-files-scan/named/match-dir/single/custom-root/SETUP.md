# Scenario

**Leaf**: `--root` scopes where named directories are searched

## Preconditions

- Two separate subtrees each contain a `node_modules` directory, but only one is under `--root`.

## Steps

1. Create `~/selected/app/.git` with `app/node_modules/file.txt` (5 bytes).
2. Create `~/unselected/app/.git` with `app/node_modules/file.txt` (10 bytes).
3. Run `scan --root <selected-path> --name node_modules`.

## Context

- `--root` limits repo discovery and thus named-hunt scope to the specified subtree.

```go
func Setup(t *testing.T, req *Request) error {
	selected := mkdir(t, req.HomeDir, "selected")
	unselected := mkdir(t, req.HomeDir, "unselected")
	selApp := repo(t, selected, "app")
	mkdir(t, selApp, "node_modules")
	writeData(t, selApp, "node_modules/file.txt", []byte("five\n"))
	unsApp := repo(t, unselected, "app")
	mkdir(t, unsApp, "node_modules")
	writeData(t, unsApp, "node_modules/file.txt", []byte("longlong\n"))
	req.Args = []string{"scan", "--root", selected, "--name", "node_modules"}
	return nil
}
```
