# Analyse CLI

CLI tests for `disk-usage-analyser analyse [DIR]`: portable directory-tree walk,
per-immediate-child TSV rows (`du -h -d 1` style), symlink and hardlink visibility,
summary line, always-on TSV header row, `pathfmt.Short` paths in TSV (JSON keeps
absolute paths), and `--json` output mode.

## Version

0.0.2

# DSN (Domain Specific Notion)

The **analyse command** is a CLI subcommand that walks a directory tree without
following symlinks (BSD `du -P` semantics). Output depth matches **`du -h -d 1`**:
one row per **immediate child** of the analysed root (real subdirectories and
symlinks; regular files in the root contribute only to the summary). Each child
row's metrics are the **subtree total** for that child path (real dirs: full
recursive walk; symlinks: do not follow — count the link and apply dir-symlink
symlink-count rules). Regular files contribute `st_size` **once per inode**;
unreadable entries are skipped. **DirResult** captures subtree metrics: apparent
unique bytes (`size`), symlink counts as `Ff+Dd`, extra hard-link references
(`Σ(nlink-1)`), hard-link inode bytes (`hardlink_size`), multiply-referenced
hard-link bytes (`shared_hardlink = Σ(size×nlink)` for `nlink>1`), and APFS
clone-family bytes (`shared_clone = Σ(size×ref_count)` for `doc_id` groups with
`ref_count>1` spanning >1 inode; hard-link inodes contribute only to
`shared_hardlink`), pnpm-store clone bytes (`pnpm_shared` = sum of unique
matching clone keys per `node_modules` walk, bubbled up to parent rows; store
index from `DISK_USAGE_ANALYSER_PNPM_STORE` or default
`~/Library/pnpm/store/v*/files`), and Bun install-cache clone bytes
(`bun_shared` = same Option A semantics against Bun's global cache; index from
`DISK_USAGE_ANALYSER_BUN_CACHE` or default `~/.bun/install/cache`; darwin
only in v1). **Result** holds `rows` (one per
immediate child, sorted by name) plus a **summary** row for the analysed root.
**RunCLI** formats TSV (always prints a header row, then data rows and a final
summary duplicate of root totals; TSV `path` column uses **pathfmt.Short** for
human-readable display while **JSON** keeps full absolute paths) or a single JSON
object; **Analyse** returns structured data for programmatic use.
The command is wired in **run** before the default web-server branch.

## Decision Tree

```
analyse-cli/
├── basics/
│   ├── empty-dir/           # no subdirs: summary only, zero bytes
│   └── regular-files/       # small files in root + one subdir, no links
├── output/
│   └── multi-level/         # nested subdirs: immediate children only + summary
├── symlinks/
│   └── file-and-dir/        # 1 file symlink + 1 dir symlink (no follow)
├── hardlinks/
│   ├── two-refs/            # one inode, nlink=2
│   └── three-refs/          # one inode, nlink=3
├── clones/
│   └── apfs-three-refs/     # darwin: cp -c original + 2 clones
├── pnpm-shared/
│   ├── store-match/
│   │   ├── node-modules-child/      # darwin: immediate node_modules row pnpm_shared=4K
│   │   └── nested-under-project/    # darwin: project row aggregates nested node_modules
│   ├── root-node-modules/
│   │   ├── summary-preserves-full/       # darwin: analyse node_modules root; summary pnpm_shared=4K, pkg row 4K; child sum = summary
│   │   └── breakdown-attributes-per-package/  # darwin: pkg row must carry store bytes; breakdown sum = summary
│   ├── nested-node-modules/
│   │   └── dedup-once/              # darwin: nested .pnpm node_modules must not double-count clone keys
│   ├── no-store-match/              # regular node_modules file, no clone match
│   ├── missing-store-env/           # darwin: clone but no store path → 0
│   └── non-darwin/                  # skip darwin; non-darwin always pnpm_shared=0
├── bun-shared/
│   ├── store-match/
│   │   ├── node-modules-child/      # darwin: immediate node_modules row bun_shared=4K
│   │   └── nested-under-project/    # darwin: project row aggregates nested node_modules
│   ├── root-node-modules/
│   │   ├── summary-preserves-full/       # darwin: analyse node_modules root; summary bun_shared=4K, pkg row 4K; child sum = summary
│   │   └── breakdown-attributes-per-package/  # darwin: pkg row must carry cache bytes; breakdown sum = summary
│   ├── nested-node-modules/
│   │   └── dedup-once/              # darwin: nested .pnpm node_modules must not double-count clone keys
│   ├── no-store-match/              # regular node_modules file, no clone match
│   ├── missing-store-env/           # darwin: clone but no cache path → 0
│   └── non-darwin/                  # skip darwin; non-darwin always bun_shared=0
├── combined/
│   └── mixed-tree/          # files + symlinks + hardlinks together
├── cli/
│   ├── default-cwd/         # No DIR: cwd fixture; TSV header + Short path "."
│   ├── header-flag/         # --header still accepted; TSV header always present
│   ├── json-flag/           # --json emits one JSON object
│   ├── help/                # -h/--help documents flags and exits 0
│   └── dispatch/            # run.Run dispatches analyse without web server
└── errors/
    └── missing-path/        # non-existent DIR exits 2
```

## Test Index

