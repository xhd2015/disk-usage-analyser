# Scenario

**Leaf**: RAW COMMANDS includes in-tool `disk-usage-analyser scan` for the path

## Steps

1. Build a small generic directory fixture.
2. Run human `explain.RunCLI <fixture-dir>`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "data.bin", 80)
	req.TargetPath = req.FixtureDir
	req.Args = []string{req.FixtureDir}
	return nil
}
```
