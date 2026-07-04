# Scenario

**Feature**: hard-link reference visibility

```
analyse -> regular file inode -> track nlink -> hardlinks=Σ(nlink-1), shared_hardlink=Σ(size×nlink)
```

## Preconditions

- Hard links are created with `os.Link` within the fixture tree.
- Each test uses 4096-byte files for predictable 4K metrics.

## Context

- Apparent `size` counts each inode once (`st_size` deduped).
- `hardlink_size` sums `size` once per multiply-linked inode.
- `shared_hardlink` multiplies `size × nlink` for inodes with `nlink > 1`.
- `shared_clone` is always 0 in hard-link-only fixtures.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "analyse"
	return nil
}
```