# Scenario

**Leaf**: `run.Run` dispatches `explain` without starting the web server

## Steps

1. Create a small generic directory fixture with `note.txt` (128 bytes).
2. Call `run.RunWithOptions` with `explain <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "note.txt", 128)
	req.Mode = "dispatch"
	req.TargetPath = req.FixtureDir
	req.Args = []string{"explain", req.FixtureDir}
	return nil
}
```
