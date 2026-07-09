# Scenario

**Leaf**: `explain --kind android-sdk <scope>` measures `{scope}/Library/Android/sdk` under a home-like PATH

## Steps

1. Build Android SDK fixture under `req.FixtureDir/Library/Android/sdk` (fixture root = fake home scope).
2. Run `explain.RunCLI --kind android-sdk <fixtureDir>` (PATH = home-like scope).

```go
func Setup(t *testing.T, req *Request) error {
	sdkDir, _ := writeAndroidSDKFixture(t, req.FixtureDir)
	// TargetPath is the measured SDK root; scope arg is the home-like fixture root.
	req.TargetPath = sdkDir
	req.Args = []string{"--kind", "android-sdk", req.FixtureDir}
	return nil
}
```
