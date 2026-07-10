---
name: analyse-my-disk-space
description: >
  Analyse macOS disk usage with one disk-usage-analyser scan --json --max-depth 6,
  then inspect/query that JSON offline via disk-usage-analyser scan --inspect (no
  multi-phase re-scans). Use when the user asks to analyse disk space, free space,
  what's using storage, or runs a scan workflow. For Android emulator/SDK/AVD
  paths or disk image suffixes, also load: disk-usage-analyser skill --show android-images.
---

# SKILL: Analyse disk usage with scan + --inspect

Agent playbook for investigating disk usage with `disk-usage-analyser scan` (live walk
or offline `--inspect`).

## Topics

Load nested topics with `skill --show` (both flag orders work):

```bash
disk-usage-analyser skill --show
disk-usage-analyser skill --show android-images
disk-usage-analyser skill android-images --show
disk-usage-analyser skill --list
```

| Path | When to load |
|------|----------------|
| `android-images` | Paths under `~/.android/`, `~/Library/Android/`, AVD names, emulator image suffixes (`.qcow2`, `ram.bin`, …) |

## Important: one scan, then offline analysis

**Do not multi-phase re-scan** (e.g. `scan ~ --max-depth 1`, then `scan ~/Library`, then `scan ~/.android`, …).

Each live `scan` invocation walks the filesystem **fully** to compute sizes. `--max-depth` only limits how deep the **emitted tree** is expanded in the JSON/text output — it does not make the walk shallower. Extra scans waste wall time (often 2–3×).

**Preferred workflow:** one deep-enough JSON capture, then query it with **`scan --inspect`** (no Python required). There is no standalone `inspect` subcommand.

## Quick start

1. Run a **single** JSON scan from the user's home directory (or the path they care about). Prefer depth **6** so typical reclaim targets (caches, App Support, AVD disks) appear as nodes without a second scan:

   ```bash
   disk-usage-analyser scan ~ --json --max-depth 6 > /tmp/disk-scan.json 2>/tmp/disk-scan.err
   ```

2. **Inspect offline** (do not re-scan):

   ```bash
   # Tree view (inspect defaults: max-depth 1, min 0) plus SOURCE line
   disk-usage-analyser scan --inspect /tmp/disk-scan.json

   # Top consumers (Option B: tree section + TOP N; default N=20, root excluded)
   disk-usage-analyser scan --inspect /tmp/disk-scan.json --top 40

   # Drill into a path already present in the JSON (tree only when --at alone)
   disk-usage-analyser scan --inspect /tmp/disk-scan.json --at "$HOME/Library"
   disk-usage-analyser scan --inspect /tmp/disk-scan.json --at "$HOME/Library/Caches"

   # Search by name/path substring or file suffix (tree + match section)
   disk-usage-analyser scan --inspect /tmp/disk-scan.json --find go-build
   disk-usage-analyser scan --inspect /tmp/disk-scan.json --suffix .qcow2 --min 100M

   # Machine-readable ViewResult for agents
   disk-usage-analyser scan --inspect /tmp/disk-scan.json --json --top 30
   disk-usage-analyser scan --inspect /tmp/disk-scan.json --json --at "$HOME/.android"
   ```

3. Report ranked findings and practical cleanup suggestions (old caches, build artifacts, duplicate downloads, AVDs, etc.).

4. **Re-scan only if needed:** if a large directory sits at `depth == maxDepth` with no `children` (branch truncated) and you still need its internals, run **one** targeted scan, then inspect that file:

   ```bash
   disk-usage-analyser scan <that-path> --json --max-depth 6 > /tmp/disk-scan-focus.json
   disk-usage-analyser scan --inspect /tmp/disk-scan-focus.json --top 30
   ```

5. Optional cross-check with the web UI `/tmp-analyse` flow when the user wants a second view.

## Topic: Android images and AVDs

When the scan or user focuses on **Android** storage, load and follow:

```bash
disk-usage-analyser skill --show android-images
```

Load that topic when any of the following apply:

- Paths under `~/.android/`, `~/Library/Android/`, or names like `*.avd`, AVD `*.ini`
- Mentions of AVD, emulator, Android SDK, system-images, `avdmanager`, `sdkmanager`
- Disk image / snapshot suffixes typical of emulator VMs: `.img`, `.qcow2`, `ram.bin`, `sdcard.img`, `userdata-qemu.img*`, `cache.img*`, `encryptionkey.img*`

Prefer querying the **existing** home JSON with `scan --inspect --find` / `--suffix` / `--at`. Only scan `~/.android` alone if the home tree did not cover it. That topic covers AVD wipe vs delete, recreate, and backup/restore.

