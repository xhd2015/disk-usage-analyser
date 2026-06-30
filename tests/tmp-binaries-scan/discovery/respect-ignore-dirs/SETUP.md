# Scenario

**Leaf**: binaries under `vendor/` are skipped

## Steps

1. Create `~/Projects/ignore-app/.git`.
2. Write Mach-O to `bin/keep` and `vendor/tool`.
3. Run `binaries-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/ignore-app")
	writeMachO(t, app, "bin/keep")
	writeMachO(t, app, "vendor/tool")
	writeMachO(t, app, "node_modules/pkg/tool")
	req.Op = "binaries-scan"
	return nil
}
```
