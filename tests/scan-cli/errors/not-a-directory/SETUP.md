# Scenario

**Leaf**: regular file as PATH is rejected as not a directory

## Steps

1. Write `not-a-dir.txt` (64 bytes) in the fixture root.
2. Run `RunCLI` with that file path as the positional argument.

```go
func Setup(t *testing.T, req *Request) error {
	filePath := writeSizedFile(t, req.FixtureDir, "not-a-dir.txt", 64)
	req.Args = []string{filePath}
	return nil
}
```