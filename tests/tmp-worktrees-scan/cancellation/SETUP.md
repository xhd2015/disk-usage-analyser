# Scenario

**Decision**: scan cancellation on client disconnect

```
cancellation -> client closes SSE connection -> scan aborts
```

## Preconditions

- Handler respects `r.Context().Done()` when the client closes the response body.

## Steps

1. Start a worktree scan over httptest.
2. Read until the first `worktree` event.
3. Cancel the request context and close the body.

## Context

- The server must not hang after disconnect; partial results are acceptable.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "worktrees-disconnect"
	return nil
}
```