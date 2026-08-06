# Go Best Practice Review — disk-usage-analyser

**Date:** 2026-08-06  
**Scope:** codebase structure, CLI design, flags handling, package layout  
**Reference recipes:** `go-best-practice` topics — `go-embed-assets`, `cli/*`, `flags-parsing/*`, `cmd-exec`, `kool-create`  
**Method:** read-only inspection of layout, entrypoints, flag parsers, embed setup, and external-command usage. No production code changes in this pass.

---

## Executive summary

The project is a mature **go + React** disk-usage tool with a strong secondary surface of **agent skill** packaging and a family of **node_modules** satellite CLIs. Several areas already align with the skill recipes:

| Area | Status |
|------|--------|
| Skill CLI (`skill --show` / `--install` / `--list`) | Strong — Shape 3 via `skillcmd.SingleSkill` |
| less-flags on main subcommands (`scan`, `analyse`, `explain`, `tmp-files`) | Mostly good |
| Dry-run on nm pipelines | Good one-path + side-effect gate |
| Streaming JSONL (`nmscan`, `nmcacheshared`) | Good |
| `script/{dev,build,install}` external commands | Good use of `xgo/support/cmd` |

The highest-impact gaps are **compile-breaking empty `//go:embed` of `dist/`**, **layering (`server` imported by library/CLI packages)**, and **CLI consistency** (stdlib `flag` vs less-flags, root dispatch UX, color API divergence).

Verified on this worktree:

```text
go build .
# main.go:12:12: pattern disk-usage-analyser-react/dist: no matching files found
```

---

## Project map (current layout)

```text
main.go                 # embed UI + server.Init + run.Run
run/                    # root CLI dispatch
analyse/                # analyse subcommand + APFS/pnpm/bun shared metrics
usagescan/              # scan / --inspect offline query
explain/                # reclaim-kind explain
tmpfiles/               # tmp-files scan library + CLI
server/                 # HTTP UI + disk APIs + PM detection + many scans (~4k LOC)
skill/                  # embedded SKILL.md + topics (Shape 3)
cmd/                    # satellite binaries (node-modules-*, probes)
nmscan|nminventory|nmcacheshared|nmmigrate|nmpipeline/
script/{dev,build,install}
disk-usage-analyser-react/
tests/                  # large doctest trees
```

This matches the spirit of a **kool-create `go-react`** product (Go backend + React frontend + `script/dev`), with substantial domain packages layered beside the server.

---

## Findings (severity-ordered)

### 1. Critical — `//go:embed dist` has no git-tracked content (breaks bare `go build` / `go install`)

**Topic:** `go-embed-assets` (Layer 1 always required)

**Evidence:**

- `main.go` embeds `disk-usage-analyser-react/dist` and `template.html`.
- `disk-usage-analyser-react/dist` is **missing** and **gitignored** (`disk-usage-analyser-react/.gitignore`: `dist`; root `.gitignore` also has `dist/`).
- No `placeholder.txt` (or any tracked file) under the embed root.
- Clean build fails: `pattern disk-usage-analyser-react/dist: no matching files found`.
- Fat path exists only via `script/build` + `script/install` (Layer 3 partial); **no** completeness check, **no** runtime hydrate (Layer 4), **no** release-asset story (Layer 5).

**Why it matters:** Module consumers, CI clean checkouts, and `go install <module>@…` cannot compile the main package without first running a Node frontend build.

**Recommended changes (grounded in `go-embed-assets`):**

1. **Layer 1 — always compile:**
   - Prefer a dedicated embed root (e.g. `embedded/ui/**` or keep `disk-usage-analyser-react/dist/**`) with a committed `placeholder.txt`.
   - Gitignore generated assets; un-ignore the placeholder:
     ```gitignore
     disk-usage-analyser-react/dist/**
     !disk-usage-analyser-react/dist/placeholder.txt
     ```
   - Use recursive embed: `//go:embed disk-usage-analyser-react/dist/**`.
