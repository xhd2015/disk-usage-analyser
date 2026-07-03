# Scenario

**Leaf**: single directory matches `--name=node_modules`, reported with correct metadata

## Preconditions

- One git repository exists with a `node_modules` directory containing a small file.

## Steps

1. Create `~/Projects/app/.git`.
2. Create `app/node_modules/pkg/file.txt` with "hello\n".
3. Run `scan --name node_modules`.

## Context

- Simplest happy path for named directory detection.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/app")
	mkdir(t, app, "node_modules/pkg")
	writeData(t, app, "node_modules/pkg/file.txt", []byte("hello\n"))
	req.Args = []string{"scan", "--name", "node_modules"}
	return nil
}
```
