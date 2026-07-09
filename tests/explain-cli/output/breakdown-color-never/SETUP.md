# Scenario

**Leaf**: `--color=never` keeps BREAKDOWN plain (no ANSI) while RECLAIMABLE Unicode checkboxes remain

```
# Force color off
explain --color=never PATH(android-avd) ->
  BREAKDOWN table with ☑/☐ and bare roles; zero ANSI escapes
```

## Steps

1. Build AVD fixture.
2. Run `explain.RunCLI --color=never <avd-dir>`.

```go
func Setup(t *testing.T, req *Request) error {
	avdDir, _ := writeAVDFixture(t, req.FixtureDir)
	req.TargetPath = avdDir
	req.Args = []string{"--color=never", avdDir}
	return nil
}
```
