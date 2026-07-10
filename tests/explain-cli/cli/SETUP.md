# Scenario

**Feature**: CLI surface for explain (help and run dispatch)

```
# Help documents PATH, --kind, --all-kinds, --json, --color
explain -h|--help -> usage text -> exit 0

# Root help lists subcommands, server flags, nested-help pointer
disk-usage-analyser -h
  -> subcommands: analyse | scan | explain | tmp-files
  -> --dev, --component
  -> … <command> --help
  -> no StartServer

# Dispatch before web server
run.RunWithOptions(["explain", PATH]) -> explain.RunCLI (no StartServer)
```

## Preconditions

- CLI leaves set `req.Mode` to `cli` or `dispatch`.
- `RunCLI` receives args after the `explain` token.
- Dispatch leaves pass full argv including `explain` to `run.RunWithOptions`.

## Context

- Help must document PATH (required unless `--kind` or `--all-kinds`), **`--kind`**
  (supported values include `xcode`, `grok`, `android-sdk`, `iterm2`, `codex`), **`--all-kinds`**,
  `--json`, and `--color` (`always|never|auto`).
- Root help must list `explain [PATH]` among subcommands, document `--dev` / `--component`,
  and include a nested-help pointer (`command` + `--help`).
- Successful dispatch must not invoke the web server hook.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Mode == "" || req.Mode == "cli" {
		// leave as cli unless leaf overrides to dispatch
		req.Mode = "cli"
	}
	return nil
}
```
