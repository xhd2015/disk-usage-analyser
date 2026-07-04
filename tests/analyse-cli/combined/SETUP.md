# Scenario

**Feature**: integrated tree with files, symlinks, and hard links

```
analyse -> mixed entries -> all metric columns populated on summary
```

## Preconditions

- Fixture combines regular files, symlinks, and hard links in one subtree.

## Context

- Validates that symlink and hardlink accounting compose without double-counting inode bytes.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "analyse"
	return nil
}
```