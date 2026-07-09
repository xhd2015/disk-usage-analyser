# Scenario

**Leaf**: `--color=always` colors the base command green on human command lines

```
# Force color even on non-TTY test buffers
explain --color=always PATH(go-build-cache) ->
  $ <green>go</green> clean -cache
  $ <green>disk-usage-analyser</green> scan …
```

## Steps

1. Build `Caches/go-build` fixture.
2. Run `explain.RunCLI --color=always <go-build-dir>`.

```go
func Setup(t *testing.T, req *Request) error {
	mkdir(t, req.FixtureDir, "Caches/go-build")
	writeSizedFile(t, req.FixtureDir, "Caches/go-build/aa/entry1", 100)
	writeSizedFile(t, req.FixtureDir, "Caches/go-build/bb/entry2", 200)
	target := mkdir(t, req.FixtureDir, "Caches/go-build")
	req.TargetPath = target
	// --color=always forces green base-command ANSI even when stdout is a buffer (non-TTY).
	req.Args = []string{"--color=always", target}
	return nil
}
```