## CLI reference

### `scan` — live walk or offline inspect

```bash
disk-usage-analyser scan [PATH] [flags]
disk-usage-analyser scan --inspect FILE [flags]
```

| Flag | Live default (text) | Live `--json` capture | `--inspect` default | Behavior |
|------|---------------------|------------------------|---------------------|----------|
| `--min SIZE` | `1M` | `1M` | `0` | Hide tree/match nodes below SIZE (`0` = show all). Replaces removed `--threshold`. |
| `--max-depth N` | `3` | `24` | `1` | Max **tree section** depth (`0` = unlimited). Live size walk is always full recursive. Match ranking uses the full loaded tree. |
| `--json` | off | — | — | Live pure capture → bare `TreeResult` with field **`min`**. Otherwise → `ViewResult`. |
| `--inspect FILE` | — | — | — | Phase 1 loads JSON (FILE `-` = stdin); no FS walk. |
| `--top N` | off | off | off | Option B match section; default N=`20` when match section is active via `--find`/`--suffix` without N. |
| `--at PATH` | off | off | off | Focus tree at PATH; alone = tree only (no TOP). With `--top`/`--find`/`--suffix` also emits matches. |
| `--find SUBSTR` | off | off | off | Case-insensitive path/name match → match section |
| `--suffix SUFFIX` | off | off | off | Name/path ends with suffix → match section |
| `--include-root` | off | off | off | Include scan root in global TOP/find rankings |

**Removed:** `--threshold`, standalone `inspect` subcommand.

Live query flags (`--top` / `--find` / `--suffix`) also work on a live scan (same Option B output).

### Examples

```bash
# Preferred agent capture
disk-usage-analyser scan ~ --json --max-depth 6 > /tmp/disk-scan.json

# Offline queries (no disk re-walk)
disk-usage-analyser scan --inspect /tmp/disk-scan.json --top 30
disk-usage-analyser scan --inspect /tmp/disk-scan.json --at ~/Library/Application\ Support
disk-usage-analyser scan --inspect /tmp/disk-scan.json --find SeaTalk --json
disk-usage-analyser scan --inspect /tmp/disk-scan.json --suffix .qcow2 --min 50M

# Live tree + top without writing JSON first
disk-usage-analyser scan ~ --min 1M --top 20

# Rare: one follow-up if a branch was truncated at max-depth
disk-usage-analyser scan ~/Projects/gopath/src/git.some.com --json --max-depth 6 > /tmp/focus.json
disk-usage-analyser scan --inspect /tmp/focus.json --top 20
```

## JSON shapes

### `scan --json` (pure capture, no query) → `TreeResult`

- `path` — absolute scanned root
- `totalSize` — full recursive bytes at root (includes hidden/filtered nodes; **independent of maxDepth**)
- `min` — min filter in bytes (JSON field is **`min`**, not `threshold`)
- `maxDepth` — configured tree emission depth limit
- `tree` — nested root node (`name` is `"."`) with `children` sorted by size descending (dirs before files on tie)

There is no flat `items` array. Use **`scan --inspect`** instead of hand-rolling a flatten loop.

### `scan --inspect --json` / live query `--json` → `ViewResult`

- `scanPath`, `totalSize`, `min`, `maxDepth`
- `sourceFile` — present for `--inspect` when FILE is a path
- `tree` — view-pruned tree (inspect default maxDepth 1)
- `matches[]` — optional `{name, path, size, isDir, depth}` when match section is active

## Human text (Option B)

Summary: `PATH:`, `TOTAL:`, `MIN:`, `MAX-DEPTH:`, and `SOURCE:` only when `--inspect`.

Then a box-drawing tree (name, then aligned size column). When `--top` / `--find` / `--suffix` are set, a `TOP N` section follows with lines `size  kind  d=N  path`. Stdout ends with a trailing blank line.

## Workflow tips

- **One live `scan --json`**, then only **`scan --inspect`** for ranking, drills, and suffix/name search.
- Depth **6** is usually enough for reclaim targets under `~`. JSON stays small (a few MB for a full home).
- `totalSize` at every included dir already reflects the full subtree — even when children are omitted past max-depth.
- Use `--min 0` on scan only when hunting small but numerous files (larger JSON).
- Avoid `--max-depth 1` as a first **live** scan if you plan to drill: you will just re-walk everything again. Prefer capture at 6 + inspect.
- Compare CLI totals with the web app when the user wants visual confirmation.
- If Android/AVD or image suffixes appear in top consumers, open **android-images** (`skill --show android-images`) before recommending delete.
