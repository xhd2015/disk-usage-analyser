# Scenario

**Leaf**: `--json` emits plain Explanation JSON for iTerm2 Application Support
(kind, roles, reclaimable, howToPurge; no ANSI, no `$`)

## Steps

1. Build iTerm2 fixture (`{fixture}/Library/Application Support/iTerm2`).
2. Run `explain.RunCLI --json <abs-path-to-iTerm2>`.

```go
func Setup(t *testing.T, req *Request) error {
	iterm2Dir, _ := writeITerm2Fixture(t, req.FixtureDir)
	req.TargetPath = iterm2Dir
	req.Args = []string{"--json", iterm2Dir}
	return nil
}
```
