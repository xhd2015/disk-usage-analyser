# Scenario

**Feature**: fallback reclaim kinds `generic-dir` and `generic-file`

```
# No specialized detector matches → generic fallback by file type
explain random-dir/ -> kind=generic-dir
explain random-file.bin -> kind=generic-file
```

## Preconditions

- Fixtures must not look like AVD, caches, or other specialized kinds.
- Human sections and RAW COMMANDS with `scan` still required.

## Context

- Generic fallbacks still measure size and offer conservative reclaim guidance.

```go
func Setup(t *testing.T, req *Request) error {
	// Generic fallback leaves use human CLI mode; TargetPath set by dir/ or file/.
	req.Mode = "cli"
	return nil
}
```

