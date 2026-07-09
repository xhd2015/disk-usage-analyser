# Scenario

**Leaf**: explain on Application Support/SeaTalk directory yields kind seatalk-app-support with roles, reclaim tiers, and HOW TO PURGE

## Steps

1. Build SeaTalk fixture under `req.FixtureDir/Application Support/SeaTalk`.
2. Run `explain.RunCLI <SeaTalk-dir>`.

```go
func Setup(t *testing.T, req *Request) error {
	seatalkDir, _ := writeSeaTalkFixture(t, req.FixtureDir)
	req.TargetPath = seatalkDir
	req.Args = []string{seatalkDir}
	return nil
}
```
