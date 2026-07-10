# Scenario

**Decision**: CLI command mode

```
tmp-files -> parent help (-h/--help) | scan help | scan dispatch
```

## Preconditions

- Command tests exercise argument parsing and top-level dispatch.
- Parent-level help (`-h` / `--help` with no `scan` token) is a first-class success path.

## Steps

1. Select parent help, scan help rendering, or run package dispatch.

## Context

- Help output is expected to exit cleanly (including parent-level `-h` / `--help`).
- Dispatch must not fall through to web server startup.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{"scan"}
	}
	return nil
}
```
