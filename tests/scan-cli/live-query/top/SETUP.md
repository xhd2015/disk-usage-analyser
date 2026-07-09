# Scenario

**Leaf**: live `scan --top` emits tree + TOP section

## Steps

1. Write `a.bin` (300 B), `b.bin` (200 B), `c.bin` (50 B) at fixture root.
2. Run `RunCLI --min 1B --top 2 <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "a.bin", 300)
	writeSizedFile(t, req.FixtureDir, "b.bin", 200)
	writeSizedFile(t, req.FixtureDir, "c.bin", 50)
	req.Args = []string{"--min", "1B", "--top", "2", req.FixtureDir}
	return nil
}
```
