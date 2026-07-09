# Scenario

**Leaf**: path ending with `go-build` (GOCACHE-like) detects go-build-cache

## Steps

1. Create `Caches/go-build` with two dummy cache entry files (100 B, 200 B).
2. Run `explain.RunCLI` on the `go-build` directory.

```go
func Setup(t *testing.T, req *Request) error {
	mkdir(t, req.FixtureDir, "Caches/go-build")
	writeSizedFile(t, req.FixtureDir, "Caches/go-build/aa/entry1", 100)
	writeSizedFile(t, req.FixtureDir, "Caches/go-build/bb/entry2", 200)
	target := mkdir(t, req.FixtureDir, "Caches/go-build")
	req.TargetPath = target
	req.Args = []string{target}
	return nil
}
```
