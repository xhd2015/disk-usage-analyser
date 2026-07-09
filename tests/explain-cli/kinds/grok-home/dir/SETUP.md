# Scenario

**Leaf**: explain on `.grok` directory yields kind `grok-home` with roles, reclaim tiers, and HOW TO PURGE

## Steps

1. Build Grok home fixture under `req.FixtureDir/.grok`.
2. Run `explain.RunCLI <abs-path-to-.grok>`.

```go
func Setup(t *testing.T, req *Request) error {
	grokDir, _ := writeGrokHomeFixture(t, req.FixtureDir)
	req.TargetPath = grokDir
	req.Args = []string{grokDir}
	return nil
}
```
