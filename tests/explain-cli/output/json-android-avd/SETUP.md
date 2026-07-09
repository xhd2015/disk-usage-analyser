# Scenario

**Leaf**: `--json` emits plain Explanation JSON for an Android AVD (no ANSI, no `$` in commands)

## Steps

1. Build AVD fixture (`MediumPhone.avd`).
2. Run `explain.RunCLI --json <avd-dir>`.


```go
func Setup(t *testing.T, req *Request) error {
	avdDir, _ := writeAVDFixture(t, req.FixtureDir)
	req.TargetPath = avdDir
	req.Args = []string{"--json", avdDir}
	return nil
}
```
