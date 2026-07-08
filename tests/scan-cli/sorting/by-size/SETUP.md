# Scenario

**Leaf**: items are sorted by size descending with directory tie-break

## Steps

1. Write `large.txt` (500 B).
2. Write `dir-medium/inner` (200 B).
3. Write `dir-tie/inner` (100 B).
4. Write `tie.txt` (100 B).
5. Write `small.txt` (50 B).
6. Run `usagescan.Scan`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "large.txt", 500)
	mkdir(t, req.FixtureDir, "dir-medium")
	writeSizedFile(t, req.FixtureDir, "dir-medium/inner", 200)
	mkdir(t, req.FixtureDir, "dir-tie")
	writeSizedFile(t, req.FixtureDir, "dir-tie/inner", 100)
	writeSizedFile(t, req.FixtureDir, "tie.txt", 100)
	writeSizedFile(t, req.FixtureDir, "small.txt", 50)
	return nil
}
```