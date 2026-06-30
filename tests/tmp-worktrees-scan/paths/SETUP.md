# Scenario

**Decision**: path rendering for worktree scan results

```
paths -> home-relative paths use ~/ prefix
```

## Preconditions

- All paths under the fixture home must render with a `~/` prefix in SSE payloads.

## Steps

1. Create a discoverable repo under the fake home.
2. Run `worktrees-scan` and inspect `path` and `repoPath` fields.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "worktrees-scan"
	return nil
}
```