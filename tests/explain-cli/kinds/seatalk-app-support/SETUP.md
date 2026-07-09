# Scenario

**Feature**: SeaTalk Application Support reclaim kind (`seatalk-app-support`)

```
# Signals: …/Application Support/SeaTalk (dir or file under tree);
# optional signatures: main_*.sqlite / search_*.sqlite + Chromium-ish Cache / Service Worker
explain "Application Support/SeaTalk" -> kind=seatalk-app-support
explain "…/SeaTalk/main_1.sqlite" -> kind=seatalk-app-support (ContextRoot=SeaTalk)

# Roles: web-cache, web-storage, chat-db, search-index, backup, session, config, diagnostics, runtime, other
# SAFE TO RECLAIM: web caches usually safe; backups conditional; chat-db/search-index caution only
# HOW TO PURGE: osascript quit SeaTalk; reclaim Cache/Service Worker…; separate sqlite-backup/idb-backup
```

## Preconditions

- Fixture from `writeSeaTalkFixture`: `Application Support/SeaTalk` with
  Service Worker, Cache, main_1.sqlite, search_1.sqlite, sqlite-backup, config.json.
- Content payload sum is `seatalkContentBytes` (510).
- Detection runs before `generic-dir` / `generic-file`.

## Context

- Breakdown assigns roles (web-cache, chat-db, search-index, backup, config, …).
- SAFE TO RECLAIM must not treat primary chat/search DBs as usually-safe purge.
- HOW TO PURGE is CLI-first: `$ osascript -e 'quit app "SeaTalk"'`, then cache
  remove basenames, then backup dirs; never `rm -rf`.


```go
func Setup(t *testing.T, req *Request) error {
	// Mark mode for SeaTalk human-explain leaves; concrete TargetPath is set by dir/ or file/.
	req.Mode = "cli"
	return nil
}
```
