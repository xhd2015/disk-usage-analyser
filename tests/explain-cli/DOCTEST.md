# Explain CLI

CLI tests for `disk-usage-analyser explain [PATH] [--kind KIND] [--all-kinds] [--json] [--color=always|never|auto]`:
reclaim-kind detection, **`--kind` packs** (v1: **`xcode`**, **`grok`** → kind id **`grok-home`**,
**`android-sdk`**, **`iterm2`**, **`codex`** → kind id **`codex-home`**), **`--all-kinds`**
multi-pack report under optional scope, semantic size breakdown, **SAFE TO RECLAIM** advice
(never `rm -rf`), **HOW TO PURGE** (CLI-first official purge recipes + what files are removed;
UI only as optional Notes), and **RAW COMMANDS** (in-tool `scan` plus system/ecosystem commands,
informational even if binaries are missing). Human formatter adds a **`$ `** shell-prompt prefix
on runnable command lines and optionally **green** ANSI on the base command token when color is
enabled. **codex-home** with readable `logs_*.sqlite` always emits a **LOGS DB** section
(ROWS + SAMPLE last 3) and JSON `logsDb`.

## Version

0.0.8

# DSN (Domain Specific Notion)

The **explain command** is a CLI subcommand that classifies a target into a **reclaim kind**,
measures size (via a modest walk or kind-specific known basenames; does not re-walk home unless
a multi-root pack scopes under home), and prints human sections or JSON.

**PATH** (directory or file) is **required unless `--kind` or `--all-kinds` is set**. With
**`--kind`** or **`--all-kinds`**, PATH is an optional **scope** root (home-like directory under
which pack-relative paths resolve). When either flag is set and PATH is omitted, default scope
is **`CLIOptions.HomeDir`** when the harness injects it, else the process home
(`os.UserHomeDir()`). Unknown `--kind` values exit non-zero. **`--all-kinds` and `--kind` are
mutually exclusive** (combined → non-zero error). Without `--kind` and without `--all-kinds`,
empty args still fail (PATH required).

**KindDetector registry** matches path patterns and signature files; the first
high-confidence match wins when PATH is given without `--kind`. Covered auto-detect kinds
include **android-avd**, **android-sdk** (Android SDK root under `…/Android/sdk` or signature
dirs: `platform-tools` plus one of `platforms`/`system-images`/`build-tools`/`cmdline-tools`/
`emulator`), **go-build-cache**, **npm-cache**, **homebrew-cache**, **generic-qcow2**,
**seatalk-app-support** (SeaTalk `Application Support/SeaTalk` tree), **grok-home** (Grok /
xAI CLI `…/.grok` with signatures), **codex-home** (OpenAI Codex CLI `…/.codex` with
signatures including `logs_*.sqlite`), **iterm2** (iTerm2 `Application Support/iTerm2` tree
with signatures), and fallbacks **generic-dir** / **generic-file**. A file inside an `*.avd`
tree is preferred as **android-avd** (parent context) over a bare disk-image kind. A file under
SeaTalk Application Support prefers **seatalk-app-support** (ContextRoot = SeaTalk dir) over
**generic-file**. A file under `.grok` prefers **grok-home** (ContextRoot = `.grok` dir) over
**generic-file**. A file under `.codex` prefers **codex-home** (ContextRoot = `.codex` dir)
over **generic-file**. A file under an Android SDK tree prefers **android-sdk** (ContextRoot =
SDK root) over **generic-file**. A file under iTerm2 Application Support prefers **iterm2**
(ContextRoot = iTerm2 dir) over **generic-file**.

**`--kind` multi-root packs / forced kinds** (v1 registry: **`xcode`**, **`grok`**,
**`android-sdk`**, **`iterm2`**, **`codex`**). Unknown `--kind` values exit non-zero and list
supported kinds (`xcode`, `grok`, `android-sdk`, `iterm2`, `codex`).

**`--kind xcode`** forces the Xcode pack id and measures an ordered list of relative roots
under scope. Missing roots are **omitted** from BREAKDOWN (no empty invented rows). Xcode pack
roles / relatives (match frontend + `server/tmp_analyse.go`):

| Role | Relative path under scope | reclaimable |
|------|---------------------------|-------------|
| `derived-data` | `Library/Developer/Xcode/DerivedData` | ☑ true |
| `simulator` | `Library/Developer/CoreSimulator/Devices` | ☑ true (caution: wipe devices) |
| `device-support` | `Library/Developer/Xcode/iOS DeviceSupport` | ☑ true |
| `archives` | `Library/Developer/Xcode/Archives` | ☐ false (signed builds) |
| `docs-cache` | `Library/Developer/Xcode/DocumentationCache` | ☑ true |

**`--kind grok`** forces the Grok CLI home pack/alias: PATH is optional scope (default
`CLIOptions.HomeDir` / process home). Scope is home-like → measure **`{scope}/.grok`**; if
scope is already a `.grok` directory, use it. Kind id in output/JSON is **`grok-home`**
(not the CLI alias). Auto-detect signatures for a directory named `.grok`: at least one of
`config.toml`, `auth.json`, `sessions/`, `downloads/`. Top-level role map (fixture-covered):