2. **Layer 2 — completeness:** `EmbedComplete()` requiring real `index.html` + built JS (not only placeholder).
3. **Layer 3 — already partly present:** keep `script/build` / `script/install`; gate install with incompleteness auto-bundle (optional `--force-bundle`).
4. **Layer 4 (if bare `go install` should still serve UI):** version-pinned release tar hydrate + `assets status|ensure`.
5. Document in README: `go run ./script/install` for fat binary; bare build is thin/placeholder-only.

---

### 2. High — Domain / CLI packages import `server` (wrong package boundary)

**Topic:** package layout (implicit Go modularity; compounds CLI/library reuse)

**Evidence:**

| Importer | Uses from `server` |
|----------|--------------------|
| `nmmigrate/migrate.go` | `DetectHasPackageJSON`, `DetectGitTrackedPackageJSON` |
| `nmscan/record.go` | package-manager / enrichment helpers |
| `cmd/analyse-node-modules` | `AnalyseNodeModules` |
| `cmd/detect-package-manager` | `TracePackageManager`, detection helpers |

`server` is also the HTTP UI host (~27 Go files, ~4k lines): API handlers, SSE, disk manager, tmp scans, iTerm2 open, etc.

**Why it matters:** Library CLIs and migration tools cannot stay free of HTTP/UI concerns. Import graph pulls “server” for pure filesystem/git detection. Testing and extraction get harder; circular-risk grows as `server` already depends on `tmpfiles` / `analyse`.

**Recommended changes:**

1. Extract detection / analysis pure logic into domain packages, e.g.:
   - `pkg/packagemanager` or `internal/pm` — lockfile / PM detection  
   - `internal/gitmeta` or keep small helpers beside `nminventory` — git-tracked package.json  
   - keep `server` as thin HTTP adapters only  
2. Have `nmmigrate`, `nmscan`, and satellite `cmd/*` depend on those domain packages, **not** on `server`.
3. Longer term, split `server` by concern (`server/http`, `server/tmp`, `server/diskapi`) once pure logic is out.

---

### 3. High — Root CLI dispatch treats unknown words as scan roots / UI paths

**Topic:** `flags-parsing/subcommand` (dispatch + every-level help)

**Evidence (`run/run.go`):**

- Known verbs (`skill`, `install`, `scan`, `explain`, `analyse`, `tmp-files`) are handled first.
- Remaining args are parsed only for server flags (`--dev`, `--dev-idle-life`, `--component`).
- First remaining positional becomes `server.InitialDir` via `filepath.Abs`.
- Empty args start the web server (no automatic help).
- There is **no** “unknown command” path: a typo like `disk-usage-analyser scann` becomes “open UI with InitialDir=…/scann”.

**Contrast with recipe:** subcommand switch with explicit `default: unknown command`, and empty-args → level help when the binary is primarily a multi-command CLI.

**Recommended changes:**

1. Decide product default explicitly:
   - **A (tool CLI):** empty args and unknown command → print root help + non-zero on unknown.
   - **B (go-react app primary):** document that bare invocation starts the UI; still reject **unknown flags** and optionally require `serve` / `ui` as the server subcommand so typos are not silent paths.
2. Prefer a dedicated `serve` / `ui` subcommand for server mode; keep bare `disk-usage-analyser` as alias if desired, documented in help.
3. Optionally adopt `StopOnFirstArg()` if global flags ever coexist with subcommands (today globals are only server-mode flags after verb checks).
4. Keep the good pattern already present: each subpackage implements its own `-h/--help` with `HelpNoExit`.

---

### 4. Medium — Dual flag stacks: less-flags vs stdlib `flag`

**Topic:** `flags-parsing` (project standard = less-flags)

**Evidence:**

| Surface | Parser |
|---------|--------|
| `run`, `usagescan`, `analyse`, `explain`, `tmpfiles`, `nminventory`, `nmpipeline`, `nmscan` | `github.com/xhd2015/less-flags` |
| `cmd/analyse-node-modules` | stdlib `flag` (`-json`, `-quick`, `-v`, `-home`) |
| `cmd/detect-package-manager` | stdlib `flag` (`-home`, `-walk-up`) |

