# Scenario

**Decision**: CLI command mode

```
tmp-files -> scan help | scan dispatch
```

## Preconditions

- Command tests exercise argument parsing and top-level dispatch.

## Steps

1. Select either help rendering or run package dispatch.

## Context

- Help output is expected to exit cleanly.
- Dispatch must not fall through to web server startup.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{"scan"}
	}
	return nil
}
```
