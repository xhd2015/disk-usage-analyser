# Scenario

**Feature**: invalid scan / inspect inputs surface as CLI errors

```
RunCLI -> resolve PATH or --inspect FILE -> validate -> non-zero exit on failure
```

## Preconditions

- Missing or invalid paths must not panic.
- `RunCLI` returns a non-zero exit code for path, flag, and inspect load failures.
- `--threshold` is not a valid flag (breaking rename to `--min`).

## Context

- These leaves exercise the CLI surface because exit codes are user-visible.
- Inspect errors: missing file, invalid JSON, combining `--inspect` with a positional PATH.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	return nil
}
```
