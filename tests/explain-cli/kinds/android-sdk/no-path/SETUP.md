# Scenario

**Leaf**: `explain --kind android-sdk` with no PATH uses injected `CLIOptions.HomeDir` as scope → `{home}/Library/Android/sdk`

## Steps

1. Build Android SDK fixture under `req.FixtureDir/Library/Android/sdk` (fixture root acts as fake home).
2. Set `req.HomeDir` to the fixture root (harness injects `CLIOptions.HomeDir`).
3. Run `explain.RunCLI --kind android-sdk` with **no PATH** argument.

```go
func Setup(t *testing.T, req *Request) error {
	sdkDir, _ := writeAndroidSDKFixture(t, req.FixtureDir)
	req.TargetPath = sdkDir
	req.HomeDir = req.FixtureDir
	req.Args = []string{"--kind", "android-sdk"}
	return nil
}
```
