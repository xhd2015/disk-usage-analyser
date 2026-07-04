# Scenario

**Feature**: CLI surface for analyse subcommand

```
RunCLI(args) -> parse flags -> Analyse(path) -> format TSV or JSON -> exit code
```

## Preconditions

- CLI leaves set `req.Mode` to `cli` or `dispatch`.
- `RunCLI` receives args **after** the `analyse` token (no `analyse` prefix in `req.Args`).

## Context

- Default path is the process current working directory when DIR is omitted.
- TSV always includes a header row; `--header` remains accepted (no-op).
- TSV `path` column uses `pathfmt.Short`; JSON keeps absolute paths.
- `--json` emits one JSON object.
- `run.Run` must dispatch `analyse` before the web-server branch.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Mode == "" || req.Mode == "analyse" {
		req.Mode = "cli"
	}
	return nil
}
```