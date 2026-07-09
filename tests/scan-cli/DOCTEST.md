# Scan CLI

CLI tests for the two-phase `disk-usage-analyser scan` pipeline:

1. **Phase 1 — TreeSource**: live FS discovery, or load a saved JSON via `--inspect FILE`.
2. **Phase 2 — View**: shared tree formatting and optional match ranking (`--top` / `--find` / `--suffix` / `--at`).

Covers `--min` (replaces `--threshold`), text and JSON output, inspect defaults, Option B
(tree section + match section when query flags warrant it), dispatch before the web-server
branch, and removal of the standalone `inspect` subcommand.

## Version

0.0.2

# DSN (Domain Specific Notion)

The **scan command** is a **two-phase** pipeline.

**Phase 1 (TreeSource)** produces a **TreeResult**: absolute **Path**, **TotalSize** (full
recursive root bytes), **Min** (display filter used when the tree was produced or re-viewed),
**MaxDepth**, and nested **Tree** (root name `"."`). Sources:

- **LiveTreeSource**: walk the filesystem from `PATH` (default: cwd) with concurrent
  **GlobalCache** deduplication. Counting always includes every file and directory. Live
  **emission** omits nodes with `size < min` and stops expanding past **maxDepth** (`0` =
  unlimited). Live defaults: **min** `1M`, text **maxDepth** `3`, JSON capture **maxDepth**
  `24`.
- **JSONTreeSource** (`--inspect FILE`, FILE `-` = stdin): decode a previously captured
  TreeResult JSON (field **`min`**, not `threshold`). No FS walk.

**Phase 2 (View)** renders a **tree section** and, when warranted, a **match section**
(Option B):

- Tree-only when no `--top` / `--find` / `--suffix` (including **`--at` alone**, which focuses
  the tree on a subtree with **no** match section).
- Tree **and** matches when `--top` and/or `--find` and/or `--suffix` (also when combined with
  `--at`). Match ranking uses the **full in-memory tree**, not the view max-depth prune.
  Default **`--top`** cap is **20** when the match section is active. Root is skipped in
  global rankings unless **`--include-root`**.
- Inspect defaults for the view: **maxDepth 1**, **min 0**. Live query views keep live min /
  depth defaults unless overridden.
- Human summary: `PATH:`, `TOTAL:`, `MIN:`, `MAX-DEPTH:`, and **`SOURCE:`** only when
  `--inspect`. Then blank line, box-drawing tree (name then aligned size column). Match
  block header `TOP N` then lines `size  kind  d=N  path`. Stdout ends with a trailing blank
  line.
- **`--json`**: pure capture (live, no query/view extras) → bare **TreeResult** with JSON
  field **`min`**. Otherwise → **ViewResult** (`scanPath`, `totalSize`, `min`, `maxDepth`,
  optional `sourceFile`, `tree`, `matches`).

**CLI**: only subcommand `scan` (no standalone `inspect`). Flags: `--inspect`, `--json`,
`--min`, `--max-depth`, `--top`, `--at`, `--find`, `--suffix`, `--include-root`, `-h/--help`.
**Removed**: `--threshold`, subcommand `inspect`. Invalid `--min` values and inspect load
errors exit non-zero. **`run.Run`** routes `scan` (including `--inspect`) before the web
server branch.

## Decision Tree

