# Scenario

**Leaf**: root and one subdirectory with regular files, no links

## Steps

1. Write `root.txt` (1024 bytes) in fixture root.
2. Write `sub/child.txt` (2048 bytes) in a subdirectory.
3. Run `analyse.Analyse`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "root.txt", 1024)
	mkdir(t, req.FixtureDir, "sub")
	writeSizedFile(t, req.FixtureDir, "sub/child.txt", 2048)
	return nil
}
```