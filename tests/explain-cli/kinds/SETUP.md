# Scenario

**Feature**: reclaim-kind detection and `--kind` packs for human `explain` output

```
# Kind registry: path patterns + signature files; first high-confidence wins
explain PATH -> KindDetector -> kind id + confidence
-> measure sizes -> breakdown roles
-> SAFE TO RECLAIM + HOW TO PURGE (CLI-first, $ prefix) + RAW COMMANDS ($ scan + system)

# Multi-root packs / forced kinds (v1: xcode; grok→grok-home; android-sdk; iterm2; codex→codex-home)
explain --kind xcode [PATH?] -> pack measure under scope (PATH or HomeDir)
explain --kind grok [PATH?]  -> {scope}/.grok measure (kind id grok-home)
explain --kind codex [PATH?] -> {scope}/.codex measure (kind id codex-home; LOGS DB when sqlite)
explain --kind android-sdk [PATH?] -> {scope}/Library/Android/sdk (or SDK root)
explain --kind iterm2 [PATH?] -> {scope}/Library/Application Support/iTerm2 (or iTerm2 root)
```

## Preconditions

- All kinds leaves use human (non-JSON) output unless noted.
- `req.Mode` is `cli`.
- Fixtures are tiny exact-byte trees under `req.FixtureDir`.
- Successful explains must include locked human section headers in order.
- Runnable command lines use `$ `; default non-TTY has no ANSI.
- Specialized kinds include CLI purge recipes (`go` / `npm` / `brew` / `emulator` /
  `sdkmanager` / `xcrun` / scan / `du`…).
- Xcode pack leaves use `--kind xcode` (not path auto-detect).
- Grok leaves auto-detect `…/.grok` or force via `--kind grok` (never real `~/.grok`).
- Codex leaves auto-detect `…/.codex` or force via `--kind codex` (never real `~/.codex`).
- Android SDK leaves auto-detect `…/Library/Android/sdk` or force via `--kind android-sdk`
  (never real `~/Library/Android/sdk`).
- iTerm2 leaves auto-detect `…/Library/Application Support/iTerm2` or force via
  `--kind iterm2` (never real `~/Library/Application Support/iTerm2`).

## Context

- Sibling kind directories are mutually exclusive by target shape/kind id (or forced pack).
- Fallbacks: `generic-dir` / `generic-file` when no specialized detector matches.
- File under `*.avd` prefers parent AVD context (`android-avd`) over bare qcow2.
- File under `Application Support/SeaTalk` prefers `seatalk-app-support` (SeaTalk
  ContextRoot) over `generic-file`. Random dirs stay `generic-dir` (not seatalk).
- File under `.grok` prefers `grok-home` (ContextRoot = `.grok`) over `generic-file`.
- File under `.codex` prefers `codex-home` (ContextRoot = `.codex`) over `generic-file`.
- File under Android SDK prefers `android-sdk` (ContextRoot = SDK root) over `generic-file`.
- File under iTerm2 Application Support prefers `iterm2` (ContextRoot = iTerm2 dir)
  over `generic-file`.
- `kinds/xcode/*` force pack composition under a home-like scope (never real `~/Library`).
- `kinds/grok-home/*` auto-detect or `--kind grok` under fixture home (never real `~/.grok`).
- `kinds/codex-home/*` auto-detect or `--kind codex` under fixture home (never real `~/.codex`).
- `kinds/android-sdk/*` auto-detect or `--kind android-sdk` under fixture home
  (never real `~/Library/Android/sdk`).
- `kinds/iterm2/*` auto-detect or `--kind iterm2` under fixture home
  (never real `~/Library/Application Support/iTerm2`).


```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	return nil
}
```
