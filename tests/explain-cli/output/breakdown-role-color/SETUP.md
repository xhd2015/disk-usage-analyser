# Scenario

**Leaf**: `--color=always` colors ROLE cells and reclaimable `☑` (green reclaimable / yellow caution; `☐` plain)

```
# Force color on non-TTY buffers
explain --color=always PATH(android-avd) ->
  ROLE snapshot   green
  ROLE user-data  yellow
  ROLE sdcard     yellow
  ROLE config     yellow
  RECLAIMABLE ☑ green-wrapped; ☐ not green/yellow-wrapped
  # $ base commands may still be green elsewhere (ok)
```

## Steps

1. Build AVD fixture.
2. Run `explain.RunCLI --color=always <avd-dir>`.

```go
func Setup(t *testing.T, req *Request) error {
	avdDir, _ := writeAVDFixture(t, req.FixtureDir)
	req.TargetPath = avdDir
	req.Args = []string{"--color=always", avdDir}
	return nil
}
```
