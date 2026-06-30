# Scenario

**Feature**: Repository Scans section group on tmp-analyse page

```
Developer Tools -> Repository Scans -> Git Worktrees + Binary files subsections
```

## Preconditions

- Repository Scans group appears after Developer Tools.
- Each subsection has independent Start/Stop scan buttons.

## Steps

1. Navigate to `/tmp-analyse`.
2. Verify section structure via playwright-debug.

```go
func Setup(t *testing.T, req *Request) error {
	_ = req.ScriptFile
	return nil
}
```