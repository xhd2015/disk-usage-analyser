# Scenario

**Leaf**: human explain stdout must never contain `rm -rf`

## Steps

1. Build AVD fixture (has SAFE TO RECLAIM content).
2. Run human `explain.RunCLI <avd-dir>` (no `--json`).

```go
func Setup(t *testing.T, req *Request) error {
	avdDir, _ := writeAVDFixture(t, req.FixtureDir)
	req.TargetPath = avdDir
	req.Args = []string{avdDir}
	return nil
}
```
