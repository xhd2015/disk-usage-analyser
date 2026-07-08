# Scenario

**Leaf**: explicit positional PATH scans that directory

## Steps

1. Write `data.bin` (128 bytes) in the fixture root.
2. Run `RunCLI <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "data.bin", 128)
	req.Args = []string{"--threshold", "1B", req.FixtureDir}
	return nil
}
```