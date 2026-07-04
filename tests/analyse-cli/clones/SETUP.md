# Scenario

**Feature**: APFS clone-family reference visibility

```
analyse -> regular file -> cp -c clones -> track doc_id groups -> shared_clone=Σ(size×ref_count)
```

## Preconditions

- APFS clone fixtures use `cp -c` on darwin only.
- Each test uses 4096-byte files for predictable 4K metrics.

## Context

- Apparent `size` counts each inode once (`st_size` deduped).
- `shared_clone` multiplies `size × ref_count` for APFS `doc_id` groups with `ref_count > 1` spanning >1 inode.
- Hard-link inodes take precedence: they contribute to `shared_hardlink`, not `shared_clone`.
- On non-darwin platforms `shared_clone` is always 0.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "analyse"
	return nil
}
```