```
scan-cli/
├── tree/                                   # Live Scan API: TreeResult + ScanOptions
│   ├── basic-nested/                       # nested children with correct recursive sizes
│   ├── min-default/                        # default min 1M hides sub-1M nodes
│   ├── min-override/                       # Min=10M filters display
│   ├── max-depth-2/                        # branches stop at depth 2; sizes include deeper
│   ├── max-depth-zero/                     # maxDepth 0 = unlimited expansion
│   ├── max-depth-one/                      # only immediate children shown
│   └── root-total-includes-filtered/       # totalSize counts sub-min bytes
├── basics/                                 # fundamental tree shapes (min 0)
│   ├── empty-dir/
│   ├── files-only/
│   ├── mixed/
│   └── nested-dirs/
├── sorting/
│   └── by-size/
├── errors/
│   ├── missing-path/
│   ├── not-a-directory/
│   ├── invalid-min/                        # --min foo exits non-zero
│   ├── unknown-threshold-flag/             # --threshold rejected (no alias)
│   ├── inspect-missing-file/
│   ├── inspect-invalid-json/
│   └── inspect-with-path/                  # --inspect + positional PATH → error
├── cli/                                    # live CLI surface (no query section)
│   ├── default-cwd/
│   ├── json-flag/                          # capture TreeResult; field min
│   ├── json-tree-shape/                    # path/totalSize/min/maxDepth/tree; no items
│   ├── text-tree-format/                   # PATH/TOTAL/MIN/MAX-DEPTH + aligned tree
│   ├── text-tree-alignment/
│   ├── help/                               # scan, PATH, --json
│   ├── help-flags/                         # --min, --max-depth (not --threshold)
│   ├── help-min-and-inspect/               # --inspect, --top, --find, --suffix, --at
│   ├── explicit-path/
│   ├── dispatch/                           # run.Run scan without web server
│   ├── dispatch-inspect/                   # run.Run scan --inspect without web server
│   └── root-help-no-inspect-subcommand/    # root -h has no inspect subcommand
├── inspect/                                # Phase 1 = JSONTreeSource
│   ├── default-depth-1/                    # tree max-depth 1; SOURCE; min 0
│   ├── max-depth-2/                        # deeper nodes appear
│   ├── top-option-b/                       # tree section + TOP 2 matches
│   ├── at-alone/                           # focused tree; no TOP
│   ├── at-with-top/                        # focused tree + match section
│   ├── find/                               # tree + find matches
│   ├── suffix/                             # tree + suffix matches
│   ├── json-tree-only/                     # ViewResult tree; inspect defaults
│   └── json-view-top/                      # ViewResult tree + matches
└── live-query/                             # live FS + match section
    ├── top/                                # scan --top N: tree + TOP
    └── parity-inspect/                     # live ranking matches inspect of same tree
```

## Test Index

