# Scenario

**Leaf**: explain on a file inside `.grok` prefers `grok-home` parent context

## Steps

1. Build Grok home fixture (`{fixture}/.grok`).
2. Run `explain.RunCLI` with the absolute path to `config.toml` under `.grok`.

```go
func Setup(t *testing.T, req *Request) error {
	_, configPath := writeGrokHomeFixture(t, req.FixtureDir)
	req.TargetPath = configPath
	req.Args = []string{configPath}
	return nil
}
```
