# Scenario

**Leaf**: one 4K inode referenced by two hard-link paths

## Steps

1. Create `data.bin` (4096 B).
2. Hard-link `alias.bin` to `data.bin`.
3. Run `analyse.Analyse`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "data.bin", file4K)
	hardlinkFile(t, req.FixtureDir, "alias.bin", "data.bin")
	return nil
}
```