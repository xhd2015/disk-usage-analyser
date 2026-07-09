# Scenario

**Leaf**: `explain --kind xcode` with no PATH uses injected `CLIOptions.HomeDir` as scope

## Steps

1. Build Xcode scope fixture under `req.FixtureDir`.
2. Set `req.HomeDir` to the fixture scope (harness injects `CLIOptions.HomeDir`).
3. Run `explain.RunCLI --kind xcode` with **no PATH** argument.

```go
func Setup(t *testing.T, req *Request) error {
	scope := writeXcodeScopeFixture(t, req.FixtureDir)
	req.TargetPath = scope
	req.HomeDir = scope
	req.Args = []string{"--kind", "xcode"}
	return nil
}
```
