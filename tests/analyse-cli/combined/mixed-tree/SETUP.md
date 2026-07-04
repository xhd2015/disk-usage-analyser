# Scenario

**Leaf**: immediate children `pkg` (real dir) and `link-pkg` (dir symlink)

## Steps

1. Under `pkg/`: write `readme.txt` (512 B), `data.bin` (4096 B), hard-link `copy.bin`.
2. Add `pkg/link-readme` -> `readme.txt` (file symlink).
3. Add root `link-pkg` -> `pkg` (directory symlink).
4. Run `analyse.Analyse`.

```go
func Setup(t *testing.T, req *Request) error {
	mkdir(t, req.FixtureDir, "pkg")
	writeSizedFile(t, req.FixtureDir, "pkg/readme.txt", 512)
	writeSizedFile(t, req.FixtureDir, "pkg/data.bin", file4K)
	hardlinkFile(t, req.FixtureDir, "pkg/copy.bin", "pkg/data.bin")
	symlinkTo(t, req.FixtureDir, "pkg/link-readme", "pkg/readme.txt")
	symlinkTo(t, req.FixtureDir, "link-pkg", "pkg")
	return nil
}
```