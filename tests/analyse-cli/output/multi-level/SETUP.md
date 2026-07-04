# Scenario

**Leaf**: nested subdirectories emit immediate-child rows only plus summary

## Steps

1. Create `alpha/x.bin` (1024 B), `alpha/beta/y.bin` (1024 B), `gamma/z.bin` (1024 B).
2. Run `analyse.Analyse`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "alpha/x.bin", 1024)
	mkdir(t, req.FixtureDir, "alpha/beta")
	writeSizedFile(t, req.FixtureDir, "alpha/beta/y.bin", 1024)
	mkdir(t, req.FixtureDir, "gamma")
	writeSizedFile(t, req.FixtureDir, "gamma/z.bin", 1024)
	return nil
}
```