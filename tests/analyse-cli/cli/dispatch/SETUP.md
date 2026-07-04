# Scenario

**Leaf**: `run.Run` dispatches `analyse` without starting the web server

## Steps

1. Create fixture with `note.txt` (512 B).
2. Call `run.RunWithOptions` with `analyse <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "note.txt", 512)
	req.Mode = "dispatch"
	req.Args = []string{"analyse", req.FixtureDir}
	return nil
}
```