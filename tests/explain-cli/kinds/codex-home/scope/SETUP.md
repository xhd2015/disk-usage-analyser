# Scenario

**Leaf**: `explain --kind codex <scope>` measures `{scope}/.codex` under a home-like PATH

## Steps

1. Build Codex home fixture under `req.FixtureDir/.codex` (fixture root = fake home scope).
2. Run `explain.RunCLI --kind codex <fixtureDir>` (PATH = home-like scope).

```go
func Setup(t *testing.T, req *Request) error {
	codexDir, _ := writeCodexHomeFixture(t, req.FixtureDir)
	// TargetPath is the measured .codex; scope arg is the home-like fixture root.
	req.TargetPath = codexDir
	req.Args = []string{"--kind", "codex", req.FixtureDir}
	return nil
}
```