| Basename / pattern | Role | reclaimable |
|--------------------|------|-------------|
| `downloads` | `installer-cache` | ☑ true |
| `sessions` | `session-logs` | ☑ true (caution in SAFE/HOW TO PURGE) |
| `marketplace-cache` | `cache` | ☑ true |
| `logs` | `logs` | ☑ true |
| `models_cache.json` | `cache` | ☑ true |
| `upload_queue` | `tmp` | ☑ true |
| `projects` | `project-meta` | ☐ false |
| `skills`, `vendor`, `bundled`, `docs`, `completions` | `app-data` | ☐ false |
| `config.toml`, `auth.json`, `trusted_folders.toml`, `agent_id`, `*.lock` | `config` | ☐ false |
| `worktrees.db`, `active_sessions.json` | `runtime-db` | ☐ false |

HOW TO PURGE for **grok-home** is CLI-first: inspect with `disk-usage-analyser scan` on
`.grok` / downloads / sessions; reclaim installer leftovers / caches; optional old sessions;
**never** `rm -rf`; never mark auth/config as usually-safe.

**`--kind codex`** forces the Codex CLI home pack/alias: PATH is optional scope (default
`CLIOptions.HomeDir` / process home). Scope is home-like → measure **`{scope}/.codex`**; if
scope is already a `.codex` directory, use it. Kind id in output/JSON is **`codex-home`**
(not the CLI alias). Auto-detect signatures for a directory named `.codex`: at least one of
`config.toml`, `auth.json`, `logs_*.sqlite`, `sessions/`, `history.jsonl`. Top-level role map
(fixture-covered):

| Basename / pattern | Role | reclaimable |
|--------------------|------|-------------|
| `logs_*.sqlite` (+ wal/shm) | `app-logs-db` | ☑ true (caution) |
| `sessions` | `session-logs` | ☑ true |
| `.tmp` | `tmp` | ☑ true |
| `cache`, `models_cache.json`, `cloud-*-cache.json` | `cache` | ☑ true |
| `shell_snapshots` | `snapshots` | ☑ true caution |
| `history.jsonl`, `session_index.jsonl` | `history` | ☑ true caution |
| `state_*.sqlite`, `goals_*.sqlite`, `memories_*.sqlite` | `app-state-db` | ☐ false |
| `config.toml`, `auth.json`, `hooks.json`, `installation_id` | `config` | ☐ false |
| `skills`, `plugins`, `vendor_imports`, `prompts`, `rules` | `app-data` | ☐ false |

**Product shape A — LOGS DB**: when a readable `logs_*.sqlite` is present under the measured
`.codex`, human output always includes a **`LOGS DB`** section (after `BREAKDOWN`, before
`SAFE TO RECLAIM`) with PATH, SIZE (file size), **ROWS** = `COUNT(*)` on table `logs`, and
**SAMPLE (last 3, newest first)** via `ORDER BY ts DESC, ts_nanos DESC, id DESC LIMIT 3`
(body truncated ~120–200 chars). JSON Explanation gains optional/required **`logsDb`**:

```
logsDb: { path, size, rows, samples:[{id, ts, level, target, body}] }
```

HOW TO PURGE for **codex-home** is CLI-first: inspect with `disk-usage-analyser scan` on
`.codex` / sessions / logs; reclaim logs/sessions/cache/tmp with caution; **never** `rm -rf`;
never mark auth/config as usually-safe. **Safe `logs_2.sqlite` reclaim**: quit Codex fully
first; prefer `mv` backup of `logs_2.sqlite` + `-wal`/`-shm`, **or** after quit
`sqlite3 … "DELETE FROM logs; VACUUM;"`; notes: diagnostic-only (not `state_5`/auth/config);
Codex recreates the DB; may regrow (e.g. TRACE); do not truncate while running. Sessions
remain a separate caution step. v1 does not require `--logs-follow` / `--logs-stats`.

**`--kind android-sdk`** forces the Android SDK pack/id: PATH is optional scope (default
`CLIOptions.HomeDir` / process home). Scope is home-like → measure
**`{scope}/Library/Android/sdk`**; if scope is already an SDK root (path ends with
`Android/sdk` or signatures present), use it. Kind id in output/JSON is **`android-sdk`**.
Auto-detect without `--kind` uses the same signatures. Top-level role map (fixture-covered):

| Basename / pattern | Role | reclaimable |
|--------------------|------|-------------|
| `system-images` | `system-images` | ☑ true |
| `sources` | `sources` | ☑ true |
| `skins` | `skins` | ☑ true |
| `.temp`, `.downloadIntermediates` | `tmp` | ☑ true |
| `emulator` | `emulator` | ☐ false |
| `build-tools` | `build-tools` | ☐ false |
| `platforms` | `platforms` | ☐ false |
| `platform-tools` | `platform-tools` | ☐ false |
| `cmdline-tools` | `cmdline-tools` | ☐ false |
| `licenses` | `licenses` | ☐ false |
| `.knownPackages` | `meta` | ☐ false |

