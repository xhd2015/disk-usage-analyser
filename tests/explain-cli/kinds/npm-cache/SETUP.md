# Scenario

**Leaf**: `.npm` with `_cacache` layout detects npm-cache

## Steps

1. Create `.npm/_cacache/content-v2` with a dummy blob (150 B).
2. Run `explain.RunCLI` on the `.npm` directory (or `_cacache` if path-based).

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, ".npm/_cacache/content-v2/sha512/ab/blob", 150)
	// Prefer explaining the .npm root so signals (_cacache under .npm) match the kind table.
	target := mkdir(t, req.FixtureDir, ".npm")
	req.TargetPath = target
	req.Args = []string{target}
	return nil
}
```
