# Scenario

**Leaf**: `--json` emits plain Explanation JSON for SeaTalk Application Support (kind, roles, howToPurge; no ANSI, no `$`)

## Steps

1. Build SeaTalk fixture (`Application Support/SeaTalk`).
2. Run `explain.RunCLI --json <SeaTalk-dir>`.


```go
func Setup(t *testing.T, req *Request) error {
	seatalkDir, _ := writeSeaTalkFixture(t, req.FixtureDir)
	req.TargetPath = seatalkDir
	req.Args = []string{"--json", seatalkDir}
	return nil
}
```
