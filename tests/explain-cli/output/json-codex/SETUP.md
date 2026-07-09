# Scenario

**Leaf**: `--json` emits plain Explanation JSON for Codex CLI home (kind, roles, reclaimable, logsDb, howToPurge; no ANSI, no `$`)

## Steps

1. Build Codex home fixture (`{fixture}/.codex` with `logs_2.sqlite` 5 rows).
2. Run `explain.RunCLI --json <abs-path-to-.codex>`.

```go
func Setup(t *testing.T, req *Request) error {
	codexDir, _ := writeCodexHomeFixture(t, req.FixtureDir)
	req.TargetPath = codexDir
	req.Args = []string{"--json", codexDir}
	return nil
}
```
