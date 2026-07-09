# Scenario

**Leaf**: explain on iTerm2 Application Support directory yields kind `iterm2` with roles,
reclaim tiers, hardlink wording, and HOW TO PURGE

## Steps

1. Build iTerm2 fixture under `req.FixtureDir/Library/Application Support/iTerm2`.
2. Run `explain.RunCLI <abs-path-to-iTerm2>`.

```go
func Setup(t *testing.T, req *Request) error {
	iterm2Dir, _ := writeITerm2Fixture(t, req.FixtureDir)
	req.TargetPath = iterm2Dir
	req.Args = []string{iterm2Dir}
	return nil
}
```
