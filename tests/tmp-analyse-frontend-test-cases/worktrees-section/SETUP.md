# Scenario

**Feature**: Git Worktrees subsection tree UI

```
worktrees-scan-btn -> SSE repo events (main size) -> linked worktree children with sizes
```

## Preconditions

- Worktree section is read-only (no delete in v1).
- Repo parent row carries main checkout size from `repo` SSE event.
- Only linked worktrees appear as child rows; no `worktree-main-badge`.
- Tree content is left-aligned within the card.

## Steps

1. Click worktrees scan button.
2. Wait for scan completion.
3. Verify tree rows and sizes.

```go
func Setup(t *testing.T, req *Request) error {
	_ = req.ScriptFile
	return nil
}
```