# SKILL: Analyse disk usage with scan

Agent playbook for investigating disk usage with `disk-usage-analyser scan`.

## Quick start

1. Run a JSON tree scan from the user's home directory (or the path they care about):

   ```bash
   disk-usage-analyser scan ~ --json
   ```

2. Parse `tree.children` and rank entries by `size` (descending).

3. Drill into large directories:

   ```bash
   disk-usage-analyser scan <path> --json
   ```

4. Optional cross-check with the web UI `/tmp-analyse` flow when the user wants a second view.

5. Report ranked findings and practical cleanup suggestions (old caches, build artifacts, duplicate downloads, etc.).

## CLI reference

```bash
disk-usage-analyser scan [PATH] [--json] [--threshold SIZE] [--max-depth N]
```

| Flag | Default (text) | Default (`--json`) | Behavior |
|------|----------------|----------------------|----------|
| `--threshold` | `1M` | `1M` | Hide nodes with recursive size below threshold (`0` = show all) |
| `--max-depth` | `3` | `24` | Max branch expansion (`0` = unlimited) |
| `--json` | off | — | Emit nested `TreeResult` JSON instead of text tree |

### Examples

```bash
# Text tree with defaults (threshold 1M, depth 3)
disk-usage-analyser scan ~/Library

# Show everything, unlimited depth
disk-usage-analyser scan ~/Downloads --threshold 0 --max-depth 0

# Machine-readable drill-down
disk-usage-analyser scan ~/Library/Developer --json --max-depth 24
```

## JSON shape

`--json` emits one `TreeResult` object:

- `path` — absolute scanned root
- `totalSize` — full recursive bytes at root (includes hidden/filtered nodes)
- `threshold` — threshold in bytes
- `maxDepth` — configured depth limit
- `tree` — nested root node (`name` is `"."`) with `children` sorted by size descending (dirs before files on tie)

There is no flat `items` array in JSON mode.

## Workflow tips

- Start broad (`~` or `/`), then recurse into the largest `isDir: true` children.
- Use `--threshold 0` when hunting small but numerous files.
- Use `--max-depth 1` for a quick top-level overview before deep scans.
- Compare CLI totals with the web app when the user wants visual confirmation.