HOW TO PURGE for **android-sdk** is CLI-first: `$ sdkmanager --list_installed` (or path under
`cmdline-tools/…/bin`), `$ sdkmanager --uninstall "…"` for unused packages (e.g. system-images),
`$ disk-usage-analyser scan` on SDK / system-images; optional Notes for Android Studio SDK
settings UI; **never** `rm -rf`. SAFE TO RECLAIM: usually-safe temp; usually-safe-with-caution
system-images/sources/skins; caution keep platform-tools/build-tools/platforms/emulator bulk.

**`--kind iterm2`** forces the iTerm2 Application Support pack/id: PATH is optional scope
(default `CLIOptions.HomeDir` / process home). Scope is home-like → measure
**`{scope}/Library/Application Support/iTerm2`**; if scope is already an iTerm2 root (path ends
with `Application Support/iTerm2` or signatures present), use it. Kind id in output/JSON is
**`iterm2`**. Auto-detect without `--kind` uses the same signatures:

- Path ends with `Application Support/iTerm2` (basename `iTerm2`, parent `Application Support`),
  **or**
- Directory named `iTerm2` containing `iterm2env` and/or `version.txt` and/or `iTermServer-*`

Top-level role map (fixture-covered):

| Basename / pattern | Role | reclaimable |
|--------------------|------|-------------|
| `iterm2env` (exact) | `python-env` | ☑ true |
| `iterm2env-*` | `python-env-alias` | ☑ true |
| `log*.txt` / name starts with `log` and ends `.txt` | `logs` | ☑ true |
| `iTermServer-*` | `helper-binary` | ☐ false |
| `chatdb.sqlite*` | `app-db` | ☐ false |
| `SavedState` | `state` | ☐ false |
| `Scripts`, `DynamicProfiles` | `user-config` | ☐ false |
| `parsers`, `private`, sockets, locks, `version.txt` | `runtime` / `meta` | ☐ false |

**SUMMARY / SAFE TO RECLAIM** for **iterm2** must document that multiple `iterm2env*` trees may
share APFS hardlink / shared-inode blocks; the logical TOTAL can **overstate** freeable space;
confirm with `du -sh` on the parent. Expect roughly one env of reclaim, not the sum of aliases.
HOW TO PURGE is CLI-first: `$ disk-usage-analyser scan` on iTerm2; `$ du -sh` parent / envs;
optional logs after quit; optional remove all `iterm2env*` together after quit (not a single
alias only); **never** `rm -rf`; never mark DynamicProfiles/Scripts/user-config as usually-safe
bulk purge. RAW may include scan, du, optional `ls -li` notes.

**Explanation** (in-memory / `--json`) carries: absolute **path** (scope for packs), **kind**,
**totalSize**, **confidence** (`high`|`medium`|`low`), **summary** lines, **breakdown** entries
(name, path, size, role, notes, **`reclaimable` bool**), **reclaim** advice items
(`safeToReclaim`, title, detail), **howToPurge** steps (`title`, `officialCommand`, `removes`,
optional `notes`), **rawCommands** (grouped `disk-usage-analyser` vs `system` / ecosystem),
and for **codex-home** with a readable logs db optional/required **`logsDb`**
(`path`, `size`, `rows`, `samples[]`).

**howToPurge.officialCommand** is **CLI-primary plain text** (may be multi-line; `#` comments
OK). UI navigation (e.g. Android Studio Device Manager, Xcode Settings → Devices) belongs only
in **Notes** as `UI (optional): …`, never as the main official command body. For **xcode**,
prefer **`xcrun simctl …`** and/or **`disk-usage-analyser scan …`** — **never** `rm -rf`
(even though the frontend cleanup popover historically listed shell nukes). For **grok-home**,
prefer **`disk-usage-analyser scan …`** inspect-first reclaim steps — **never** `rm -rf`; never
treat auth/config as usually-safe. For **android-sdk**, prefer **`sdkmanager …`** and/or
**`disk-usage-analyser scan …`** — **never** `rm -rf`. For **iterm2**, prefer
**`disk-usage-analyser scan …`** and/or **`du -sh …`** inspect-first reclaim steps — **never**
`rm -rf`; document hardlink overcount in summary/reclaim. For **codex-home**, prefer
**`disk-usage-analyser scan …`** inspect-first plus safe **`logs_2.sqlite`** reclaim
(quit Codex; `mv` backup + wal/shm and/or `sqlite3 DELETE FROM logs; VACUUM;`) for
logs/sessions/cache/tmp — **never** `rm -rf`; never treat auth/config/`state_5` as usually-safe
logs reclaim. JSON stores plain strings: **no ANSI**, and **no leading `$ `** (the human
formatter adds `$` and color).

