# Scenario

**Feature**: Codex CLI home reclaim kind (`codex-home`) for `~/.codex`

```
# Auto-detect: dir named .codex + signatures (config.toml / auth.json / logs_*.sqlite / sessions)
explain "…/.codex" -> kind=codex-home, ContextRoot=.codex, confidence high
explain "…/.codex/config.toml" -> kind=codex-home (ContextRoot=.codex)

# Force via short pack/alias (PATH optional)
explain --kind codex [SCOPE?]
  scope = PATH | CLIOptions.HomeDir | os.UserHomeDir()
  measure {scope}/.codex when scope is home-like; if scope is already .codex use it

# Roles (top-level): app-logs-db, session-logs, cache, tmp, config, …
# Product shape A: when logs_*.sqlite present → LOGS DB section (ROWS + SAMPLE last 3)
# SAFE TO RECLAIM: logs/sessions/cache/tmp reclaimable; config/auth never usually-safe
# HOW TO PURGE: scan first; safe logs_2.sqlite reclaim (quit + mv/sqlite3); sessions/cache/tmp;
#   never rm -rf; never purge auth/config/state_5 as part of logs reclaim
```

## Preconditions

- Fixture from `writeCodexHomeFixture`: `{parent}/.codex` with sessions, cache, .tmp,
  config.toml, auth.json, and `logs_2.sqlite` (5 rows via `sqlite3` CLI when available).
- Non-DB content payload sum is `codexHomeNonDBContentBytes` (398).
- Detection runs before `generic-dir` / `generic-file`.
- Tests must never use the real user home for fixtures; use `req.FixtureDir` / `req.HomeDir`.
- Leaves that need the sqlite fixture **Skip** if `sqlite3` CLI is missing.

## Context

- Kind id in output/JSON is **`codex-home`**; CLI force alias is **`codex`**.
- Breakdown assigns roles (app-logs-db, session-logs, cache, tmp, config, …).
- LOGS DB human section always shown when logs db present (shape A).
- SAFE TO RECLAIM must not treat auth/config as usually-safe purge.
- HOW TO PURGE is CLI-first (`disk-usage-analyser scan` inspect) and documents **safe
  `logs_2.sqlite` reclaim**: quit Codex; `mv` backup of db + wal/shm **or**
  `sqlite3 … "DELETE FROM logs; VACUUM;"`; diagnostic-only; not state_5/auth/config;
  recreates / may regrow; no live truncate; **never** `rm -rf`.


```go
func Setup(t *testing.T, req *Request) error {
	// Mark mode for codex-home leaves; concrete TargetPath / --kind is set by child leaves.
	req.Mode = "cli"
	return nil
}
```
