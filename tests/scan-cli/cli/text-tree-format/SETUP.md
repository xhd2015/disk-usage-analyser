# Scenario

**Leaf**: text output uses summary block and tree(1) box-drawing with name-then-aligned-size column

## Steps

1. Write `big.txt` (400 bytes) and `small.txt` (100 bytes) in the fixture root.
2. Run `RunCLI <fixture>` (text mode, default threshold 1M, maxDepth 3).

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "big.txt", 400)
	writeSizedFile(t, req.FixtureDir, "small.txt", 100)
	req.Args = []string{"--threshold", "1B", req.FixtureDir}
	return nil
}
```