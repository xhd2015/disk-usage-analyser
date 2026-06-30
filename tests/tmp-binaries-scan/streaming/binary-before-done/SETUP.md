# Scenario

**Leaf**: at least one `binary` SSE event arrives before `done`

## Preconditions

- Repository contains two binaries to make buffering visible.

## Steps

1. Create `~/Projects/stream-app/.git` with Mach-O and ELF binaries.
2. Run `binaries-sse-order`.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/stream-app")
	writeMachO(t, app, "bin/first")
	writeELF(t, app, "bin/second")
	req.Op = "binaries-sse-order"
	return nil
}
```
