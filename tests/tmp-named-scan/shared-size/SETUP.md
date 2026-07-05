# Scenario

**Feature**: shared size metrics on named SSE events from per-hit analyse

```
named hit node_modules -> analyse(absPath) -> pnpm_shared + bun_shared -> sharedSize/sharedHuman on SSE named event
```

## Preconditions

- Darwin clone leaves require `cp -c` (APFS) and skip on non-darwin.
- Non-darwin leaf skips on darwin and expects `sharedSize=0`.
- Store and cache paths are isolated via `DISK_USAGE_ANALYSER_PNPM_STORE` and `DISK_USAGE_ANALYSER_BUN_CACHE`.
- Fixture files use 4096-byte payloads for deterministic 4K metrics.

## Context

- Analyse errors or non-darwin platforms yield shared fields `0` / `"0 B"` without failing the SSE stream.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "named-scan"
	req.Name = "node_modules"
	return nil
}
```