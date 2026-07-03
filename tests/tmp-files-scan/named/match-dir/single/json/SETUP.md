# Scenario

**Leaf**: `--json --name=node_modules` emits valid NDJSON with type "named"

## Preconditions

- A git repository contains a `node_modules` directory with files.

## Steps

1. Create `~/Projects/app/.git`.
2. Create `app/node_modules/pkg/file.txt` with "hello\n" (6 bytes).
3. Run `scan --json --name node_modules`.

## Context

- JSON output for named hits uses `"type":"named"` to distinguish from binary hits.
- The text summary follows the JSON lines.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/app")
	mkdir(t, app, "node_modules/pkg")
	writeData(t, app, "node_modules/pkg/file.txt", []byte("hello\n"))
	req.Args = []string{"scan", "--json", "--name", "node_modules"}
	return nil
}
```
