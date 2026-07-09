# Scenario

**Leaf**: `explain --kind grok` with no PATH uses injected `CLIOptions.HomeDir` as scope → `{home}/.grok`

## Steps

1. Build Grok home fixture under `req.FixtureDir/.grok` (fixture root acts as fake home).
2. Set `req.HomeDir` to the fixture root (harness injects `CLIOptions.HomeDir`).
3. Run `explain.RunCLI --kind grok` with **no PATH** argument.

```go
func Setup(t *testing.T, req *Request) error {
	grokDir, _ := writeGrokHomeFixture(t, req.FixtureDir)
	req.TargetPath = grokDir
	req.HomeDir = req.FixtureDir
	req.Args = []string{"--kind", "grok"}
	return nil
}
```
