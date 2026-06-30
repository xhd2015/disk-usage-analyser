# Scenario

**Feature**: re-sort when a larger item arrives via SSE

```
# each new repo/worktree/binary event triggers filter then sort
new larger SSE item -> repositoryScansLayout -> moves to top immediately
```

## Preconditions

- Re-sort runs on every state update, not only on scan done.

## Steps

1. Leaf provides initial repos plus an insert payload.

```go
func Setup(t *testing.T, req *Request) error {
	if req.FixtureFile == "" {
		req.FixtureFile = "testdata/fixture.json"
	}
	return nil
}
```