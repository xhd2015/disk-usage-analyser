# Scenario

**Decision**: SSE event ordering for worktree scan

```
streaming -> worktree before done | repo before worktree
```

## Preconditions

- The handler streams SSE events as worktrees are sized.
- At least one git repository with worktrees exists in the fixture home.

## Steps

1. Build a fixture with discoverable repos and worktrees.
2. Run `worktrees-sse-order` and inspect parsed event order.

## Context

- Event format: `event: <type>\ndata: <json>\n\n`
- Expected event types: `repo`, `worktree`, `summary`, `done`

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "worktrees-sse-order"
	return nil
}
```