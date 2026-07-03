# Scenario

**Leaf**: repo has both a `node_modules` directory and a Mach-O binary; both appear in output

## Preconditions

- A git repository contains both a named directory (`node_modules`) and a binary file.

## Steps

1. Create `~/Projects/app/.git`.
2. Create `app/node_modules/pkg/file.txt` with "hello\n" (6 bytes).
3. Write a Mach-O executable at `app/bin/app` (104 bytes).
4. Run `scan --name node_modules`.

## Context

- Named hits and binary hits are additive and interleaved.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/app")
	mkdir(t, app, "node_modules/pkg")
	writeData(t, app, "node_modules/pkg/file.txt", []byte("hello\n"))
	writeMachO(t, app, "bin/app")
	req.Args = []string{"scan", "--name", "node_modules"}
	return nil
}
```
