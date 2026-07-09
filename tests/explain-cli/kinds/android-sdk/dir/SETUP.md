# Scenario

**Leaf**: explain on Android SDK directory yields kind `android-sdk` with roles, reclaim tiers, and HOW TO PURGE

## Steps

1. Build Android SDK fixture under `req.FixtureDir/Library/Android/sdk`.
2. Run `explain.RunCLI <abs-path-to-sdk>`.

```go
func Setup(t *testing.T, req *Request) error {
	sdkDir, _ := writeAndroidSDKFixture(t, req.FixtureDir)
	req.TargetPath = sdkDir
	req.Args = []string{sdkDir}
	return nil
}
```
