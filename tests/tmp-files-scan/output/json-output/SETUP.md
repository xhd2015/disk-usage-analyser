# Scenario

**Leaf**: `--json` emits NDJSON hit objects followed by text summary

## Preconditions

- The scan has one Mach-O binary hit.

## Steps

1. Create `~/Projects/json-app/.git`.
2. Add one Mach-O binary.
3. Run `scan --json`.
4. Parse the first stdout line as JSON and inspect fields.

## Context

- Summary remains the human text summary after all JSON hits.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/json-app")
	writeMachO(t, app, "bin/json-app")
	req.Args = []string{"scan", "--json"}
	return nil
}
```
