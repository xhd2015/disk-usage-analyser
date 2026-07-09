# Scenario

**Leaf**: explain on an Android AVD directory yields kind android-avd with full human sections

## Steps

1. Build AVD fixture under `req.FixtureDir/MediumPhone.avd`.
2. Run `explain.RunCLI <avd-dir>`.

```go
func Setup(t *testing.T, req *Request) error {
	avdDir, _ := writeAVDFixture(t, req.FixtureDir)
	req.TargetPath = avdDir
	req.Args = []string{avdDir}
	return nil
}
```
