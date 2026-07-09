# Scenario

**Feature**: CLI surface for scan subcommand (live source, capture, help, dispatch)

```
RunCLI(args) -> parse flags -> LiveTreeSource | JSONTreeSource -> View -> text or JSON -> exit code
```

## Preconditions

- CLI leaves set `req.Mode` to `cli` or `dispatch`.
- `RunCLI` receives args **after** the `scan` token (no `scan` prefix in `req.Args`).
- Default path is the process current working directory when PATH is omitted.
- Dispatch leaves pass a full argv including the `scan` token to `run.RunWithOptions`.

## Context

- Human text includes `PATH:`, `TOTAL:`, `MIN:`, `MAX-DEPTH:` summary, blank line, then tree lines.
- Pure capture `--json` (no query extras) emits one JSON object matching `TreeResult` with field **`min`**.
- All stdout ends with a trailing blank line after the last content line.
- `run.Run` must dispatch `scan` (including `--inspect`) before the web-server branch.
- Flag rename: **`--min`** replaces **`--threshold`** (no alias).

```go
func Setup(t *testing.T, req *Request) error {
	if req.Mode == "" || req.Mode == "scan" {
		req.Mode = "cli"
	}
	return nil
}
```
