# Scenario

**Leaf**: one root file plus one subdirectory with nested content

## Steps

1. Write `root.txt` (50 bytes) in the fixture root.
2. Write `subdir/nested.txt` (150 bytes).
3. Run `usagescan.Scan`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "root.txt", 50)
	mkdir(t, req.FixtureDir, "subdir")
	writeSizedFile(t, req.FixtureDir, "subdir/nested.txt", 150)
	return nil
}
```