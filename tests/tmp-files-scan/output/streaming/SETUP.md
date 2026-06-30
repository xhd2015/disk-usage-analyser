# Scenario

**Leaf**: hit output is written before the final summary

## Preconditions

- A recording stdout writer captures every Write call.
- The fixture contains two hits to make buffering visible.

## Steps

1. Create `~/Projects/stream-app/.git`.
2. Add Mach-O and ELF binaries.
3. Run `scan` with the recording writer.
4. Inspect the first Write call.

## Context

- If the implementation buffers all output and writes once at the end, the first write will include the summary and this test fails.

```go
func Setup(t *testing.T, req *Request) error {
	app := repo(t, req.HomeDir, "Projects/stream-app")
	writeMachO(t, app, "bin/first")
	writeELF(t, app, "bin/second")
	req.StreamStdout = &recordingWriter{}
	req.Args = []string{"scan"}
	return nil
}
```
