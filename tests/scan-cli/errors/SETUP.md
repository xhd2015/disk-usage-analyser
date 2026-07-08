# Scenario

**Feature**: invalid scan paths surface as CLI errors

```
RunCLI -> resolve PATH -> validate exists and is directory -> non-zero exit on failure
```

## Preconditions

- Missing or invalid paths must not panic.
- `RunCLI` returns a non-zero exit code for path validation failures.

## Context

- These leaves exercise the CLI surface because exit codes are user-visible.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	return nil
}
```