# Scenario

**Leaf**: `run.Run` dispatches `tmp-files scan` without starting the web server

## Preconditions

- The run package exposes an option-aware entry point for tests.
- A fake `StartServer` hook fails the test if called.

## Steps

1. Create a fixture repository with one Mach-O binary.
2. Call `run.RunWithOptions` with `tmp-files scan --root <home>`.
3. Verify scan output appears and the server hook is not invoked.

## Context

- This guards against the current web-app default path swallowing `tmp-files` as an initial directory.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "cli-dispatch"
	app := repo(t, req.HomeDir, "Projects/app")
	writeMachO(t, app, "bin/app")
	req.Args = []string{"tmp-files", "scan", "--root", req.HomeDir}
	return nil
}
```
