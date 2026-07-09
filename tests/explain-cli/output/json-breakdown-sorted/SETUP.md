# Scenario

**Leaf**: `--json` breakdown is size DESC and each entry has `reclaimable` bool (no ANSI/glyphs)

```
# JSON capture
explain --json PATH(android-avd) ->
  breakdown: [{name, size, role, reclaimable: bool, …}, …]
  # sorted size DESC; reclaimable true for snapshot, false for config/user-data/sdcard
  # no ANSI; reclaimable is boolean not "☑"/"☐" or "[x]"/"[ ]"
```

## Steps

1. Build AVD fixture.
2. Run `explain.RunCLI --json <avd-dir>`.

```go
func Setup(t *testing.T, req *Request) error {
	avdDir, _ := writeAVDFixture(t, req.FixtureDir)
	req.TargetPath = avdDir
	req.Args = []string{"--json", avdDir}
	return nil
}
```
