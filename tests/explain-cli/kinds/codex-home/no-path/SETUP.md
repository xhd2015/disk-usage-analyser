# Scenario

**Leaf**: `explain --kind codex` with no PATH uses injected `CLIOptions.HomeDir` as scope → `{home}/.codex`

## Steps

1. Build Codex home fixture under `req.FixtureDir/.codex` (fixture root acts as fake home).
2. Set `req.HomeDir` to the fixture root (harness injects `CLIOptions.HomeDir`).
3. Run `explain.RunCLI --kind codex` with **no PATH** argument.

```go
func Setup(t *testing.T, req *Request) error {
	codexDir, _ := writeCodexHomeFixture(t, req.FixtureDir)
	req.TargetPath = codexDir
	req.HomeDir = req.FixtureDir
	req.Args = []string{"--kind", "codex"}
	return nil
}
```
