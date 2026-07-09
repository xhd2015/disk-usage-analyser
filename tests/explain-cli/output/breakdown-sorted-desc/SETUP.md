# Scenario

**Leaf**: human BREAKDOWN rows are sorted by size descending (largest first)

```
# Same AVD fixture: 400B userdata, 200B sdcard, 100B snapshots, 32B config
explain PATH(android-avd) ->
  BREAKDOWN rows name order:
    userdata-qemu.img.qcow2  (400)
    sdcard.img               (200)
    snapshots                (100)
    config.ini               (32)
```

## Steps

1. Build AVD fixture with exact unequal sizes.
2. Run `explain.RunCLI <avd-dir>`.

```go
func Setup(t *testing.T, req *Request) error {
	avdDir, _ := writeAVDFixture(t, req.FixtureDir)
	req.TargetPath = avdDir
	req.Args = []string{avdDir}
	return nil
}
```
