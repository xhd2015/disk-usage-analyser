# Scenario

**Leaf**: `--json --kind xcode <scope>` emits plain Explanation JSON for the Xcode multi-root pack

## Steps

1. Build Xcode scope fixture (all five roots).
2. Run `explain.RunCLI --json --kind xcode <scope>`.

```go
func Setup(t *testing.T, req *Request) error {
	scope := writeXcodeScopeFixture(t, req.FixtureDir)
	req.TargetPath = scope
	req.Args = []string{"--json", "--kind", "xcode", scope}
	return nil
}
```