**Human output** (fixed section order, exact headers):

1. `PATH:`, `KIND:`, `TOTAL:`, `CONFIDENCE:` summary lines
2. `SUMMARY` then indented lines
3. `BREAKDOWN` then an **aligned table** (sorted **size DESC**, ties **name ASC**):
   columns **`SIZE`** (right, `FormatCompactHumanSize`), **`NAME`**, **`ROLE`** (bare role,
   no `role=` prefix), **`RECLAIMABLE`** (Unicode **`☑`** / **`☐` only** — never ASCII
   `[x]`/`[ ]`, never the words `true`/`false`), optional **`NOTES`** (omit column when all
   notes empty). Empty breakdown still prints `  (empty)`. Indent body lines under
   `BREAKDOWN` (e.g. two spaces).
4. `LOGS DB` (codex-home only, when `logs_*.sqlite` present) — PATH, SIZE, ROWS, SAMPLE
   (last 3, newest first). Omitted for other kinds / missing db.
5. `SAFE TO RECLAIM` bullet advice — must **never** contain `rm -rf` (or equivalent
   destructive one-liners)
6. `HOW TO PURGE` — numbered options with **Official command** and **Removes** (what files/data
   each purge deletes). Runnable official lines are shown with a leading **`$ `**; comment
   lines (`# …`) do not get `$`.
7. `RAW COMMANDS` with at least one `disk-usage-analyser scan` line for the path/scope, plus
   optional system/ecosystem commands (informational; tools need not be installed). Group
   headers remain `# group`; runnable lines use **`$ `**.

**BREAKDOWN reclaim tiers** (role → checkbox + ROLE/checkbox color when color is on):

| Tier | Roles | RECLAIMABLE | ROLE color | Checkbox color |
|------|-------|-------------|------------|----------------|
| **Reclaimable** | `cache`, `web-cache`, `tmp`, `temp`, `snapshot`, `backup`, `derived-data`, `docs-cache`, `device-support`, `simulator`, `installer-cache`, `session-logs`, `logs`, `app-logs-db`, `history`, `snapshots`, `system-images`, `skins`, `sources`, `python-env`, `python-env-alias` | `☑` | green | **green** |
| **Caution** | `chat-db`, `search-index`, `user-data`, `sdcard`, `config`, `session`, `web-storage`, `archives`, `project-meta`, `app-data`, `runtime-db`, `app-state-db`, `build-tools`, `platforms`, `emulator`, `platform-tools`, `cmdline-tools`, `licenses`, `meta`, `user-config`, `helper-binary`, `app-db`, `state` | `☐` | yellow | plain |
| **Neutral** | `directory`, `file`, `other`, `runtime`, `diagnostics`, `child`, … | `☐` | plain | plain |

JSON `breakdown[].reclaimable` is the same tier signal as a **boolean** (no checkbox glyphs,
no ANSI in any JSON field). Sort is shared for human + JSON after measure.

**Generic-dir basename → role**: when measuring top-level children that would otherwise be
only `directory`/`file`, remap basenames (case-insensitive) so reclaim signals work without
a specialized kind — e.g. `tmp`/`temp`/`.tmp` → `tmp`; `cache`/`Caches`/names ending in
`Cache`/`_cacache` → `cache` (keep `web-cache` if already specialized); `snapshots` →
`snapshot`; `sqlite-backup`/`idb-backup` → `backup`. Specialized kinds keep their own role
maps; tiers apply the same way.

**Shell prompt + color (human only)**:

- Runnable command lines (HOW TO PURGE official lines and RAW COMMANDS) use **`$ `** prefix.
- When color is **on**, only the **base command** token (argv0 / tool name, e.g. `emulator`,
  `sdkmanager`, `go`, `npm`, `brew`, `avdmanager`, `xcrun`, `disk-usage-analyser`) is **green**
  (`\x1b[32m` or bold green); the rest of the line is default. **`$` is not green**.
- When color is **on**, BREAKDOWN **ROLE cells** may be green (reclaimable tier) or yellow
  (caution tier). Reclaimable **`☑` is also green-wrapped** (same SGR as reclaimable ROLE);
  non-reclaimable **`☐` stays plain** (never green/yellow). Column padding uses **visible**
  width (ANSI / multi-byte glyphs must not break alignment).
- Color is **on** when `--color=always` (or bare `--color` if accepted as always), or when
  mode is **auto** (default) and stdout is a TTY.
- Color is **off** when `--color=never`, when auto and stdout is not a TTY, or when
  **`NO_COLOR`** is set (auto disabled). **`--color=always` forces color on** even if
  `NO_COLOR` is set (explicit user override). When color is off, glyphs are still `☑`/`☐`
  with zero ANSI in BREAKDOWN.
- Test harness writes to **non-TTY buffers** → default **auto** yields **no ANSI** unless
  `--color=always` is passed.

