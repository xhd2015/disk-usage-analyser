# Scenario

**Leaf**: `explain --all-kinds` on an empty home (no pack roots) exits 0 with all kinds missing

## Steps

1. Leave `req.FixtureDir` empty (no pack fixtures).
2. Set `req.HomeDir` to the empty fixture root.
3. Run `explain.RunCLI --all-kinds` with no PATH.

```go
func Setup(t *testing.T, req *Request) error {
	// FixtureDir is already an empty directory from root Setup.
	req.HomeDir = req.FixtureDir
	req.TargetPath = req.FixtureDir
	req.Args = []string{"--all-kinds"}
	return nil
}
```