| Leaf | Mode | Description |
|------|------|-------------|
| basics/empty-dir | analyse | Empty root: no subdirectory rows, summary size 0, link cols 0. |
| basics/regular-files | analyse | Root file + subdir file: one subdir row + summary, no links. |
| output/multi-level | analyse | Two immediate-child rows (`alpha`, `gamma`); no grandchild rows. |
| symlinks/file-and-dir | analyse | `link-file` and `link-dir` rows; summary `1f+1d`, targets not followed. |
| hardlinks/two-refs | analyse | `hardlinks=1`, `hardlink_size=4K`, `shared_hardlink=8K`, `shared_clone=0`. |
| hardlinks/three-refs | analyse | `hardlinks=2`, `shared_hardlink=12K` (4K×3), `shared_clone=0`. |
| clones/apfs-three-refs | analyse | darwin: `cp -c` original + 2 clones; `shared_hardlink=0`, `shared_clone=size×3`. |
| pnpm-shared/store-match/node-modules-child | analyse | darwin: store clone in immediate `node_modules`; `pnpm_shared=4096`. |
| pnpm-shared/store-match/nested-under-project | analyse | darwin: nested `project/node_modules`; `project` row `pnpm_shared=4096`. |
| pnpm-shared/root-node-modules/summary-preserves-full | analyse | darwin: analyse `node_modules` root; summary `pnpm_shared=4096`, `pkg` row `pnpm_shared=4096`; child sum = summary. |
| pnpm-shared/root-node-modules/breakdown-attributes-per-package | analyse | darwin: analyse `node_modules` root; `pkg-a` row `pnpm_shared=4096`, `pkg-b` row `0`; child sum = summary. |
| pnpm-shared/nested-node-modules/dedup-once | analyse | darwin: outer + nested `.pnpm` `node_modules` same clone key; `pnpm_shared=4096` not 8192; `pnpm_shared <= size`. |
| pnpm-shared/no-store-match | analyse | `node_modules` regular file; store env set but no clone match; `pnpm_shared=0`. |
| pnpm-shared/missing-store-env | analyse | darwin: `cp -c` clone but no store path; `pnpm_shared=0`. |
| pnpm-shared/non-darwin | analyse | non-darwin: `node_modules` + store env; `pnpm_shared=0` (skip on darwin). |
| bun-shared/store-match/node-modules-child | analyse | darwin: cache clone in immediate `node_modules`; `bun_shared=4096`. |
| bun-shared/store-match/nested-under-project | analyse | darwin: nested `project/node_modules`; `project` row `bun_shared=4096`. |
| bun-shared/root-node-modules/summary-preserves-full | analyse | darwin: analyse `node_modules` root; summary `bun_shared=4096`, `pkg` row `bun_shared=4096`; child sum = summary. |
| bun-shared/root-node-modules/breakdown-attributes-per-package | analyse | darwin: analyse `node_modules` root; `pkg-a` row `bun_shared=4096`, `pkg-b` row `0`; child sum = summary. |
| bun-shared/nested-node-modules/dedup-once | analyse | darwin: outer + nested `.pnpm` `node_modules` same clone key; `bun_shared=4096` not 8192; `bun_shared <= size`. |
| bun-shared/no-store-match | analyse | `node_modules` regular file; cache env set but no clone match; `bun_shared=0`. |
| bun-shared/missing-store-env | analyse | darwin: `cp -c` clone but no cache path; `bun_shared=0`. |
| bun-shared/non-darwin | analyse | non-darwin: `node_modules` + cache env; `bun_shared=0` (skip on darwin). |
| combined/mixed-tree | analyse | `pkg` and `link-pkg` rows; hardlink bytes in `shared_hardlink` only. |
| cli/default-cwd | cli | No DIR: analyses cwd; TSV header + summary; path uses `pathfmt.Short` (`.`). |
| cli/header-flag | cli | TSV always has header row including `pnpm_shared` and `bun_shared`; `--header` remains valid (redundant). |
| cli/json-flag | cli | `--json` returns valid JSON with `root`, `rows`, `summary`, `pnpmSharedSize`, and `bunSharedSize` fields. |
| cli/help | cli | `-h` prints usage and documented flags; exit 0. |
| cli/dispatch | dispatch | `run.RunWithOptions` routes `analyse` without starting server. |
| errors/missing-path | cli | Missing path returns exit code 2. |

## How to Run

```sh
doctest vet ./tests/analyse-cli
doctest test ./tests/analyse-cli
```

```go
import (
	"bytes"
	"context"
	"strings"
	"testing"

	"disk-usage-analyser/analyse"
	"disk-usage-analyser/run"
)

type Request struct {
	Mode       string
	FixtureDir string
	Args       []string
	Stdout     *bytes.Buffer
	Stderr     *bytes.Buffer
}

type Response struct {
	Result           *analyse.Result
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
			Stdout:  req.Stdout,
			Stderr:  req.Stderr,
			StartServer: func(context.Context, run.ServerOptions) error {
				serverStarted = true
				t.Fatal("analyse dispatch must not start web server")
				return nil
			},
		})
		return &Response{
			Stdout:           req.Stdout.String(),
			Stderr:           req.Stderr.String(),
			Err:              err,
			ServerWasStarted: serverStarted,
		}, nil

	case "cli":
		exitCode, err := analyse.RunCLI(req.Args, analyse.CLIOptions{
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
		result, err := analyse.Analyse(req.FixtureDir)
		return &Response{
			Result: &result,
			Err:    err,
		}, nil
	}
}
```