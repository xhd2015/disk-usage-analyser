# Scenario

**Leaf**: explain on a file inside the Android SDK prefers `android-sdk` parent context

## Steps

1. Build Android SDK fixture (`{fixture}/Library/Android/sdk`).
2. Run `explain.RunCLI` with the absolute path to `platform-tools/f` under the SDK.

```go
func Setup(t *testing.T, req *Request) error {
	_, filePath := writeAndroidSDKFixture(t, req.FixtureDir)
	req.TargetPath = filePath
	req.Args = []string{filePath}
	return nil
}
```
