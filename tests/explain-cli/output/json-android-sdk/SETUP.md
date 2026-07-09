# Scenario

**Leaf**: `--json` emits plain Explanation JSON for Android SDK (kind, roles, reclaimable, howToPurge; no ANSI, no `$`)

## Steps

1. Build Android SDK fixture (`{fixture}/Library/Android/sdk`).
2. Run `explain.RunCLI --json <abs-path-to-sdk>`.

```go
func Setup(t *testing.T, req *Request) error {
	sdkDir, _ := writeAndroidSDKFixture(t, req.FixtureDir)
	req.TargetPath = sdkDir
	req.Args = []string{"--json", sdkDir}
	return nil
}
```
