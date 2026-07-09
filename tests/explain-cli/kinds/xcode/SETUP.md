# Scenario

**Feature**: Xcode multi-root reclaim pack via `--kind xcode`

```
# Force pack (PATH = optional scope under which relatives resolve)
explain --kind xcode [PATH?] -> kind=xcode
  scope = PATH | CLIOptions.HomeDir | os.UserHomeDir()
  measure existing:
    Library/Developer/Xcode/DerivedData              → derived-data ☑
    Library/Developer/CoreSimulator/Devices          → simulator ☑
    Library/Developer/Xcode/iOS DeviceSupport        → device-support ☑
    Library/Developer/Xcode/Archives                 → archives ☐
    Library/Developer/Xcode/DocumentationCache       → docs-cache ☑
  omit missing roots; TOTAL = sum of existing measured roots
  HOW TO PURGE: CLI-first (xcrun simctl / disk-usage-analyser scan); never rm -rf
```

## Preconditions

- Fixture from `writeXcodeScopeFixture`: all five roots with deterministic payload sizes
  (`xcodeContentBytes` = 830 file bytes).
- Tests must never use the real user home for fixtures; use `req.FixtureDir` / `req.HomeDir`.
- `req.Mode` is `cli`.
- Missing roots must not appear as empty BREAKDOWN rows.

## Context

- v1 registry accepts kind ids `xcode`, `grok`, `android-sdk`, `iterm2`, and `codex`
  (other ids → errors/unknown-kind).
- Auto-detect of a single Xcode leaf path without `--kind` is out of scope for these leaves.
- Archives are signed builds → non-reclaimable `☐`; other pack roles reclaimable `☑`
  (simulator purge still cautioned in SAFE TO RECLAIM / HOW TO PURGE).
- Size DESC name order for the full fixture: DerivedData, Devices, Archives,
  DeviceSupport, DocumentationCache.


```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	return nil
}
```
