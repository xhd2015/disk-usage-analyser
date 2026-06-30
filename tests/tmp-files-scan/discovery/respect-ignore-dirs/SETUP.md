# Scenario

**Leaf**: binaries under ignored directories are skipped

## Preconditions

- The repository contains a root binary and binaries under ignored basenames.

## Steps

1. Create `~/Projects/ignore-app/.git`.
2. Write one Mach-O hit under `bin/keep`.
3. Write Mach-O files under `vendor`, `node_modules`, and `.venv`.
4. Run `scan`.

## Context

- Ignored basenames should be reused from scan_repo defaults.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/ignore-app")
	writeMachO(t, app, "bin/keep")
	writeMachO(t, app, "vendor/tool")
	writeMachO(t, app, "node_modules/pkg/tool")
	writeMachO(t, app, ".venv/bin/tool")
	req.Args = []string{"scan"}
	return nil
}
```
