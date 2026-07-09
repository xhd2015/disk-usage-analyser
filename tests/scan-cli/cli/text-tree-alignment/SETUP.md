# Scenario

**Leaf**: nested text tree aligns size column across short and long entry names

## Steps

1. Write `long-name-file.txt` (200 bytes), `a.txt` (100 bytes), and `shortdir/b.txt` (50 bytes).
2. Run `RunCLI --min 1B <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "long-name-file.txt", 200)
	writeSizedFile(t, req.FixtureDir, "a.txt", 100)
	mkdir(t, req.FixtureDir, "shortdir")
	writeSizedFile(t, req.FixtureDir, "shortdir/b.txt", 50)
	req.Args = []string{"--min", "1B", req.FixtureDir}
	return nil
}
```
