# Scenario

**Feature**: fundamental directory walk without links

```
analyse -> walk root -> count unique inode bytes -> emit subdir rows + summary
```

## Preconditions

- Fixture trees contain only regular files and directories.
- No symlinks or hard links in this branch.

## Context

- Empty directory produces zero-byte summary with no subdirectory rows.
- Regular files contribute `st_size` once per inode to ancestor subdirectory totals.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "analyse"
	return nil
}
```