**Why it matters:** Inconsistent UX (`-json` vs `--json`, different help styles), harder shared testing, and satellite tools diverge from the house style already used in `nminventory.ParseFlags`.

**Recommended changes:**

1. Migrate both satellite mains (or better: move flag parse into library packages and keep `main` thin) to less-flags with long options (`--json`, `--home`, `--walk-up`, `-v,--verbose`).
2. Reuse shared helpers (`nminventory.ParseFlags` pattern or a small `internal/cliopts`) for home/json/verbose.
3. Keep `HelpFunc` + `HelpNoExit` for injectable stdout in tests (already a strength of the main CLIs).

---

### 5. Medium — Color API only on `explain`, and it diverges from `cli/color`

**Topic:** `cli/color`

**Evidence (`explain/run.go`):**

- Implements `--color=always|never|auto` (string), with bare `--color` rewritten to `--color=always` via hand-scanning argv (`normalizeColorArgs`).
- Auto path respects `NO_COLOR` and TTY (good).
- Does **not** implement the recipe’s `--color` / `--no-color` bool pair or mutual-exclusion error message.
- Other human CLIs (`analyse`, `scan` text tree, `tmp-files`) have no color policy.

**Recommended changes:**

1. Align `explain` with the recipe:

   ```text
   --color / --no-color  (mutually exclusive)
   auto: TTY && NO_COLOR empty
   ```

2. Prefer parsing bools with less-flags and rejecting conflict from parse results (do not hand-scan argv for `"--color"`).
3. Optional: extract shared `cli/color` helpers (`ColorMode`, `ResolveColor`, style wrap) for any future colored tables in `analyse` / `scan`.
4. Keep JSON output colorless (already true).

---

### 6. Medium — External commands: mixed `xgo/support/cmd` and raw `os/exec`

**Topic:** `cmd-exec`

**Good (use `cmd.Debug()` / `cmd.Output`):**

- `script/dev`, `script/build`, `script/install`
- `server/disk_manager.go`, `server/disk/disk.go`
- `server/tmp_runtime.go`, `tmp_simulator_runtime.go`, `tmp_podman_vm.go`

**Raw `os/exec` (production paths):**

| Location | Command |
|----------|---------|
| `server/server.go` `EnsureFrontendDevServer` | `bun run dev` |
| `server/detect_git_tracked.go` | `git ls-files` |
| `server/detect_belongs_to_git.go` | `git rev-parse` |
| `server/usage.go` | `osascript` |
| `nmmigrate/run.go` `ExecRunner` | `corepack use pnpm@latest` |

**Recommended changes:**

1. Prefer `github.com/xhd2015/xgo/support/cmd` for process spawn (Dir, Output capture, Debug prefix, consistent I/O).
2. Keep the **injected `CommandRunner` interface** in `nmmigrate` (excellent for dry-run tests); implement `ExecRunner` with `cmd` rather than `exec.Command` if Debug/Dir ergonomics help.
3. For silent probes (`git ls-files`), `cmd` with discarded stderr/stdout is fine; tests may keep raw `exec` for fixtures.

---

### 7. Medium — Satellite node_modules CLIs are outside the main binary

**Topic:** CLI UX / `flags-parsing/subcommand` / product cohesion

**Evidence:**

- Main binary subcommands: `analyse`, `scan`, `explain`, `tmp-files`, `skill`, `install`, (+ implicit server).
- Separate binaries under `cmd/`:
  - `node-modules-scan`
  - `node-modules-cache-shared`
  - `node-modules-migrate-pnpm`
  - `node-modules-migration-report`
  - `analyse-node-modules`, `detect-package-manager`
  - `tmp-named-scan-lifecycle-probe` (dev probe)

Logic packages (`nmscan`, …) are solid; only the **surface** is fragmented. README does not document these tools.

**Recommended changes:**

1. Either:
   - **Integrate** as `disk-usage-analyser node-modules <scan|cache-shared|migrate|report> …`, or  
   - **Document** them as intentional separate tools with install paths (`go run ./cmd/…`).
