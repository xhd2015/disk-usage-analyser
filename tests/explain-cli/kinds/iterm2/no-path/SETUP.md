# Scenario

**Leaf**: `explain --kind iterm2` with no PATH uses injected `CLIOptions.HomeDir` as scope →
`{home}/Library/Application Support/iTerm2`

## Steps

1. Build iTerm2 fixture under `req.FixtureDir/Library/Application Support/iTerm2`
   (fixture root acts as fake home).
2. Set `req.HomeDir` to the fixture root (harness injects `CLIOptions.HomeDir`).
3. Run `explain.RunCLI --kind iterm2` with **no PATH** argument.

```go
func Setup(t *testing.T, req *Request) error {
	iterm2Dir, _ := writeITerm2Fixture(t, req.FixtureDir)
	req.TargetPath = iterm2Dir
	req.HomeDir = req.FixtureDir
	req.Args = []string{"--kind", "iterm2"}
	return nil
}
```