Stdout ends with a **trailing blank line** after the last content line. Sizes use the same
human formatting style as scan (`FormatCompactHumanSize` / similar).

**`--all-kinds` multi-pack mode**: analyses **all registered v1 packs** (`xcode`, `grok`,
`android-sdk`, `iterm2`, `codex`) under optional scope (default home). Missing pack roots get
status `missing` and size 0; overall **exit 0**. Human layout (mode A):

1. Header: `SCOPE:`, `MODE: all-kinds`, `TOTAL (present):`, counts
2. `INDEX` table: SIZE, KIND, STATUS, PATH (optional short NOTE) — all 5 cli kinds;
   **present** rows sorted **size DESC**
3. Present kinds each with separator + mini-explain (`PATH`/`KIND`/`TOTAL`/`SUMMARY`/
   `BREAKDOWN` + short SAFE/HOW TO PURGE OK)
4. Optional `MISSING` list
5. Trailing blank line
6. Never `rm -rf`

JSON (`--json --all-kinds`) is an **AllKindsResult** envelope (not a single Explanation):

```
{scope, totalSize, kinds:[{kind, cliKind, path, status, totalSize, ...}]}
```

`totalSize` sums **present** packs only. `cliKind` is the registry id (`grok` / `codex`);
`kind` is the output kind id (`grok-home` / `codex-home`).

**CLI / dispatch**: PATH is required unless `--kind` or `--all-kinds` is set (missing PATH
without either, or non-existent PATH → non-zero exit). `-h/--help` documents PATH (required
unless `--kind` or `--all-kinds`), **`--kind`**, **`--all-kinds`**, `--json`, and **`--color`**.
**`run.RunWithOptions`** routes `args[0] == "explain"` to the explain package before the
web-server branch. Root help lists `explain [PATH]`. No auto-delete. Test harness may inject
**`CLIOptions.HomeDir`** via `req.HomeDir` so `--kind` / `--all-kinds` without PATH never
touches the real user home.

## Decision Tree

