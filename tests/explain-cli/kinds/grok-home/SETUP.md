# Scenario

**Feature**: Grok CLI home reclaim kind (`grok-home`) for `~/.grok`

```
# Auto-detect: dir named .grok + signatures (config.toml / auth.json / sessions / downloads)
explain "…/.grok" -> kind=grok-home, ContextRoot=.grok, confidence high
explain "…/.grok/config.toml" -> kind=grok-home (ContextRoot=.grok)

# Force via short pack/alias (PATH optional)
explain --kind grok [SCOPE?]
  scope = PATH | CLIOptions.HomeDir | os.UserHomeDir()
  measure {scope}/.grok when scope is home-like; if scope is already .grok use it

# Roles (top-level children): installer-cache, session-logs, cache, logs, tmp,
# project-meta, app-data, config, runtime-db, …
# SAFE TO RECLAIM: downloads/caches/logs usually reclaimable; sessions caution;
#   config/auth never usually-safe
# HOW TO PURGE: scan first; reclaim downloads *.tmp / old installers; optional sessions;
#   marketplace/logs/models_cache; never rm -rf; never purge auth/config as usually-safe
```

## Preconditions

- Fixture from `writeGrokHomeFixture`: `{parent}/.grok` with downloads, sessions,
  marketplace-cache, logs, config.toml, auth.json.
- Content payload sum is `grokHomeContentBytes` (798).
- Detection runs before `generic-dir` / `generic-file`.
- Tests must never use the real user home for fixtures; use `req.FixtureDir` / `req.HomeDir`.

## Context

- Kind id in output/JSON is **`grok-home`**; CLI force alias is **`grok`**.
- Breakdown assigns roles (installer-cache, session-logs, cache, logs, config, …).
- SAFE TO RECLAIM must not treat auth/config as usually-safe purge.
- HOW TO PURGE is CLI-first (`disk-usage-analyser scan` inspect); never `rm -rf`.


```go
func Setup(t *testing.T, req *Request) error {
	// Mark mode for grok-home leaves; concrete TargetPath / --kind is set by child leaves.
	req.Mode = "cli"
	return nil
}
```
