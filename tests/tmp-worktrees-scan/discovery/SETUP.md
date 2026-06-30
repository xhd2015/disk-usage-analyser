# Scenario

**Decision**: worktree discovery and metadata

```
discovery -> main+linked | main only | multi repos | no repos
```

## Preconditions

- Repo discovery scans `~` using scan_repo defaults.
- Worktree listing uses `git worktree list --porcelain`.

## Steps

1. Build an isolated fixture tree per leaf.
2. Run `worktrees-scan` and inspect parsed worktree hits.

## Context

- Main checkout sizing is on the `repo` event; it is not a `worktree` child.
- Linked worktrees have `isMain=false` and non-empty `head` values.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "worktrees-scan"
	return nil
}
```