2. Keep thin `cmd/*/main.go` wrappers either way (good already for packages that use `RunCLI`).
3. Move probe-only binaries under `script/` or `internal/tools/` so they are not mistaken for product CLIs.

---

### 8. Medium — Duplicated human-size helpers

**Topic:** package layout / DRY (supporting CLI consistency)

**Evidence:** near-duplicate `ParseCompactHumanSize` / `FormatCompactHumanSize` in:

- `nminventory/humansize.go`
- `usagescan/humansize.go`

Plus additional formatters:

- `analyse.formatHumanSize`
- `tmpfiles.FormatHumanSize`
- `usagescan/format.go` helpers

Error strings already drift (`empty size string` vs `invalid min size: empty size string`).

**Recommended changes:**

1. Single package, e.g. `internal/humansize` or reuse one exported package and delete the twin.
2. All CLIs (`--min`, `--size-threshold`) call the same parser so doctests and UX stay MECE.

---

### 9. Low — `tmpfiles` is a single ~770-line package file

**Topic:** package layout

**Evidence:** `tmpfiles/tmpfiles.go` holds CLI parse, scan orchestration, binary classification, named hits, and formatting.

**Recommended changes (incremental):**

- `tmpfiles/cli.go` — `RunCLI` / flag parse  
- `tmpfiles/scan.go` — repo walk  
- `tmpfiles/binary.go` / `named.go` — hit types  
- Keep public API stable for `server` and tests.

---

### 10. Low — less-flags type usage: string duration instead of `Duration`

**Topic:** `flags-parsing/types`

**Evidence:** `run/run.go` parses `--dev-idle-life` as `string` + `time.ParseDuration`, with special cases `"off"` / `"0"`.

**Recommended changes:**

- If only Go durations are allowed, use `lessflags.Duration`.
- If `off` remains a product requirement, keep string parse but document it next to help (already partly documented); optional custom type later.

Also: usagescan correctly uses `**int` / `**string` for “unset vs zero” on `--max-depth` / `--min` / `--inspect` — **good**; treat as the house pattern for optional flags.

---

### 11. Low — Module path is bare `disk-usage-analyser`

**Topic:** `kool-create` / installability conventions

**Evidence:** `go.mod` → `module disk-usage-analyser` (not `github.com/xhd2015/disk-usage-analyser`). README clone URL uses GitHub.

**Why it matters:** `go install github.com/…@latest` module path will not match without a vanity or full module path; import paths in docs and skills assume the short name.

**Recommended changes:** Align module path with the public repo path when publishing; kool-create server templates rewrite module path from git remote for this reason.

---

### 12. Low — Exit-code fidelity at the root dispatcher

**Topic:** CLI UX

**Evidence:** `run.RunWithOptions` maps subcommand non-zero codes to:

```text
fmt.Errorf("scan exited with code %d", exitCode)
```

and `main` always `os.Exit(1)` on any error.

**Recommended changes:** Preserve exit codes 2 (usage) vs 1 (runtime) by returning a typed exit error or `(code int, err error)` from `run.Run`, matching what subpackages already return.

---

## What already matches go-best-practice well

### Skill CLI — `cli/skill-cli` Shape 3

`skill/skill.go` embeds `SKILL.md` + `android-images/` tree, uses `skillcmd.SingleSkill`, supports `--show` / `--install` / `--list`, both arg orders, top-level `install` alias, and nested topics as `TOPIC.md`. Help text lists topics. This is a **reference-quality** implementation relative to the recipe.

### Dry-run — `cli/dry-run`

`nmcacheshared`, `nmmigrate`, and `nmpipeline` share discovery/filter/limit via `nminventory`, then gate side effects:

- cache-shared: dry-run logs would-scan lines, skips counting  
- migrate: `migrateOne(..., dryRun)` sets plan fields, skips `RemoveAll` / corepack  
- report: dry-run measures before only, skips migrate  

One pipeline + `dryRun bool` — **correct**.

### Streaming — `cli/streaming`

