# Scenario

**Bug**: parent-level `tmp-files -h` / `--help` must print usage (today only empty args and `scan -h` work)

```
# Parent-level help (args after tmp-files token)
tmp-files                 -> same scan usage help -> exit 0
tmp-files -h|--help       -> same scan usage help -> exit 0  # currently broken (unknown command)
tmp-files scan -h|--help  -> same scan usage help -> exit 0  # covered by command/help
tmp-files <unknown>       -> non-zero (unknown command)
```

## Preconditions

- Args are those **after** the `tmp-files` token (as `tmpfiles.RunCLI` receives them).
- Parent help must not require the `scan` subcommand token to show usage.

## Steps

1. Run `tmpfiles.RunCLI` with parent-level help args (`[]`, `-h`, or `--help`).
2. Capture stdout, exit code, and error.

## Context

- Expected help content matches `scan -h` / empty-args help: documents
  `tmp-files scan [OPTIONS]` and scan flags (`--go-binaries`, `--root`, `--max-depth`,
  `--json`, `-v/--verbose`, `-h/--help`).
- Help is successful CLI usage rendering, not a filesystem scan.
- User-facing stdout ends with a trailing `\n`.

```go
func Setup(t *testing.T, req *Request) error {
	// Leaves override Args with parent-level help forms.
	req.Op = "parent-help"
	return nil
}
```
