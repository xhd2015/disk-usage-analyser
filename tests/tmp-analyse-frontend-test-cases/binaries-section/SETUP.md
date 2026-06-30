# Scenario

**Feature**: Binary files subsection with selection and delete

```
binaries-scan-btn -> SSE binary events -> repo-grouped tree -> checkbox selection -> delete modal
```

## Preconditions

- Binaries grouped by `repoPath` with kind badges.
- Parent checkbox selects all binaries in a repo.
- Delete always shows a confirmation modal before `POST /api/tmp-binaries-delete`.
- Only paths from the current scan session may be deleted.
- Tree content is left-aligned within the card.

## Steps

1. Run binaries scan.
2. Interact with checkboxes and delete flow per leaf.

```go
func Setup(t *testing.T, req *Request) error {
	_ = req.ScriptFile
	return nil
}
```