| Leaf | Mode | Description |
|------|------|-------------|
| tree/basic-nested | scan | Nested `tree.children` with correct recursive sizes. |
| tree/min-default | scan | Default min 1M hides sub-1M nodes from display. |
| tree/min-override | scan | `Min=10M` hides nodes below 10M. |
| tree/max-depth-2 | scan | `MaxDepth=2` truncates branches; ancestor sizes include deeper bytes. |
| tree/max-depth-zero | scan | `MaxDepth=0` shows all levels (subject to min). |
| tree/max-depth-one | scan | `MaxDepth=1` shows only root immediate children. |
| tree/root-total-includes-filtered | scan | `totalSize` includes sub-min bytes not shown in tree. |
| basics/empty-dir | scan | Empty root: `tree.children` empty, `totalSize` 0. |
| basics/files-only | scan | Two root files (100 B, 200 B): `totalSize` 300. |
| basics/mixed | scan | Root file 50 B + subdir with 150 B nested file: `totalSize` 200. |
| basics/nested-dirs | scan | Nested dirs roll up to one dir child 300 B. |
| sorting/by-size | scan | Root children: size desc; dirs before files on tie. |
| errors/missing-path | cli | Non-existent PATH → non-zero exit. |
| errors/not-a-directory | cli | File as PATH → not a directory error. |
| errors/invalid-min | cli | `--min foo` → non-zero exit and clear error. |
| errors/unknown-threshold-flag | cli | `--threshold` is not accepted (breaking rename). |
| errors/inspect-missing-file | cli | Missing inspect FILE → non-zero + error. |
| errors/inspect-invalid-json | cli | Corrupt inspect JSON → non-zero + error. |
| errors/inspect-with-path | cli | `--inspect FILE` with positional PATH → non-zero + clear error. |
| cli/default-cwd | cli | No args: cwd scan; summary uses `MIN:`. |
| cli/json-flag | cli | `--json` capture TreeResult; sizes nested under `tree`. |
| cli/json-tree-shape | cli | Capture JSON has `min` (not `threshold`), `maxDepth` 24. |
| cli/text-tree-format | cli | Text summary PATH/TOTAL/MIN/MAX-DEPTH + aligned tree. |
| cli/text-tree-alignment | cli | Nested tree sizes share one column. |
| cli/help | cli | `-h` documents scan, `[PATH]`, `--json`. |
| cli/help-flags | cli | `-h` documents `--min` and `--max-depth`; not `--threshold`. |
| cli/help-min-and-inspect | cli | `-h` documents `--inspect`, `--top`, `--at`, `--find`, `--suffix`. |
| cli/explicit-path | cli | Positional PATH; `MIN:` in summary. |
| cli/dispatch | dispatch | `run.RunWithOptions(["scan", ...])` does not start web server. |
| cli/dispatch-inspect | dispatch | `run.RunWithOptions(["scan", "--inspect", file])` works offline. |
| cli/root-help-no-inspect-subcommand | dispatch | Root `-h` does not list `inspect` as a subcommand. |
| inspect/default-depth-1 | cli | Inspect text: depth-1 children only; SOURCE; MIN 0; trailing blank. |
| inspect/max-depth-2 | cli | `--max-depth 2` shows depth-2 nodes. |
| inspect/top-option-b | cli | Tree section present and TOP 2 (root skipped). |
| inspect/at-alone | cli | Focused subtree tree; no TOP section. |
| inspect/at-with-top | cli | Focused tree plus match section. |
| inspect/find | cli | Tree + find match lines. |
| inspect/suffix | cli | Tree + suffix match lines. |
| inspect/json-tree-only | cli | Inspect `--json` ViewResult with tree; defaults. |
| inspect/json-view-top | cli | Inspect `--json --top` ViewResult with matches. |
| live-query/top | cli | Live `--top` emits tree + TOP section. |
| live-query/parity-inspect | cli | Live top ranking agrees with inspect of equivalent capture. |

## How to Run

```sh
doctest vet ./tests/scan-cli
doctest test ./tests/scan-cli
```

```go
import (
	"bytes"
	"context"
	"strings"
	"testing"

	"disk-usage-analyser/run"
	"disk-usage-analyser/usagescan"
)

type Request struct {
	Mode       string
	FixtureDir string
	Args       []string
	ScanOpts   *usagescan.ScanOptions
	Stdout     *bytes.Buffer
	Stderr     *bytes.Buffer
}

type Response struct {
	TreeResult       *usagescan.TreeResult
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
				t.Fatal("scan dispatch must not start web server")
				return nil
			},
		})
		return &Response{
			Stdout:           strings.ReplaceAll(req.Stdout.String(), "\r\n", "\n"),
			Stderr:           strings.ReplaceAll(req.Stderr.String(), "\r\n", "\n"),
			Err:              err,
			ServerWasStarted: serverStarted,
		}, nil

	case "cli":
		exitCode, err := usagescan.RunCLI(req.Args, usagescan.CLIOptions{
			Stdout: req.Stdout,
			Stderr: req.Stderr,
		})
		return &Response{
			Stdout:   strings.ReplaceAll(req.Stdout.String(), "\r\n", "\n"),
			Stderr:   strings.ReplaceAll(req.Stderr.String(), "\r\n", "\n"),
			ExitCode: exitCode,
			Err:      err,
		}, nil

	default:
		opts := defaultScanOpts()
		if req.ScanOpts != nil {
			opts = *req.ScanOpts
		}
		result, err := usagescan.Scan(req.FixtureDir, opts)
		return &Response{
			TreeResult: &result,
			Err:        err,
		}, nil
	}
}
```
