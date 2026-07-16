# Scenario

**Bug**: `go run` wrapper SIGTERM orphans compiled `--dev` listener

```
# harness starts dev server, tears down, must leave no TCP listener
doctest Run -> go build binary -> exec --dev -> discover port -> cleanup kills listener PID
```

## Preconditions

- Harness will build compiled binary into session cache (`DOCTEST_SESSION_ID`).
- Teardown must kill the process listening on the discovered port, not only `go` parent.

## Steps

1. Descendant leaves set `req.TeardownOnly = true` (no playwright).

## Context

- Fast leaves under this group prove post-teardown port and listener PID are gone.

```go
func Setup(t *testing.T, req *Request) error {
	req.TeardownOnly = true
	return nil
}
```