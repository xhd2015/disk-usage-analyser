# Scenario

**Leaf**: `--json` emits one JSON object matching TreeResult

## Steps

1. Write `big.txt` (400 bytes) and `small.txt` (100 bytes) in the fixture root.
2. Run `RunCLI --json <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "big.txt", 400)
	writeSizedFile(t, req.FixtureDir, "small.txt", 100)
	req.Args = []string{"--json", "--threshold", "1B", req.FixtureDir}
	return nil
}
```