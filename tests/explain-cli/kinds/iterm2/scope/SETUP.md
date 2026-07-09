# Scenario

**Leaf**: `explain --kind iterm2 <scope>` measures `{scope}/Library/Application Support/iTerm2`
under a home-like PATH

## Steps

1. Build iTerm2 fixture under `req.FixtureDir/Library/Application Support/iTerm2`
   (fixture root = fake home scope).
2. Run `explain.RunCLI --kind iterm2 <fixtureDir>` (PATH = home-like scope).

```go
func Setup(t *testing.T, req *Request) error {
	iterm2Dir, _ := writeITerm2Fixture(t, req.FixtureDir)
	// TargetPath is the measured iTerm2 root; scope arg is the home-like fixture root.
	req.TargetPath = iterm2Dir
	req.Args = []string{"--kind", "iterm2", req.FixtureDir}
	return nil
}
```
