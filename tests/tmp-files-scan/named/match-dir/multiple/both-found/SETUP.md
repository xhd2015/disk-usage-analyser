# Scenario

**Leaf**: `--name=node_modules --name=vendor` finds both named dirs in the same repo

## Preconditions

- A git repository contains both a `node_modules` and a `vendor` directory, each with files.

## Steps

1. Create `~/Projects/app/.git`.
2. Create `app/node_modules/pkg/a.txt` (6 bytes) and `app/vendor/tool/b.txt` (5 bytes).
3. Run `scan --name node_modules --name vendor`.

## Context

- Both names are reported as separate named hits in a single scan.
- Both dirs are normally in `ignoredDirBasenames` but are overridden by `--name`.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/app")
	mkdir(t, app, "node_modules/pkg")
	mkdir(t, app, "vendor/tool")
	writeData(t, app, "node_modules/pkg/a.txt", []byte("hello\n"))
	writeData(t, app, "vendor/tool/b.txt", []byte("vend\n"))
	req.Args = []string{"scan", "--name", "node_modules", "--name", "vendor"}
	return nil
}
```
