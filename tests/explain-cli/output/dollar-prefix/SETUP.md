# Scenario

**Leaf**: human RAW COMMANDS and HOW TO PURGE use `$ ` on runnable lines (not on `#` comments)

```
# Human shell-prompt style
explain go-build-cache PATH ->
  HOW TO PURGE Official command: $ go clean -cache
  RAW COMMANDS:
    # disk-usage-analyser
    $ disk-usage-analyser scan …
```

## Steps

1. Build `Caches/go-build` fixture (predictable `$ go clean -cache` + scan).
2. Run human `explain.RunCLI <go-build-dir>` without color flags.

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
