# Scenario

**Feature**: Named directory scans (node_modules) with selection and delete

```
node-modules-scan-btn -> SSE named events -> repo-grouped tree -> checkbox selection -> delete modal
```

## Preconditions

- Named hits grouped by `repoPath` with size badges.
- Parent checkbox selects all hits in a repo.
- Delete always shows a confirmation modal before `POST /api/tmp-named-delete`.
- Only paths from the current scan session may be deleted.
- Tree content is left-aligned within the card.

## Steps

1. Run node_modules scan.
2. Interact with checkboxes and delete flow per leaf.

```go
func Setup(t *testing.T, req *Request) error {
	_ = req.ScriptFile
	return nil
}
```
