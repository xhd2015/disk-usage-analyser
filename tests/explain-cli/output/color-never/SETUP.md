# Scenario

**Leaf**: `--color=never` forces plain human output (no ANSI) even if color would otherwise apply

```
# Force color off
explain --color=never PATH -> human sections with $ prefixes but no ESC/CSI sequences
```

## Steps

1. Build `Caches/go-build` fixture (same shape as color-force for contrast).
2. Run `explain.RunCLI --color=never <go-build-dir>`.

```go
func Setup(t *testing.T, req *Request) error {
	mkdir(t, req.FixtureDir, "Caches/go-build")
	writeSizedFile(t, req.FixtureDir, "Caches/go-build/aa/entry1", 100)
	writeSizedFile(t, req.FixtureDir, "Caches/go-build/bb/entry2", 200)
	target := mkdir(t, req.FixtureDir, "Caches/go-build")
	req.TargetPath = target
	req.Args = []string{"--color=never", target}
	return nil
}
```
