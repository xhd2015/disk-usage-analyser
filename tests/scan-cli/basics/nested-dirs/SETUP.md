# Scenario

**Leaf**: one subdirectory aggregates bytes from nested paths

## Steps

1. Write `subdir/a` (100 bytes).
2. Write `subdir/b/nested` (200 bytes).
3. Run `usagescan.Scan`.

```go
func Setup(t *testing.T, req *Request) error {
	mkdir(t, req.FixtureDir, "subdir")
	writeSizedFile(t, req.FixtureDir, "subdir/a", 100)
	mkdir(t, req.FixtureDir, "subdir/b")
	writeSizedFile(t, req.FixtureDir, "subdir/b/nested", 200)
	return nil
}
```