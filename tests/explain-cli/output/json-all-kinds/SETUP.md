# Scenario

**Leaf**: `--json --all-kinds` emits AllKindsResult envelope (scope, totalSize, kinds[])

## Steps

1. Build multi-pack home fixture under `req.FixtureDir`.
2. Set `req.HomeDir` to the fixture root.
3. Run `explain.RunCLI --json --all-kinds` with no PATH.

```go
func Setup(t *testing.T, req *Request) error {
	home := writeAllKindsHomeFixture(t, req.FixtureDir)
	req.HomeDir = home
	req.TargetPath = home
	req.Args = []string{"--json", "--all-kinds"}
	return nil
}
```
