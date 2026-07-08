# Scenario

**Leaf**: root with two files and no subdirectories

## Steps

1. Write `a.txt` (100 bytes) and `b.txt` (200 bytes) in the fixture root.
2. Run `usagescan.Scan`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "a.txt", 100)
	writeSizedFile(t, req.FixtureDir, "b.txt", 200)
	return nil
}
```