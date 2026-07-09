# Scenario

**Leaf**: random directory falls back to generic-dir

## Steps

1. Create `project/notes.txt` (64 bytes) under the fixture root (not a known reclaim layout).
2. Run `explain.RunCLI <fixture-dir>`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "project/notes.txt", 64)
	req.TargetPath = req.FixtureDir
	req.Args = []string{req.FixtureDir}
	return nil
}
```
