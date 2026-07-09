# Scenario

**Leaf**: `explain --kind grok <scope>` measures `{scope}/.grok` under a home-like PATH

## Steps

1. Build Grok home fixture under `req.FixtureDir/.grok` (fixture root = fake home scope).
2. Run `explain.RunCLI --kind grok <fixtureDir>` (PATH = home-like scope).

```go
func Setup(t *testing.T, req *Request) error {
	grokDir, _ := writeGrokHomeFixture(t, req.FixtureDir)
	// TargetPath is the measured .grok; scope arg is the home-like fixture root.
	req.TargetPath = grokDir
	req.Args = []string{"--kind", "grok", req.FixtureDir}
	return nil
}
```
