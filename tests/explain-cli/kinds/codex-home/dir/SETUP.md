# Scenario

**Leaf**: explain on `.codex` directory yields kind `codex-home` with roles, LOGS DB, reclaim tiers, and HOW TO PURGE

## Steps

1. Build Codex home fixture under `req.FixtureDir/.codex` (includes `logs_2.sqlite` with 5 rows).
2. Run `explain.RunCLI <abs-path-to-.codex>`.

```go
func Setup(t *testing.T, req *Request) error {
	codexDir, _ := writeCodexHomeFixture(t, req.FixtureDir)
	req.TargetPath = codexDir
	req.Args = []string{codexDir}
	return nil
}
```