```
explain-cli/
├── cli/                              # Meta CLI surface (help + routing)
│   ├── help/                         # explain -h: PATH, --kind, --all-kinds, --json, --color
│   ├── root-help/                    # root -h lists explain [PATH]
│   └── dispatch/                     # run.Run routes explain; no web server
├── errors/                           # Invalid inputs
│   ├── missing-path/                 # no PATH, no --kind, no --all-kinds → non-zero
│   ├── missing-file/                 # path does not exist → non-zero
│   ├── unknown-kind/                 # --kind not-a-kind → non-zero; lists xcode, grok, android-sdk, iterm2, codex
│   └── all-kinds-with-kind/          # --all-kinds --kind xcode → mutual exclusion non-zero
├── kinds/                            # Successful human explain by reclaim kind
│   ├── android-avd/
│   │   ├── dir/                      # *.avd directory fixture
│   │   └── file/                     # file inside AVD → android-avd context
│   ├── android-sdk/                  # Android SDK ~/Library/Android/sdk reclaim kind
│   │   ├── dir/                      # …/Android/sdk dir → android-sdk; roles; ☑/☐; CLI purge
│   │   ├── file/                     # file under SDK → same kind (ContextRoot=SDK root)
│   │   ├── scope/                    # --kind android-sdk <fakeHome> → {home}/Library/Android/sdk
│   │   └── no-path/                  # --kind android-sdk; HomeDir inject; PATH optional
│   ├── generic/
│   │   ├── dir/                      # random dir → generic-dir (not seatalk)
│   │   └── file/                     # random file → generic-file
│   ├── go-build-cache/               # .../go-build layout
│   ├── npm-cache/                    # .npm/_cacache layout
│   ├── homebrew-cache/               # Caches/Homebrew layout
│   ├── generic-qcow2/                # lone *.qcow2 file (not under AVD)
│   ├── seatalk-app-support/
│   │   ├── dir/                      # Application Support/SeaTalk → seatalk-app-support
│   │   └── file/                     # main_*.sqlite under SeaTalk → same kind
│   ├── grok-home/                    # Grok CLI ~/.grok reclaim kind
│   │   ├── dir/                      # …/.grok dir → grok-home; roles; ☑/☐; CLI purge
│   │   ├── file/                     # file under .grok → same kind (ContextRoot=.grok)
│   │   ├── scope/                    # --kind grok <fakeHome> → {home}/.grok
│   │   └── no-path/                  # --kind grok; HomeDir inject; PATH optional
│   ├── codex-home/                   # Codex CLI ~/.codex reclaim kind + LOGS DB (shape A)
│   │   ├── dir/                      # …/.codex dir → codex-home; LOGS DB; safe logs_2 purge; ☑/☐
│   │   ├── file/                     # file under .codex → same kind (ContextRoot=.codex)
│   │   ├── scope/                    # --kind codex <fakeHome> → {home}/.codex + logs_2 purge
│   │   └── no-path/                  # --kind codex; HomeDir inject; PATH optional
│   ├── iterm2/                       # iTerm2 ~/Library/Application Support/iTerm2 reclaim kind
│   │   ├── dir/                      # …/iTerm2 dir → iterm2; roles; ☑/☐; hardlink; scan/du purge
│   │   ├── file/                     # file under iTerm2 → same kind (ContextRoot=iTerm2)
│   │   ├── scope/                    # --kind iterm2 <fakeHome> → {home}/…/iTerm2
│   │   └── no-path/                  # --kind iterm2; HomeDir inject; PATH optional
│   └── xcode/                        # --kind xcode multi-root pack
│       ├── scope/                    # --kind xcode <fixtureHome>: 5 roots, ☑/☐, CLI purge
│       └── no-path/                  # --kind xcode; HomeDir inject; PATH optional
└── output/                           # Cross-cutting output contracts
    ├── all-kinds-index/              # --all-kinds HomeDir: INDEX all 5; present size DESC; statuses
    ├── all-kinds-detail/             # --all-kinds SCOPE PATH: present KIND: + BREAKDOWN mini-explain
    ├── all-kinds-all-missing/        # empty home: exit 0; all missing
    ├── json-all-kinds/               # --json --all-kinds: scope/totalSize/kinds[]; no ANSI
    ├── json-android-avd/             # --json: fields + no ANSI + no $ in officialCommand
    ├── json-android-sdk/             # --json android-sdk: roles + reclaimable; sdkmanager/scan; no rm -rf
    ├── json-seatalk/                 # --json seatalk-app-support: kind, roles, howToPurge
    ├── json-xcode/                   # --json --kind xcode: roles + reclaimable; no rm -rf
    ├── json-grok/                    # --json grok-home: roles + reclaimable; scan; no rm -rf
    ├── json-codex/                   # --json codex-home: logsDb; safe logs_2 howToPurge; no rm -rf
    ├── json-iterm2/                  # --json iterm2: roles + reclaimable; hardlink; scan/du; no rm -rf
    ├── json-breakdown-sorted/        # --json: breakdown size DESC + reclaimable bool
    ├── no-rm-rf/                     # stdout must not contain rm -rf
    ├── raw-commands-scan/            # RAW COMMANDS: $ disk-usage-analyser scan
    ├── purge-cli-first/              # android-avd HOW TO PURGE is CLI-first
    ├── dollar-prefix/                # human $ prefix on runnable command lines
    ├── color-force/                  # --color=always → green base command ANSI
    ├── color-never/                  # --color=never → no ANSI
    ├── breakdown-table-align/        # BREAKDOWN aligned SIZE/NAME/ROLE/RECLAIMABLE/NOTES
    ├── breakdown-sorted-desc/        # human BREAKDOWN rows size DESC
    ├── breakdown-reclaimable-checkbox/ # ☑/☐ by reclaim tier (never [x]/[ ]/true/false)
    ├── breakdown-role-color/         # --color=always: ROLE green/yellow; ☑ green; ☐ plain
    ├── breakdown-color-never/        # --color=never: no ANSI; ☑/☐ still present
    └── generic-dir-cache-tmp/        # generic-dir Cache/tmp → cache/tmp roles + ☑
```



## Test Index

