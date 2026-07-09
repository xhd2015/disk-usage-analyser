# Scenario

**Leaf**: `explain --kind xcode <scope>` measures the Xcode multi-root pack under PATH scope

## Steps

1. Build Xcode scope fixture under `req.FixtureDir` (all five roots, exact sizes).
2. Run `explain.RunCLI --kind xcode <fixtureDir>` (PATH = scope).

```go
func Setup(t *testing.T, req *Request) error {
	scope := writeXcodeScopeFixture(t, req.FixtureDir)
	req.TargetPath = scope
	req.Args = []string{"--kind", "xcode", scope}
	return nil
}
```
