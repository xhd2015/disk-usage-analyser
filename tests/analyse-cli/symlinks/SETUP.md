# Scenario

**Feature**: symlink counting without following targets

```
analyse -> encounter symlink -> classify file vs dir -> count Ff+Dd -> do not add target bytes
```

## Preconditions

- Symlinks are created with `os.Symlink`.
- File vs directory classification uses `Readlink` plus `Stat` on the target.

## Context

- Broken symlinks count as file symlinks.
- Target file and directory content must not inflate apparent `size`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "analyse"
	return nil
}
```