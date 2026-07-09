# Scenario

**Feature**: cross-cutting explain output contracts (JSON, safety, raw cmds, `$`, color,
CLI-first purge, **BREAKDOWN table / reclaimable**, **`--all-kinds` multi-pack report**)

```
# JSON capture of Explanation (plain: no $, no ANSI)
explain PATH --json -> single JSON object (path, kind, totalSize, breakdown, reclaim, howToPurge, rawCommands)
explain --json --kind xcode [PATH] -> same shape for multi-root pack
explain --json …/.grok -> kind=grok-home with role map + reclaimable bools
explain --json …/.codex -> kind=codex-home with role map + reclaimable + logsDb
explain --json …/Library/Android/sdk -> kind=android-sdk with role map + reclaimable bools
explain --json …/Application Support/iTerm2 -> kind=iterm2 with role map + reclaimable bools
breakdown[] includes reclaimable bool; sorted size DESC

# --all-kinds multi-pack (human + JSON envelope)
explain --all-kinds [SCOPE?] -> SCOPE/MODE/TOTAL header + INDEX + present detail
explain --json --all-kinds [SCOPE?] -> {scope, totalSize, kinds:[{kind,cliKind,path,status,totalSize,…}]}
missing pack roots -> status missing, size 0; exit 0

# Safety: never destructive one-liners
SAFE TO RECLAIM / stdout -> must not contain "rm -rf"

# Human: shell prompt + optional color
RAW COMMANDS / HOW TO PURGE Official command -> "$ <cmd>"; # comments without $
--color=always -> green base command; ROLE cells green/yellow by reclaim tier;
                  reclaimable ☑ green-wrapped; ☐ plain
--color=never / non-TTY auto -> no ANSI (glyphs still ☑/☐)

# BREAKDOWN table (all human explain)
BREAKDOWN -> SIZE NAME ROLE RECLAIMABLE [NOTES]; size DESC; ☑/☐ only (never [x]/[ ])

# CLI-first purge recipes
howToPurge.officialCommand -> emulator/avdmanager/sdkmanager/go/npm/brew/xcrun/simctl… ; UI only in Notes
```

## Preconditions

- Output leaves use real fixtures (AVD, go-build-cache, SeaTalk Application Support,
  Grok `.grok` home, Codex `.codex` home with logs sqlite, Android SDK, iTerm2 Application
  Support, Xcode multi-root scope, multi-pack home for `--all-kinds`, or generic-dir
  cache/tmp) so content is non-empty
  (except `all-kinds-all-missing`, which uses an empty home).
- `req.Mode` is `cli`.
- Default buffer I/O is non-TTY: no color unless `--color=always`.

## Context

- These leaves lock product invariants independent of a single kind’s prose wording.
- JSON field names locked here for single Explanation: `path`, `kind`, `totalSize`,
  `breakdown`, `reclaim`, `howToPurge`, `rawCommands` (plus recommended `confidence`,
  `summary`).
- **`--json --all-kinds`** locks AllKindsResult keys: `scope`, `totalSize`, `kinds[]` with
  per-entry `kind`, `cliKind`, `path`, `status`, `totalSize`.
- `breakdown[].reclaimable` is a JSON **boolean** (human uses `☑`/`☐` only; never glyphs
  or strings as the JSON reclaim signal).
- Human formatter owns `$`, ANSI, and BREAKDOWN table layout; JSON stores plain strings.
- `json-seatalk` additionally locks `kind=seatalk-app-support` and role tags
  (`web-cache`, `chat-db`, `search-index`, `backup`, `config`, …).
- `json-xcode` locks `kind=xcode` and pack roles
  (`derived-data`, `simulator`, `device-support`, `archives`, `docs-cache`) with
  reclaimable bools and CLI-first howToPurge (no `rm -rf`).
- `json-grok` locks `kind=grok-home` and roles
  (`installer-cache`, `session-logs`, `cache`, `logs`, `config`) with reclaimable bools
  and CLI-first howToPurge (scan; no `rm -rf`).
- `json-iterm2` locks `kind=iterm2` and roles
  (`python-env`, `python-env-alias`, `logs`, `meta`, `user-config`) with reclaimable
  bools, hardlink/shared-block wording in summary/reclaim, and CLI-first howToPurge
  (scan/du; no `rm -rf`).
- `--all-kinds` leaves: `all-kinds-index`, `all-kinds-detail`, `all-kinds-all-missing`,
  `json-all-kinds`.
- BREAKDOWN table leaves: `breakdown-table-align`, `breakdown-sorted-desc`,
  `json-breakdown-sorted`, `breakdown-reclaimable-checkbox`, `breakdown-role-color`,
  `breakdown-color-never`, `generic-dir-cache-tmp`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	return nil
}
```
