# Scenario

**Leaf**: `explain --all-kinds <SCOPE>` prints mini-explain detail sections for each present pack

## Steps

1. Build multi-pack home fixture under `req.FixtureDir`.
2. Run `explain.RunCLI --all-kinds <fixtureDir>` (PATH = home-like scope).

```go
func Setup(t *testing.T, req *Request) error {
	home := writeAllKindsHomeFixture(t, req.FixtureDir)
	req.TargetPath = home
	req.Args = []string{"--all-kinds", home}
	return nil
}
```
