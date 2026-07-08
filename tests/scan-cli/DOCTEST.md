# Scan CLI

CLI tests for `disk-usage-analyser scan [PATH] [--json] [--threshold SIZE] [--max-depth N]`:
recursive directory tree with size annotations, threshold display filtering, depth limits,
in-memory `GlobalCache` deduplication, human-readable tree or nested JSON `TreeResult`, and
`run` dispatch before the default web-server branch.

## Version

0.0.2

# DSN (Domain Specific Notion)

The **scan command** walks a directory root recursively and builds a **TreeNode** tree.
Each node carries `name`, absolute `path`, recursive **Size** (files: `st_size`; directories:
sum of all nested content), `isDir`, `depth`, and **Children** sorted by size descending with
directories before files when sizes tie. **TreeResult** aggregates absolute **Path**,
**TotalSize** (full recursive bytes at root — includes nodes hidden by threshold or beyond
display depth), **Threshold**, **MaxDepth**, and the root **Tree** (name `"."`).

**Counting** always includes every file and directory. **Displaying** omits a node when
`size < threshold` (files and directories). **MaxDepth** limits branch expansion (`0` =
unlimited); ancestor sizes still include bytes from deeper levels even when children are not
shown.

**Scan** uses the shared **GlobalCache** / **CacheEntry** design: absolute path keys, progress
subscribers, and concurrent subtree deduplication. **RunCLI** resolves `PATH` (default: current
working directory), parses `--threshold` (default `1M`, compact binary sizes) and `--max-depth`
(default `3` text, `24` with `--json`), supports `--json` for one nested JSON object (no flat
`items` key), and prints a summary block (`PATH`, `TOTAL`, `THRESHOLD`, `MAX-DEPTH`), blank
line, then `tree(1)` box-drawing: name immediately after `├──`/`└──` (dirs with trailing `/`),
sizes in a left-aligned column after padding (no brackets; `FormatCompactHumanSize` values such
as `400B`, `1.1M`). Column alignment uses a two-pass layout: `maxLen(prefix+connector+name)` over
visible rows, pad each row to that width, then at least two spaces before the size. Root `.` has
no size. All user-facing stdout ends with a trailing blank line after the last content line.
**run.Run** wires `scan` before the web-server branch. Invalid threshold values exit non-zero
with a clear error.

## Decision Tree

```
scan-cli/
├── tree/                              # Scan API: TreeResult + ScanOptions
│   ├── basic-nested/                  # nested children with correct recursive sizes
│   ├── threshold-default/             # sub-1M nodes omitted from tree display
│   ├── threshold-override/            # --threshold 10M filters display
│   ├── max-depth-2/                   # branches stop at depth 2; sizes include deeper bytes
│   ├── max-depth-zero/                # maxDepth 0 = unlimited expansion
│   ├── max-depth-one/                 # only immediate children shown
│   └── root-total-includes-filtered/  # totalSize counts sub-threshold bytes
├── basics/                            # fundamental tree shapes (default scan opts)
│   ├── empty-dir/                     # empty root: no children, totalSize 0
│   ├── files-only/                    # two root files
│   ├── mixed/                         # root file + subdir with nested file
│   └── nested-dirs/                   # subdir aggregating nested subtree bytes
├── sorting/
│   └── by-size/                       # root children sorted size desc; dirs before files on tie
├── errors/
│   ├── missing-path/                  # non-existent PATH exits non-zero
│   ├── not-a-directory/               # regular file as PATH is rejected
│   └── invalid-threshold/             # --threshold foo exits non-zero
└── cli/
    ├── default-cwd/                   # no args: cwd scan; summary + tree lines
    ├── json-flag/                     # --json emits TreeResult JSON with nested tree
    ├── json-tree-shape/               # --json has tree/threshold/maxDepth; no items
    ├── text-tree-format/              # text summary + `.` + box-drawing + aligned size column
    ├── text-tree-alignment/           # nested fixture: sizes share one column across name lengths
    ├── help/                          # -h documents scan, PATH, --json
    ├── help-flags/                    # -h documents --threshold, --max-depth
    ├── explicit-path/                 # positional PATH scans that directory
    └── dispatch/                      # run.RunWithOptions routes scan without web server
```

## Test Index

| Leaf | Mode | Description |
|------|------|-------------|
| tree/basic-nested | scan | Nested `tree.children` with correct recursive sizes at multiple depths. |
| tree/threshold-default | scan | Default 1M threshold hides sub-1M nodes from display. |
| tree/threshold-override | scan | `Threshold=10M` hides nodes below 10M. |
| tree/max-depth-2 | scan | `MaxDepth=2` truncates branches; ancestor sizes include deeper bytes. |
| tree/max-depth-zero | scan | `MaxDepth=0` shows all levels (subject to threshold). |
| tree/max-depth-one | scan | `MaxDepth=1` shows only root immediate children. |
| tree/root-total-includes-filtered | scan | `totalSize` includes sub-threshold bytes not shown in tree. |
| basics/empty-dir | scan | Empty root: `tree.children` empty, `totalSize` 0, absolute `path`. |
| basics/files-only | scan | Two root files (100 B, 200 B): `totalSize` 300. |
| basics/mixed | scan | Root file 50 B + subdir with 150 B nested file: `totalSize` 200. |
| basics/nested-dirs | scan | `subdir/a` 100 B + `subdir/b/nested` 200 B: one dir child 300 B. |
| sorting/by-size | scan | Root children: descending size; dirs precede files when sizes tie. |
| errors/missing-path | cli | Non-existent PATH returns non-zero exit and error. |
| errors/not-a-directory | cli | Regular file as PATH returns error mentioning not a directory. |
| errors/invalid-threshold | cli | `--threshold foo` returns non-zero exit and clear error. |
| cli/default-cwd | cli | No args with cwd in fixture: summary lines + tree output. |
| cli/json-flag | cli | `--json` + fixture dir: valid JSON with nested `tree` and sizes. |
| cli/json-tree-shape | cli | `--json` includes `tree`, `threshold`, `maxDepth`; no `items`. |
| cli/text-tree-format | cli | Text has PATH/TOTAL/THRESHOLD/MAX-DEPTH, `.`, and name-then-aligned-size tree lines. |
| cli/text-tree-alignment | cli | Nested tree with varied name lengths; all size values start at the same column. |
| cli/help | cli | `-h` prints usage mentioning `scan`, `[PATH]`, `--json`. |
| cli/help-flags | cli | `-h` documents `--threshold` and `--max-depth`. |
| cli/explicit-path | cli | Explicit fixture path: PATH line matches absolute directory. |
| cli/dispatch | dispatch | `run.RunWithOptions` with `scan` never starts the web server. |

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