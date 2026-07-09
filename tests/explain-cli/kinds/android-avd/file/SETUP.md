# Scenario

**Leaf**: explain on userdata-qemu.img.qcow2 inside an AVD prefers android-avd parent context

## Steps

1. Build AVD fixture.
2. Run `explain.RunCLI` with the absolute path to `userdata-qemu.img.qcow2`.

```go
func Setup(t *testing.T, req *Request) error {
	_, userdata := writeAVDFixture(t, req.FixtureDir)
	req.TargetPath = userdata
	req.Args = []string{userdata}
	return nil
}
```
