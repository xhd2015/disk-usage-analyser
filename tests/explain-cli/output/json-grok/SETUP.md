# Scenario

**Leaf**: `--json` emits plain Explanation JSON for Grok CLI home (kind, roles, reclaimable, howToPurge; no ANSI, no `$`)

## Steps

1. Build Grok home fixture (`{fixture}/.grok`).
2. Run `explain.RunCLI --json <abs-path-to-.grok>`.

```go
func Setup(t *testing.T, req *Request) error {
	grokDir, _ := writeGrokHomeFixture(t, req.FixtureDir)
	req.TargetPath = grokDir
	req.Args = []string{"--json", grokDir}
	return nil
}
```
