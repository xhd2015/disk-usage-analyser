# Scenario

**Leaf**: explain on a file inside `.codex` prefers `codex-home` parent context

## Steps

1. Build Codex home fixture (`{fixture}/.codex`).
2. Run `explain.RunCLI` with the absolute path to `config.toml` under `.codex`.

```go
func Setup(t *testing.T, req *Request) error {
	_, configPath := writeCodexHomeFixture(t, req.FixtureDir)
	req.TargetPath = configPath
	req.Args = []string{configPath}
	return nil
}
```