| Leaf | Mode | Description |
|------|------|-------------|
| cli/help | cli | `-h` documents explain usage, PATH, `--kind` (xcode, grok, android-sdk, iterm2, codex), `--all-kinds`, `--json`, `--color`. PATH optional when `--kind` or `--all-kinds`. |
| cli/root-help | dispatch | Root `-h` lists `explain [PATH]`; no web server. |
| cli/dispatch | dispatch | `run.RunWithOptions(["explain", path])` explains without starting web server. |
| errors/missing-path | cli | No PATH, no `--kind`, no `--all-kinds` → non-zero exit + error. |
| errors/missing-file | cli | Non-existent PATH → non-zero exit + error. |
| errors/unknown-kind | cli | `--kind not-a-kind` → non-zero; lists supported kinds (`xcode`, `grok`, `android-sdk`, `iterm2`, `codex`). |
| errors/all-kinds-with-kind | cli | `--all-kinds --kind xcode` → non-zero; mutual exclusion message. |
| kinds/android-avd/dir | cli | AVD dir → `android-avd`; CLI-first purge; `$` on cmds; trailing blank. |
| kinds/android-avd/file | cli | Path to `userdata-qemu.img.qcow2` inside AVD → kind `android-avd`. |
| kinds/android-sdk/dir | cli | SDK dir → `android-sdk`; roles system-images/emulator/sources/…; ☑/☐; sdkmanager/scan purge; `$`; no `rm -rf`. |
| kinds/android-sdk/file | cli | Path under SDK → kind `android-sdk` (ContextRoot = SDK root). |
| kinds/android-sdk/scope | cli | `--kind android-sdk <fakeHome>`: measures `{home}/Library/Android/sdk`; roles + reclaim; CLI purge; no `rm -rf`. |
| kinds/android-sdk/no-path | cli | `--kind android-sdk` with `req.HomeDir` inject; no PATH; measures `{HomeDir}/Library/Android/sdk`. |
| kinds/generic/dir | cli | Random directory → `generic-dir`; breakdown + `$` scan still present. |
| kinds/generic/file | cli | Random file → `generic-file`. |
| kinds/go-build-cache | cli | Path ending in `go-build` → `go-build-cache`; `$ go clean` in HOW TO PURGE. |
| kinds/npm-cache | cli | `.npm` with `_cacache` → `npm-cache`; `$ npm` purge line. |
| kinds/homebrew-cache | cli | `Caches/Homebrew` → `homebrew-cache`; `$ brew` purge line. |
| kinds/generic-qcow2 | cli | Lone `disk.qcow2` file → `generic-qcow2`. |
| kinds/seatalk-app-support/dir | cli | SeaTalk Application Support dir → `seatalk-app-support`; roles; reclaim tiers; osascript + cache/backup HOW TO PURGE; `$`; no `rm -rf`. |
| kinds/seatalk-app-support/file | cli | Path to `main_*.sqlite` under SeaTalk → kind `seatalk-app-support` (parent context). |
| kinds/grok-home/dir | cli | `.grok` dir → `grok-home`; roles installer-cache/session-logs/cache/logs/config; ☑/☐; scan purge; `$`; no `rm -rf`. |
| kinds/grok-home/file | cli | Path to `config.toml` under `.grok` → kind `grok-home` (parent context). |
| kinds/grok-home/scope | cli | `--kind grok <fakeHome>`: measures `{home}/.grok`; roles + reclaim; CLI purge; no `rm -rf`. |
| kinds/grok-home/no-path | cli | `--kind grok` with `req.HomeDir` inject; no PATH; measures `{HomeDir}/.grok`. |
| kinds/codex-home/dir | cli | `.codex` dir → `codex-home`; roles app-logs-db/session-logs/cache/tmp/config; LOGS DB ROWS:5 + SAMPLE last 3; ☑/☐; scan + safe `logs_2.sqlite` purge (quit/mv/sqlite3); `$`; no `rm -rf`. |
| kinds/codex-home/file | cli | Path to `config.toml` under `.codex` → kind `codex-home` (parent context); LOGS DB still present. |
| kinds/codex-home/scope | cli | `--kind codex <fakeHome>`: measures `{home}/.codex`; roles + LOGS DB + reclaim; safe logs_2 purge; no `rm -rf`. |
| kinds/codex-home/no-path | cli | `--kind codex` with `req.HomeDir` inject; no PATH; measures `{HomeDir}/.codex`; safe logs_2 purge. |
| kinds/iterm2/dir | cli | iTerm2 App Support dir → `iterm2`; roles python-env/python-env-alias/logs/meta/user-config; ☑/☐; hardlink wording; scan/du purge; `$`; no `rm -rf`. |
| kinds/iterm2/file | cli | Path under iTerm2 → kind `iterm2` (ContextRoot = iTerm2 dir). |
| kinds/iterm2/scope | cli | `--kind iterm2 <fakeHome>`: measures `{home}/Library/Application Support/iTerm2`; roles + reclaim; CLI purge; no `rm -rf`. |
| kinds/iterm2/no-path | cli | `--kind iterm2` with `req.HomeDir` inject; no PATH; measures `{HomeDir}/Library/Application Support/iTerm2`. |
| kinds/xcode/scope | cli | `--kind xcode <scope>`: multi-root pack; roles; archives `☐`; others `☑`; size DESC; CLI purge (`xcrun`/`simctl`/scan); no `rm -rf`. |
| kinds/xcode/no-path | cli | `--kind xcode` with `req.HomeDir` inject; no PATH; same pack under injected home. |
| output/all-kinds-index | cli | `--all-kinds` + `HomeDir` multi-pack fixture: header SCOPE/MODE/TOTAL; INDEX lists all 5 cli kinds; present size DESC; present/missing statuses; no `rm -rf`; trailing blank. |
| output/all-kinds-detail | cli | `--all-kinds <SCOPE>`: present kinds have detail `KIND:` + `BREAKDOWN` (incl. codex-home); no detail success for missing xcode; no `rm -rf`; trailing blank. |
| output/all-kinds-all-missing | cli | `--all-kinds` empty home: exit 0; INDEX/JSON all missing (or all status missing); total 0. |
| output/json-all-kinds | cli | `--json --all-kinds`: envelope `scope`/`totalSize`/`kinds[]` (5 packs); each entry `kind`/`cliKind`/`path`/`status`/`totalSize`; present sizes; no ANSI; no `rm -rf`. |
| output/json-android-avd | cli | `--json` fields; no ANSI; officialCommand CLI-first without `$`. |
| output/json-android-sdk | cli | `--json` android-sdk: kind, path, roles + reclaimable bools; howToPurge sdkmanager/scan CLI-first without `$`/ANSI/`rm -rf`. |
| output/json-seatalk | cli | `--json` seatalk-app-support: kind, breakdown roles, howToPurge (osascript/cache/backup), no ANSI/`$`. |
| output/json-xcode | cli | `--json --kind xcode <scope>`: kind, path, roles + reclaimable bools; howToPurge CLI-first without `$`/ANSI/`rm -rf`. |
| output/json-grok | cli | `--json` grok-home: kind, path, roles + reclaimable bools; howToPurge scan CLI-first without `$`/ANSI/`rm -rf`. |
| output/json-codex | cli | `--json` codex-home: kind, path, roles + reclaimable; **logsDb** (rows=5, samples≤3 newest first); howToPurge scan + safe logs_2 reclaim (quit/mv/sqlite3; not state_5/auth) without `$`/ANSI/`rm -rf`. |
| output/json-iterm2 | cli | `--json` iterm2: kind, path, roles + reclaimable bools; hardlink wording; howToPurge scan/du CLI-first without `$`/ANSI/`rm -rf`. |
| output/json-breakdown-sorted | cli | `--json` breakdown sorted size DESC; each entry has `reclaimable` bool; AVD tier map. |
| output/no-rm-rf | cli | Full stdout must not contain `rm -rf` (case-insensitive). |
| output/raw-commands-scan | cli | RAW COMMANDS includes `$ disk-usage-analyser scan` and path. |
| output/purge-cli-first | cli | android-avd human HOW TO PURGE prefers emulator/avdmanager; UI only in Notes. |
| output/dollar-prefix | cli | Human RAW + HOW TO PURGE runnable lines use `$ `; comments do not. |
| output/color-force | cli | `--color=always` greens base command (e.g. `go`); `$` not green. |
| output/color-never | cli | `--color=never` → no ANSI escapes in stdout. |
| output/breakdown-table-align | cli | Human BREAKDOWN is aligned multi-column table; SIZE right-aligned. |
| output/breakdown-sorted-desc | cli | Human BREAKDOWN rows largest→smallest (AVD fixture name order). |
| output/breakdown-reclaimable-checkbox | cli | RECLAIMABLE `☑` for snapshot; `☐` for caution roles; never `[x]`/`[ ]`/true/false. |
| output/breakdown-role-color | cli | `--color=always`: ROLE green/yellow by tier; reclaimable `☑` green-wrapped; `☐` plain. |
| output/breakdown-color-never | cli | `--color=never`: no ANSI in BREAKDOWN; `☑`/`☐` still present. |
| output/generic-dir-cache-tmp | cli | generic-dir `Cache`/`tmp` → cache/tmp roles with `☑`; size DESC. |

