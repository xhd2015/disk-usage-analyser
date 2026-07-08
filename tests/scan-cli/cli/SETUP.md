# Scenario

**Feature**: CLI surface for scan subcommand

```
RunCLI(args) -> parse flags -> Scan(path, opts) -> format text tree or JSON TreeResult -> exit code
```

## Preconditions

- CLI leaves set `req.Mode` to `cli` or `dispatch`.
- `RunCLI` receives args **after** the `scan` token (no `scan` prefix in `req.Args`).
- Default path is the process current working directory when PATH is omitted.

## Context

- Human text includes `PATH:`, `TOTAL:`, `THRESHOLD:`, `MAX-DEPTH:` summary, blank line, then tree lines.
- `--json` emits one JSON object matching `TreeResult` (nested `tree`; no `items` key).
- All stdout ends with a trailing blank line after the last content line.
- `run.Run` must dispatch `scan` before the web-server branch.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Mode == "" || req.Mode == "scan" {
		req.Mode = "cli"
	}
	return nil
}
```