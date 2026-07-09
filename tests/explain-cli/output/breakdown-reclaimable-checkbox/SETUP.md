# Scenario

**Leaf**: human RECLAIMABLE column uses Unicode `☑` / `☐` by reclaim tier (never `[x]`/`[ ]`/true/false)

```
# AVD roles: snapshot reclaimable; user-data/sdcard/config caution
explain PATH(android-avd) ->
  … snapshot     ☑  …
  … user-data    ☐  …
  … sdcard       ☐  …
  … config       ☐  …
```

## Steps

1. Build AVD fixture (includes snapshots + config + userdata + sdcard).
2. Run `explain.RunCLI <avd-dir>`.

```go
func Setup(t *testing.T, req *Request) error {
	avdDir, _ := writeAVDFixture(t, req.FixtureDir)
	req.TargetPath = avdDir
	req.Args = []string{avdDir}
	return nil
}
```
