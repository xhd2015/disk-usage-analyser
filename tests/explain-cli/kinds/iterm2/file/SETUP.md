# Scenario

**Leaf**: explain on a file inside iTerm2 Application Support prefers `iterm2` parent context

## Steps

1. Build iTerm2 fixture (`{fixture}/Library/Application Support/iTerm2`).
2. Run `explain.RunCLI` with the absolute path to `iterm2env/f` under the iTerm2 tree.

```go
func Setup(t *testing.T, req *Request) error {
	_, filePath := writeITerm2Fixture(t, req.FixtureDir)
	req.TargetPath = filePath
	req.Args = []string{filePath}
	return nil
}
```
