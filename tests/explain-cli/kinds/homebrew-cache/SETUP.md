# Scenario

**Leaf**: `Caches/Homebrew` layout detects homebrew-cache

## Steps

1. Create `Library/Caches/Homebrew` with a dummy bottle download (300 B).
2. Run `explain.RunCLI` on the Homebrew caches directory.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "Library/Caches/Homebrew/downloads/foo--1.0.bottle.tar.gz", 300)
	target := mkdir(t, req.FixtureDir, "Library/Caches/Homebrew")
	req.TargetPath = target
	req.Args = []string{target}
	return nil
}
```
