# Scenario

**Leaf**: random regular file falls back to generic-file

## Steps

1. Write `blob.dat` (256 bytes) under the fixture root.
2. Run `explain.RunCLI` on that file path.

```go
func Setup(t *testing.T, req *Request) error {
	path := writeSizedFile(t, req.FixtureDir, "blob.dat", 256)
	req.TargetPath = path
	req.Args = []string{path}
	return nil
}
```