`nmscan` / `nmcacheshared` emit JSONL as workers complete (with `bufio.Writer`), progress on stderr when verbose. Tables (`analyse`, migration report) correctly buffer for alignment — justified exception per recipe.

### Flags — `flags-parsing`

- Widespread less-flags with `HelpNoExit` + injectable writers (test-friendly).  
- `StringSlice` for `--name` / `--root` in `tmpfiles`.  
- Optional-pointer pattern in `usagescan` for inspect/min/depth.  
- Shared `nminventory.ParseFlags` for the nm family.

### Scripts / kool shape — `cmd-exec` + `kool-create`

`script/dev` → `cmd.Debug().Run("go", "run", ".", "--dev", …)`  
`script/build` → npm resolve + install + build via `cmd.Debug().Dir`  
`script/install` → build then `go install`  

This is the local fat-embed path for a go-react app. Completing embed Layer 1–2 closes the loop.

### Inspect offline workflow

`scan --inspect FILE` / live capture JSON is a strong CLI design for agent workflows (skill mentions offline inspect). Keep as a showcase of streaming-vs-buffer discipline (capture full JSON once; query offline).

---

## Recommended change backlog (priority)

| Priority | Change | Topics |
|----------|--------|--------|
| P0 | Placeholder + gitignore un-ignore under embed root; verify `go build` on clean tree | `go-embed-assets` |
| P0 | Document fat install (`go run ./script/install`) vs thin embed | `go-embed-assets`, README |
| P1 | Extract PM/git detection out of `server`; stop `nmmigrate`/`nmscan` importing `server` | package layout |
| P1 | Root CLI: unknown command handling; optional `serve` subcommand | `flags-parsing/subcommand` |
| P2 | Migrate `cmd/analyse-node-modules` and `detect-package-manager` to less-flags | `flags-parsing` |
| P2 | Align `explain` color flags with `cli/color` | `cli/color` |
| P2 | Replace remaining production `os/exec` with `xgo/support/cmd` where practical | `cmd-exec` |
| P2 | Unify humansize into one package | layout |
| P3 | Integrate or document satellite nm CLIs; relocate probes | CLI UX |
| P3 | Split `tmpfiles.go`; preserve exit codes; consider full module path | layout, CLI UX |
| P3 | Optional Layer 4 hydrate if bare `go install` must ship UI without Node | `go-embed-assets` |

---

## Suggested verification after fixes (not run in this review)

```bash
# Layer 1
go build -o /tmp/disk-usage-analyser .
go test ./usagescan ./analyse ./nminventory ./nmpipeline ./nmcacheshared ./nmmigrate ./nmscan ./explain ./run ...

# Fat path
go run ./script/install
disk-usage-analyser --help
disk-usage-analyser scan --help
disk-usage-analyser skill --list

# Doctest suites that already guard CLI contracts
# (existing tests/scan-cli, tests/explain-cli, tests/node-modules-cache-shared, …)
```

---

## Out of scope for this review

- Frontend React architecture and doctest design quality (separate skills).  
- Performance of APFS clone / pnpm shared walks.  
- Security review of `sudo mount` / password paths in disk manager.  
- Implementing the recommended fixes (except this docs-only report).

---

## Appendix — recipe checklist snapshot

| Recipe | Project fit |
|--------|-------------|
| `go-embed-assets` | Fail Layer 1; partial Layer 3 via scripts |
| `cli/skill-cli` | Pass Shape 3 |
| `cli/dry-run` | Pass on nm tools |
| `cli/streaming` | Pass on JSONL tools; tables OK buffered |
| `cli/color` | Partial / non-standard on `explain` only |
| `flags-parsing` | Strong core; gaps in root dispatch + 2 satellite mains |
| `flags-parsing/subcommand` | Manual dispatch OK; missing unknown-command / empty-args policy |
| `cmd-exec` | Good in scripts/disk; mixed in git/bun/corepack paths |
| `kool-create` | Structurally go-react-like; module path not vanity-aligned |

---

*End of review. Implementation of fixes should be staged starting at P0 embed compile safety, then package boundaries, then CLI consistency.*
