# Scenario

**Leaf**: `run.Run` dispatches `scan` without starting the web server

## Steps

1. Write `note.txt` (512 bytes) in the fixture root.
2. Call `run.RunWithOptions` with `scan --min 1B <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "note.txt", 512)
	req.Mode = "dispatch"
	req.Args = []string{"scan", "--min", "1B", req.FixtureDir}
	return nil
}
```
