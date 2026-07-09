# Scenario

**Leaf**: text output uses summary block and tree(1) box-drawing with name-then-aligned-size column

## Steps

1. Write `big.txt` (400 bytes) and `small.txt` (100 bytes) in the fixture root.
2. Run `RunCLI --min 1B <fixture>` (text mode).

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "big.txt", 400)
	writeSizedFile(t, req.FixtureDir, "small.txt", 100)
	req.Args = []string{"--min", "1B", req.FixtureDir}
	return nil
}
```
