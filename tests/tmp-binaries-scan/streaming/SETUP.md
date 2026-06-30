# Scenario

**Decision**: SSE event ordering for binary scan

```
streaming -> binary before done
```

## Preconditions

- Handler streams each binary hit immediately over SSE.

## Steps

1. Build a fixture with at least one binary hit.
2. Run `binaries-sse-order`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "binaries-sse-order"
	return nil
}
```
