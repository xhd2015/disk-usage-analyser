# Scenario

**Leaf**: one 4K inode referenced by three hard-link paths

## Steps

1. Create `a.bin` (4096 B).
2. Hard-link `b.bin` and `c.bin` to `a.bin`.
3. Run `analyse.Analyse`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "a.bin", file4K)
	hardlinkFile(t, req.FixtureDir, "b.bin", "a.bin")
	hardlinkFile(t, req.FixtureDir, "c.bin", "a.bin")
	return nil
}
```