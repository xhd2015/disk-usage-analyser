# Scenario

**Leaf**: android-avd HOW TO PURGE is CLI-first (emulator/avdmanager); UI only in Notes

```
# CLI-primary official purge for AVD
explain MediumPhone.avd/ -> HOW TO PURGE Official command: $ emulator … / $ avdmanager …
Notes may include UI (optional): Device Manager …
```

## Steps

1. Build AVD fixture (`MediumPhone.avd`).
2. Run human `explain.RunCLI <avd-dir>` (no `--json`, default color auto/non-TTY).

```go
func Setup(t *testing.T, req *Request) error {
	avdDir, _ := writeAVDFixture(t, req.FixtureDir)
	req.TargetPath = avdDir
	req.Args = []string{avdDir}
	return nil
}
```
