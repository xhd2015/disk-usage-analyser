# Scenario

**Leaf**: `--json` TreeResult metadata keys and no flat items

## Steps

1. Write `sample.bin` (1024 bytes) in the fixture root.
2. Run `RunCLI --json <fixture>` with no extra flags (JSON defaults: threshold 1M, maxDepth 24).

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "sample.bin", 1024)
	req.Args = []string{"--json", req.FixtureDir}
	return nil
}
```