# Scenario

**Feature**: iTerm2 Application Support reclaim kind (`iterm2`) for
`~/Library/Application Support/iTerm2`

```
# Auto-detect: path ends with Application Support/iTerm2 (basename iTerm2 under
# Application Support), OR dir named iTerm2 containing iterm2env and/or version.txt
# and/or iTermServer-*
explain "…/Application Support/iTerm2" -> kind=iterm2, ContextRoot=iTerm2 dir, confidence high
explain "…/Application Support/iTerm2/iterm2env/f" -> kind=iterm2 (ContextRoot=iTerm2 dir)

# Force via kind id (PATH optional)
explain --kind iterm2 [SCOPE?]
  scope = PATH | CLIOptions.HomeDir | os.UserHomeDir()
  measure {scope}/Library/Application Support/iTerm2 when scope is home-like;
  if scope is already an iTerm2 root (signatures), use it

# Roles (top-level children): python-env ☑ (iterm2env), python-env-alias ☑ (iterm2env-*),
# logs ☑ (log*.txt), helper-binary ☐, app-db ☐, user-config ☐, state ☐, meta ☐, …
# SUMMARY/SAFE: hardlink/shared blocks among iterm2env*; logical TOTAL can overstate;
#   confirm with du -sh on parent (~one env of reclaim, not sum of aliases)
# HOW TO PURGE: scan + du; optional logs/envs after quit; never rm -rf;
#   never bulk-delete user-config as usually-safe
```

## Preconditions

- Fixture from `writeITerm2Fixture`: `{parent}/Library/Application Support/iTerm2` with
  `iterm2env`, `iterm2env-3.10`, `log.0.txt`, `version.txt`, `DynamicProfiles`.
- Content payload sum is `iTerm2ContentBytes` (674).
- Detection runs before `generic-dir` / `generic-file`.
- Tests must never use the real user home for fixtures; use `req.FixtureDir` / `req.HomeDir`.

## Context

- Kind id in output/JSON is **`iterm2`** (same as CLI force id).
- Breakdown assigns roles (`python-env`, `python-env-alias`, `logs`, `meta`, `user-config`, …).
- SUMMARY or SAFE TO RECLAIM must document APFS hardlink / shared-block overcount for
  multiple `iterm2env*` trees.
- SAFE TO RECLAIM must not treat DynamicProfiles / Scripts / user-config as usually-safe-only.
- HOW TO PURGE is CLI-first (`disk-usage-analyser scan`, `du`); never `rm -rf`.


```go
func Setup(t *testing.T, req *Request) error {
	// Mark mode for iterm2 leaves; concrete TargetPath / --kind is set by child leaves.
	req.Mode = "cli"
	return nil
}
```
