# Scenario

**Leaf**: fixture home with no git repositories emits zero worktree hits

## Preconditions

- Fake home contains only non-git directories.

## Steps

1. Create `~/Documents/notes.txt` without any `.git` directory.
2. Run `worktrees-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	writeFile(t, req.HomeDir, "Documents/notes.txt", []byte("no repos here\n"))
	req.Op = "worktrees-scan"
	return nil
}
```