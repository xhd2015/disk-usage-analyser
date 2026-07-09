# Scenario

**Feature**: invalid explain inputs surface as CLI errors

```
RunCLI -> require PATH unless --kind or --all-kinds -> resolve PATH on FS -> non-zero on missing/invalid
RunCLI --kind unknown -> non-zero (unknown kind; lists supported: xcode, grok, android-sdk, iterm2, codex)
RunCLI --all-kinds --kind xcode -> non-zero (mutually exclusive)
```

## Preconditions

- PATH is required when both `--kind` and `--all-kinds` are unset; omitting all three must not
  panic and must exit non-zero.
- A PATH that does not exist on disk must exit non-zero with a clear error.
- Unknown `--kind` values exit non-zero with a clear error that lists supported kinds
  (`xcode`, `grok`, `android-sdk`, `iterm2`, `codex`).
- Combining `--all-kinds` with `--kind` exits non-zero with a mutual-exclusion message.
- Error leaves use `req.Mode = "cli"`.

## Context

- Explain accepts both files and directories; “not a directory” is **not** an error here.
- With a valid `--kind` (`xcode`, `grok`, `android-sdk`, `iterm2`, `codex`) or with `--all-kinds`,
  PATH may be omitted (default scope = home); that is not an error leaf.
- Exit codes may be 1 or 2; tests assert non-zero only.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	return nil
}
```
