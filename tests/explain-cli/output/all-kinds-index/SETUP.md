# Scenario

**Leaf**: `explain --all-kinds` with `HomeDir` inject prints multi-kind INDEX (all 5 cli kinds;
present size DESC; present/missing status)

## Steps

1. Build multi-pack home fixture under `req.FixtureDir` (`.codex`, `.grok`, Android SDK, iTerm2; no Xcode).
2. Set `req.HomeDir` to the fixture root (harness injects `CLIOptions.HomeDir`).
3. Run `explain.RunCLI --all-kinds` with **no PATH** argument.

```go
func Setup(t *testing.T, req *Request) error {
	home := writeAllKindsHomeFixture(t, req.FixtureDir)
	req.HomeDir = home
	req.TargetPath = home // scope for asserts
	req.Args = []string{"--all-kinds"}
	return nil
}
```
