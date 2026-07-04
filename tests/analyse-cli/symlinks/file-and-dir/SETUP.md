# Scenario

**Leaf**: immediate-child rows for one file symlink and one directory symlink

## Steps

1. Create `target.txt` (4096 B) and `target-dir/inside.txt` (4096 B).
2. Add `link-file` -> `target.txt` and `link-dir` -> `target-dir`.
3. Run `analyse.Analyse`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "target.txt", file4K)
	mkdir(t, req.FixtureDir, "target-dir")
	writeSizedFile(t, req.FixtureDir, "target-dir/inside.txt", file4K)
	symlinkTo(t, req.FixtureDir, "link-file", "target.txt")
	symlinkTo(t, req.FixtureDir, "link-dir", "target-dir")
	return nil
}
```