## How to Run

```sh
doctest vet ./tests/explain-cli
doctest test ./tests/explain-cli
```

```go
import (
	"bytes"
	"context"
	"strings"
	"testing"

	"disk-usage-analyser/explain"
	"disk-usage-analyser/run"
)

type Request struct {
	Mode       string // "cli" (default) or "dispatch"
	FixtureDir string
	TargetPath string // absolute path passed as PATH / scope (dir or file)
	HomeDir    string // optional; inject CLIOptions.HomeDir for --kind/--all-kinds without PATH
	Args       []string
	Stdout     *bytes.Buffer
	Stderr     *bytes.Buffer
}

type Response struct {
	Stdout           string
	Stderr           string
	ExitCode         int
	Err              error
	ServerWasStarted bool
}

func Run(t *testing.T, req *Request) (*Response, error) {
	if req.Stdout == nil {
		req.Stdout = &bytes.Buffer{}
	}
	if req.Stderr == nil {
		req.Stderr = &bytes.Buffer{}
	}

	switch req.Mode {
	case "dispatch":
		serverStarted := false
		err := run.RunWithOptions(context.Background(), req.Args, run.Options{
			Stdout: req.Stdout,
			Stderr: req.Stderr,
			StartServer: func(context.Context, run.ServerOptions) error {
				serverStarted = true
				t.Fatal("explain dispatch must not start web server")
				return nil
			},
		})
		return &Response{
			Stdout:           strings.ReplaceAll(req.Stdout.String(), "\r\n", "\n"),
			Stderr:           strings.ReplaceAll(req.Stderr.String(), "\r\n", "\n"),
			Err:              err,
			ServerWasStarted: serverStarted,
		}, nil

	default: // "cli"
		exitCode, err := explain.RunCLI(req.Args, explain.CLIOptions{
			Stdout:  req.Stdout,
			Stderr:  req.Stderr,
			HomeDir: req.HomeDir,
		})
		return &Response{
			Stdout:   strings.ReplaceAll(req.Stdout.String(), "\r\n", "\n"),
			Stderr:   strings.ReplaceAll(req.Stderr.String(), "\r\n", "\n"),
			ExitCode: exitCode,
			Err:      err,
		}, nil
	}
}
```
