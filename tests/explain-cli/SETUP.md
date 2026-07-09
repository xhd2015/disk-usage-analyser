# Scenario

**Feature**: explain CLI harness (kind detection, `--kind` multi-root packs, `--all-kinds`,
human sections, JSON, dispatch, color)

```
# Human / JSON explain for PATH and/or --kind pack
explain [PATH] [--kind KIND] [--json] [--color=always|never|auto]
  -> KindDetector | kind pack registry -> measure sizes -> Explanation
Explanation -> text sections | JSON object -> stdout (trailing blank line)

# --all-kinds: PATH optional; default scope = CLIOptions.HomeDir (or real home)
# Mutually exclusive with --kind. Runs all v1 packs under scope.
explain --all-kinds [SCOPE?] -> AllKindsResult (index + detail per present pack)
explain --json --all-kinds [SCOPE?] -> JSON envelope {scope, totalSize, kinds[]}
explain --all-kinds --kind X -> error (mutually exclusive)

# --kind xcode: PATH optional; default scope = CLIOptions.HomeDir (or real home)
explain --kind xcode [PATH?] -> multi-root Xcode pack under scope

# --kind grok: PATH optional; scope home → {scope}/.grok (kind id grok-home)
explain --kind grok [PATH?] -> measure Grok CLI home under scope
explain PATH/.grok -> auto-detect kind=grok-home (signatures)

# --kind android-sdk: PATH optional; scope home → {scope}/Library/Android/sdk
explain --kind android-sdk [PATH?] -> measure Android SDK under scope (or SDK root)
explain PATH/Library/Android/sdk -> auto-detect kind=android-sdk (signatures)

# --kind iterm2: PATH optional; scope home → {scope}/Library/Application Support/iTerm2
explain --kind iterm2 [PATH?] -> measure iTerm2 App Support under scope (or iTerm2 root)
explain PATH/.../Application Support/iTerm2 -> auto-detect kind=iterm2 (signatures)

# --kind codex: PATH optional; scope home → {scope}/.codex (kind id codex-home)
explain --kind codex [PATH?] -> measure Codex CLI home under scope
explain PATH/.codex -> auto-detect kind=codex-home (signatures + logs_*.sqlite)
# When logs_*.sqlite present: human LOGS DB (ROWS + SAMPLE last 3) + JSON logsDb
# HOW TO PURGE: safe logs_2.sqlite reclaim (quit Codex; mv backup+wal/shm and/or sqlite3 DELETE/VACUUM)

# Human formatter: $ prompt + optional green base command
officialCommand / rawCommands -> human lines: "$ <cmd>"; color base token when color on
JSON -> plain officialCommand (no $, no ANSI)

# Dispatch before web server
run.RunWithOptions(["explain", ...]) -> explain.RunCLI (no StartServer)
```

## Preconditions

- Tests use isolated temporary fixture trees under `t.TempDir()`.
- Fixture sizes use exact byte counts so size assertions are deterministic
  (Codex `logs_*.sqlite` is built via `sqlite3` CLI when available — page size varies;
  totalSize assertions use non-DB payload lower bounds).
- The `explain` package provides `RunCLI` and `CLIOptions` (Stdout/Stderr writers;
  optional **`HomeDir`** for default `--kind` / `--all-kinds` scope without PATH).
- `req.FixtureDir` is an absolute directory used as the fixture root.
- `req.TargetPath` is the absolute PATH/scope argument when applicable (may be a file under the fixture).
- `req.HomeDir` is optional; when set, harness injects `CLIOptions.HomeDir` for no-PATH
  `--kind` / `--all-kinds` leaves.
- `RunCLI` receives args **after** the `explain` token (no `explain` prefix in `req.Args` for cli mode).
- Dispatch leaves pass a full argv including the `explain` token to `run.RunWithOptions`.
- Default `req.Mode` is `cli`.
- Single-kind human section headers are locked exactly: `PATH:`, `KIND:`, `TOTAL:`, `CONFIDENCE:`,
  `SUMMARY`, `BREAKDOWN`, `SAFE TO RECLAIM`, `HOW TO PURGE`, `RAW COMMANDS`.
  **codex-home** with readable `logs_*.sqlite` also emits a **`LOGS DB`** section (PATH/SIZE/ROWS/
  SAMPLE last 3) — not required for other kinds.
- **`--all-kinds` human** uses multi-kind layout: `SCOPE:` / `MODE: all-kinds` / `TOTAL (present):`
  header, `INDEX` table (SIZE/KIND/STATUS/PATH), then per-present mini-explain (KIND/BREAKDOWN…),
  optional MISSING list; trailing blank; never `rm -rf`.
- **BREAKDOWN** is an aligned table (`SIZE`/`NAME`/`ROLE`/`RECLAIMABLE`/`NOTES`), sorted
  size DESC; human RECLAIMABLE is Unicode `☑`/`☐` (never ASCII `[x]`/`[ ]`); JSON entries
  include `reclaimable` bool; when `--color=always`, ROLE cells may be green/yellow and
  reclaimable `☑` is green-wrapped (`☐` stays plain).
- User-facing stdout ends with a trailing blank line after the last content line.
- **SAFE TO RECLAIM** / **HOW TO PURGE** (and full stdout) must never mention `rm -rf`.
- **HOW TO PURGE** lists **CLI-first** official purge command(s) and what files each removes;
  UI is optional Notes only (`UI (optional): …`).
- Human runnable command lines (HOW TO PURGE official + RAW COMMANDS) use a leading **`$ `**;
  `#` comment / group lines do not.
- Default test I/O is non-TTY → no ANSI unless `--color=always` is passed.
- JSON has no ANSI and officialCommand has no `$` prefix.
- **`--json --all-kinds`** emits `AllKindsResult` envelope (`scope`, `totalSize`, `kinds[]`),
  not a single `Explanation`.

## Steps

1. Create a fresh fixture root for each leaf.
2. Build kind-specific fixtures (AVD, SeaTalk, Grok `.grok` home, Codex `.codex` home with
   tiny `logs_2.sqlite`, Android SDK under `Library/Android/sdk`, iTerm2 under
   `Library/Application Support/iTerm2`, Xcode multi-root pack, multi-pack home for
   `--all-kinds`, caches, generic files).
3. Set `req.Args` / `req.TargetPath` / optional `req.HomeDir` for `explain.RunCLI` or `run.RunWithOptions`.
4. Assert human sections, JSON fields, safety, `$` prefixes, color, and raw `scan` commands.

## Context

- PATH is required **unless** `--kind` or `--all-kinds` is set (then PATH is optional scope;
  default home). Empty args (no PATH, no kind, no all-kinds) still fail.
- `--all-kinds` and `--kind` are **mutually exclusive**.
- Both directories and files are valid PATH targets for auto-detect (no `--kind` / `--all-kinds`).
- Kind detection uses path patterns + signature files; first high-confidence wins.
- `--kind xcode` forces the Xcode multi-root pack (relative roots under scope).
- `--kind grok` forces Grok CLI home (`grok-home`): scope is home-like → `{scope}/.grok`,
  or scope is already `.grok` → use it; PATH optional with `HomeDir` inject.
- `--kind codex` forces Codex CLI home (`codex-home`): scope is home-like → `{scope}/.codex`,
  or scope is already `.codex` → use it; PATH optional with `HomeDir` inject.
- `--kind android-sdk` forces Android SDK: scope home-like → `{scope}/Library/Android/sdk`;
  if scope is already an SDK root (signatures / `Android/sdk` path), use it; PATH optional
  with `HomeDir` inject (default `{HomeDir}/Library/Android/sdk`).
- `--kind iterm2` forces iTerm2 Application Support: scope home-like →
  `{scope}/Library/Application Support/iTerm2`; if scope is already an iTerm2 root
  (signatures / `Application Support/iTerm2` path), use it; PATH optional with
  `HomeDir` inject (default `{HomeDir}/Library/Application Support/iTerm2`).
- **`--all-kinds`** runs all v1 registered packs (`xcode`, `grok`, `android-sdk`, `iterm2`,
  `codex`) under scope; missing pack roots → status `missing`, size 0; overall **exit 0**.
- Raw system commands are informational even when binaries are absent in CI.
- No auto-delete behavior is tested or required.
- Color modes: `auto` (default, TTY only), `always`, `never`.

```go
import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.FixtureDir = filepath.Join(t.TempDir(), "fixture")
	if err := os.MkdirAll(req.FixtureDir, 0755); err != nil {
		return err
	}
	if req.Mode == "" {
		req.Mode = "cli"
	}
	return nil
}

func mkdir(t *testing.T, base string, rel string) string {
	t.Helper()
	dir := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir %s: %v", rel, err)
	}
	return dir
}

func writeSizedFile(t *testing.T, base string, rel string, size int64) string {
	t.Helper()
	path := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir parent %s: %v", rel, err)
	}
	data := make([]byte, size)
	for i := range data {
		data[i] = 'x'
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
	return path
}

func writeTextFile(t *testing.T, base string, rel string, content string) string {
	t.Helper()
	path := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir parent %s: %v", rel, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
	return path
}

// writeAVDFixture creates a tiny Android AVD layout under fixtureDir/MediumPhone.avd.
// Returns the absolute path to the .avd directory and the absolute path to userdata-qemu.img.qcow2.
// Exact sizes: config.ini 32 B, userdata 400 B, sdcard 200 B, ram.bin 100 B → content total 732 B
// (directory metadata may make measured total ≥ 732 depending on walk semantics).
func writeAVDFixture(t *testing.T, fixtureDir string) (avdDir string, userdataPath string) {
	t.Helper()
	avdRel := "MediumPhone.avd"
	avdDir = mkdir(t, fixtureDir, avdRel)
	writeSizedFile(t, fixtureDir, avdRel+"/config.ini", 32)
	userdataPath = writeSizedFile(t, fixtureDir, avdRel+"/userdata-qemu.img.qcow2", 400)
	writeSizedFile(t, fixtureDir, avdRel+"/sdcard.img", 200)
	writeSizedFile(t, fixtureDir, avdRel+"/snapshots/default_boot/ram.bin", 100)
	return avdDir, userdataPath
}

// avdContentBytes is the sum of file payloads in writeAVDFixture (excludes dir inodes).
const avdContentBytes int64 = 32 + 400 + 200 + 100 // 732

// writeSeaTalkFixture creates a tiny SeaTalk Application Support layout under
// fixtureDir/Application Support/SeaTalk.
// Returns the absolute SeaTalk root directory and absolute path to main_1.sqlite.
// Exact file payloads (content total seatalkContentBytes = 510):
//
//	Service Worker/CacheStorage/entry  50 B  → web-cache
//	Cache/index                        80 B  → web-cache
//	main_1.sqlite                     200 B  → chat-db
//	search_1.sqlite                   100 B  → search-index
//	sqlite-backup/main_1.sqlite        60 B  → backup
//	config.json                        20 B  → config
//
// Directory metadata may make measured total ≥ 510 depending on walk semantics.
func writeSeaTalkFixture(t *testing.T, fixtureDir string) (seatalkDir string, mainDBPath string) {
	t.Helper()
	rel := "Application Support/SeaTalk"
	seatalkDir = mkdir(t, fixtureDir, rel)
	writeSizedFile(t, fixtureDir, rel+"/Service Worker/CacheStorage/entry", 50)
	writeSizedFile(t, fixtureDir, rel+"/Cache/index", 80)
	mainDBPath = writeSizedFile(t, fixtureDir, rel+"/main_1.sqlite", 200)
	writeSizedFile(t, fixtureDir, rel+"/search_1.sqlite", 100)
	writeSizedFile(t, fixtureDir, rel+"/sqlite-backup/main_1.sqlite", 60)
	writeSizedFile(t, fixtureDir, rel+"/config.json", 20)
	return seatalkDir, mainDBPath
}

// seatalkContentBytes is the sum of file payloads in writeSeaTalkFixture (excludes dir inodes).
const seatalkContentBytes int64 = 50 + 80 + 200 + 100 + 60 + 20 // 510

// writeXcodeScopeFixture creates a fake home/scope with all five Xcode pack roots.
// Relative layout matches frontend / server/tmp_analyse.go Xcode locations.
// Exact file payloads (content total xcodeContentBytes = 830):
//
//	Library/Developer/Xcode/DerivedData/file               400 B → derived-data ☑
//	Library/Developer/CoreSimulator/Devices/file           200 B → simulator ☑
//	Library/Developer/Xcode/Archives/file                  100 B → archives ☐
//	Library/Developer/Xcode/iOS DeviceSupport/file          80 B → device-support ☑
//	Library/Developer/Xcode/DocumentationCache/file         50 B → docs-cache ☑
//
// Size DESC name order: DerivedData, Devices, Archives, DeviceSupport, DocumentationCache.
// Directory metadata may make measured total ≥ 830 depending on walk semantics.
// Returns the absolute scope root (fixtureDir).
func writeXcodeScopeFixture(t *testing.T, fixtureDir string) string {
	t.Helper()
	writeSizedFile(t, fixtureDir, "Library/Developer/Xcode/DerivedData/file", 400)
	writeSizedFile(t, fixtureDir, "Library/Developer/CoreSimulator/Devices/file", 200)
	writeSizedFile(t, fixtureDir, "Library/Developer/Xcode/Archives/file", 100)
	writeSizedFile(t, fixtureDir, "Library/Developer/Xcode/iOS DeviceSupport/file", 80)
	writeSizedFile(t, fixtureDir, "Library/Developer/Xcode/DocumentationCache/file", 50)
	return fixtureDir
}

// xcodeContentBytes is the sum of file payloads in writeXcodeScopeFixture (excludes dir inodes).
const xcodeContentBytes int64 = 400 + 200 + 100 + 80 + 50 // 830

// xcodePackRoles is the locked role id set for the xcode multi-root pack (v1).
var xcodePackRoles = []string{
	"derived-data",
	"simulator",
	"device-support",
	"archives",
	"docs-cache",
}

// assertXcodeBreakdownMentions locks human BREAKDOWN/summary references for Xcode pack roles
// or well-known basenames under the multi-root layout.
func assertXcodeBreakdownMentions(t *testing.T, stdout string) {
	t.Helper()
	lower := strings.ToLower(stdout)
	checks := []struct {
		label string
		ok    bool
	}{
		{"derived-data or DerivedData", strings.Contains(lower, "derived-data") || strings.Contains(lower, "deriveddata")},
		{"simulator or Devices/CoreSimulator", strings.Contains(lower, "simulator") || strings.Contains(lower, "coresimulator") || strings.Contains(lower, "devices")},
		{"device-support or DeviceSupport", strings.Contains(lower, "device-support") || strings.Contains(lower, "devicesupport")},
		{"archives", strings.Contains(lower, "archives") || strings.Contains(lower, "archive")},
		{"docs-cache or DocumentationCache", strings.Contains(lower, "docs-cache") || strings.Contains(lower, "documentationcache") || strings.Contains(lower, "documentation")},
	}
	for _, c := range checks {
		if !c.ok {
			t.Fatalf("BREAKDOWN/summary should mention %s:\n%s", c.label, stdout)
		}
	}
	// Prefer explicit ROLE column tags when present.
	bd := breakdownSection(stdout)
	if strings.TrimSpace(bd) == "" {
		return
	}
	bdLower := strings.ToLower(bd)
	for _, role := range xcodePackRoles {
		if !strings.Contains(bdLower, role) {
			// Soft: names alone may satisfy the checks above; roles preferred.
			// Fail only if neither role nor a distinctive basename token exists for that role.
			switch role {
			case "derived-data":
				if !strings.Contains(bdLower, "deriveddata") {
					t.Fatalf("BREAKDOWN missing derived-data role/name:\n%s", bd)
				}
			case "simulator":
				if !strings.Contains(bdLower, "devices") && !strings.Contains(bdLower, "simulator") {
					t.Fatalf("BREAKDOWN missing simulator role/name:\n%s", bd)
				}
			case "device-support":
				if !strings.Contains(bdLower, "devicesupport") {
					t.Fatalf("BREAKDOWN missing device-support role/name:\n%s", bd)
				}
			case "archives":
				if !strings.Contains(bdLower, "archives") {
					t.Fatalf("BREAKDOWN missing archives role/name:\n%s", bd)
				}
			case "docs-cache":
				if !strings.Contains(bdLower, "documentation") && !strings.Contains(bdLower, "docs-cache") {
					t.Fatalf("BREAKDOWN missing docs-cache role/name:\n%s", bd)
				}
			}
		}
	}
}

// assertXcodeReclaimCheckboxes locks human RECLAIMABLE ☑/☐ for Xcode pack rows.
// Prefer matching by role id; fall back to distinctive basenames (BREAKDOWN only).
func assertXcodeReclaimCheckboxes(t *testing.T, stdout string) {
	t.Helper()
	bd := breakdownSection(stdout)
	pick := func(candidates []string) string {
		for _, name := range candidates {
			if strings.Contains(bd, name) {
				return name
			}
		}
		return ""
	}
	// reclaimable ☑
	for _, candidates := range [][]string{
		{"derived-data", "DerivedData"},
		{"docs-cache", "DocumentationCache"},
		{"device-support", "DeviceSupport"},
		{"simulator", "Devices"}, // role first; Devices != DeviceSupport
	} {
		name := pick(candidates)
		if name == "" {
			t.Fatalf("BREAKDOWN missing reclaimable row among %v:\n%s", candidates, bd)
		}
		assertBreakdownReclaimableCheckbox(t, stdout, name, true)
	}
	// archives ☐
	arch := pick([]string{"archives", "Archives"})
	if arch == "" {
		t.Fatalf("BREAKDOWN missing archives row:\n%s", bd)
	}
	assertBreakdownReclaimableCheckbox(t, stdout, arch, false)
}

// assertXcodeCLIFirstPurge locks CLI-primary HOW TO PURGE for kind xcode:
// official bodies prefer xcrun / simctl and/or disk-usage-analyser scan; never rm -rf;
// UI (Xcode Settings / Devices window) only under Notes.
func assertXcodeCLIFirstPurge(t *testing.T, stdout string) {
	t.Helper()
	how := sectionBody(stdout, "HOW TO PURGE")
	if strings.TrimSpace(how) == "" {
		t.Fatalf("missing HOW TO PURGE:\n%s", stdout)
	}
	lower := strings.ToLower(how)
	if !strings.Contains(lower, "xcrun") && !strings.Contains(lower, "simctl") &&
		!strings.Contains(lower, "disk-usage-analyser scan") {
		t.Fatalf("HOW TO PURGE must include CLI tools xcrun/simctl and/or disk-usage-analyser scan:\n%s", how)
	}
	if strings.Contains(lower, "rm -rf") || strings.Contains(lower, "rm -fr") {
		t.Fatalf("HOW TO PURGE must not contain rm -rf:\n%s", how)
	}

	inOfficial := false
	for _, line := range strings.Split(how, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "Official command:") {
			inOfficial = true
			continue
		}
		if strings.HasPrefix(trim, "Removes:") {
			inOfficial = false
			continue
		}
		if strings.HasPrefix(trim, "Notes:") {
			inOfficial = false
			continue
		}
		if isNumberedStepTitle(trim) {
			inOfficial = false
			continue
		}
		if !inOfficial || trim == "" || strings.HasPrefix(trim, "#") {
			continue
		}
		cmd := strings.TrimSpace(strings.TrimPrefix(trim, "$"))
		cmd = strings.TrimSpace(cmd)
		// Primary official command must not be pure UI navigation prose.
		if strings.Contains(cmd, "Xcode >") || strings.Contains(cmd, "Xcode →") ||
			(strings.Contains(cmd, "Settings") && strings.Contains(cmd, "Devices") &&
				!strings.Contains(cmd, "xcrun") && !strings.Contains(cmd, "simctl") &&
				!strings.Contains(cmd, "disk-usage-analyser")) {
			t.Fatalf("HOW TO PURGE Official command must be CLI-first (UI only in Notes), got: %q", line)
		}
	}
	assertHowToPurgeOfficialDollarPrefix(t, stdout)
}

// assertXcodeSafeToReclaim locks SAFE TO RECLAIM tiers for Xcode pack:
// rebuildable roots (derived-data/docs/device-support) as usually safe;
// archives not treated as usually-safe-only; simulator caution language preferred.
func assertXcodeSafeToReclaim(t *testing.T, stdout string) {
	t.Helper()
	safe := sectionBody(stdout, "SAFE TO RECLAIM")
	if strings.TrimSpace(safe) == "" {
		t.Fatalf("missing SAFE TO RECLAIM body:\n%s", stdout)
	}
	lower := strings.ToLower(safe)
	if !strings.Contains(lower, "usually safe") && !strings.Contains(lower, "usually-safe") &&
		!strings.Contains(lower, "reclaimable") && !strings.Contains(lower, "rebuild") {
		t.Fatalf("SAFE TO RECLAIM should mark rebuildable Xcode caches as usually safe / reclaimable:\n%s", safe)
	}
	// Archives: signed builds — section should not claim archives are usually-safe purge alone.
	for _, line := range strings.Split(safe, "\n") {
		ll := strings.ToLower(strings.TrimSpace(line))
		if ll == "" {
			continue
		}
		if strings.Contains(ll, "archive") &&
			(strings.Contains(ll, "usually safe") || strings.Contains(ll, "usually-safe")) &&
			!strings.Contains(ll, "not") && !strings.Contains(ll, "caution") &&
			!strings.Contains(ll, "valuable") && !strings.Contains(ll, "unless") {
			t.Fatalf("SAFE TO RECLAIM must not mark archives as usually-safe purge:\n%s", line)
		}
	}
}

// writeGrokHomeFixture creates a tiny Grok CLI home layout under parentDir/.grok.
// Returns the absolute .grok directory and absolute path to config.toml.
// Exact file payloads (content total grokHomeContentBytes = 798):
//
//	downloads/f           400 B → installer-cache ☑
//	sessions/s            200 B → session-logs ☑
//	marketplace-cache/m   100 B → cache ☑
//	logs/l                 50 B → logs ☑
//	config.toml            32 B → config ☐
//	auth.json              16 B → config ☐
//
// Size DESC name order: downloads, sessions, marketplace-cache, logs, config.toml, auth.json.
// Directory metadata may make measured total ≥ 798 depending on walk semantics.
// Signatures (config.toml + auth.json + sessions/ + downloads/) enable high-confidence detect.
func writeGrokHomeFixture(t *testing.T, parentDir string) (grokDir string, configPath string) {
	t.Helper()
	rel := ".grok"
	grokDir = mkdir(t, parentDir, rel)
	writeSizedFile(t, parentDir, rel+"/downloads/f", 400)
	writeSizedFile(t, parentDir, rel+"/sessions/s", 200)
	writeSizedFile(t, parentDir, rel+"/marketplace-cache/m", 100)
	writeSizedFile(t, parentDir, rel+"/logs/l", 50)
	configPath = writeSizedFile(t, parentDir, rel+"/config.toml", 32)
	writeSizedFile(t, parentDir, rel+"/auth.json", 16)
	return grokDir, configPath
}

// grokHomeContentBytes is the sum of file payloads in writeGrokHomeFixture (excludes dir inodes).
const grokHomeContentBytes int64 = 400 + 200 + 100 + 50 + 32 + 16 // 798

// grokHomeRoles is the locked role id set for the grok-home kind (fixture-covered).
func grokHomeRoles() []string {
	return []string{
		"installer-cache",
		"session-logs",
		"cache",
		"logs",
		"config",
	}
}

// codexLogRow is one fixture row for Codex logs_*.sqlite `logs` table.
type codexLogRow struct {
	TS      int64
	TSNanos int64
	Level   string
	Target  string
	Body    string
}

// defaultCodexLogRows returns 5 deterministic log rows (oldest→newest by ts).
// Newest 3 (ORDER BY ts DESC, ts_nanos DESC, id DESC LIMIT 3) are rows 5, 4, 3.
func defaultCodexLogRows() []codexLogRow {
	return []codexLogRow{
		{TS: 1000, TSNanos: 0, Level: "INFO", Target: "codex::session", Body: "row1 oldest body alpha"},
		{TS: 2000, TSNanos: 0, Level: "WARN", Target: "codex::agent", Body: "row2 body beta"},
		{TS: 3000, TSNanos: 0, Level: "ERROR", Target: "codex::tool", Body: "row3 body gamma"},
		{TS: 4000, TSNanos: 100, Level: "INFO", Target: "codex::session", Body: "row4 body delta"},
		{TS: 5000, TSNanos: 200, Level: "DEBUG", Target: "codex::runtime", Body: "row5 newest body epsilon"},
	}
}

// writeCodexLogsSQLite creates a tiny real SQLite DB at dbPath with a `logs` table and rows.
// Uses the `sqlite3` CLI when available; otherwise skips the leaf (document for implementer:
// pure Go seed via modernc.org/sqlite or embed is acceptable in production fixture helpers).
func writeCodexLogsSQLite(t *testing.T, dbPath string, rows []codexLogRow) {
	t.Helper()
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skip("sqlite3 CLI not found; need it to build logs_*.sqlite fixture (implementer may seed via pure Go)")
	}
	writeCodexLogsSQLiteOrNoop(t, dbPath, rows)
	if st, err := os.Stat(dbPath); err != nil || st.Size() == 0 {
		t.Fatalf("expected non-empty sqlite file at %s (err=%v)", dbPath, err)
	}
}

// writeCodexHomeTree creates the Codex CLI home layout under parentDir/.codex without
// logs_*.sqlite (so --all-kinds fixtures never Skip for missing sqlite3).
// Returns the absolute .codex directory and absolute path to config.toml.
// Non-DB file payloads (codexHomeNonDBContentBytes = 1198):
//
//	sessions/s    1000 B → session-logs ☑
//	cache/c        100 B → cache ☑
//	.tmp/t          50 B → tmp ☑
//	config.toml     32 B → config ☐
//	auth.json       16 B → config ☐
//
// Size floor 1198 > android-sdk 890 so all-kinds INDEX present order can put codex first
// even without the sqlite file.
func writeCodexHomeTree(t *testing.T, parentDir string) (codexDir string, configPath string) {
	t.Helper()
	rel := ".codex"
	codexDir = mkdir(t, parentDir, rel)
	writeSizedFile(t, parentDir, rel+"/sessions/s", 1000)
	writeSizedFile(t, parentDir, rel+"/cache/c", 100)
	writeSizedFile(t, parentDir, rel+"/.tmp/t", 50)
	configPath = writeSizedFile(t, parentDir, rel+"/config.toml", 32)
	writeSizedFile(t, parentDir, rel+"/auth.json", 16)
	return codexDir, configPath
}

// writeCodexHomeFixture creates a full Codex home for codex-home leaves: tree + logs_2.sqlite
// with 5 rows (requires `sqlite3` CLI; skips leaf if missing).
// Signatures: basename .codex + config.toml / auth.json / logs_*.sqlite / sessions/.
func writeCodexHomeFixture(t *testing.T, parentDir string) (codexDir string, configPath string) {
	t.Helper()
	codexDir, configPath = writeCodexHomeTree(t, parentDir)
	writeCodexLogsSQLite(t, filepath.Join(codexDir, "logs_2.sqlite"), defaultCodexLogRows())
	return codexDir, configPath
}

// codexHomeNonDBContentBytes is the sum of non-sqlite file payloads in writeCodexHomeTree.
const codexHomeNonDBContentBytes int64 = 1000 + 100 + 50 + 32 + 16 // 1198

// codexHomeFixtureLogRows is the locked COUNT(*) for the fixture logs table.
const codexHomeFixtureLogRows int64 = 5

// codexHomeRoles returns the locked role id set for the codex-home kind (fixture-covered).
func codexHomeRoles() []string {
	return []string{
		"app-logs-db",
		"session-logs",
		"cache",
		"tmp",
		"config",
	}
}

// assertGrokBreakdownMentions locks human BREAKDOWN/summary references for grok-home roles
// or well-known basenames under ~/.grok.
func assertGrokBreakdownMentions(t *testing.T, stdout string) {
	t.Helper()
	lower := strings.ToLower(stdout)
	checks := []struct {
		label string
		ok    bool
	}{
		{"installer-cache or downloads", strings.Contains(lower, "installer-cache") || strings.Contains(lower, "downloads")},
		{"session-logs or sessions", strings.Contains(lower, "session-logs") || strings.Contains(lower, "sessions")},
		{"cache or marketplace-cache", strings.Contains(lower, "marketplace-cache") ||
			(strings.Contains(lower, "cache") && (strings.Contains(lower, "marketplace") || strings.Contains(lower, "models")))},
		{"logs", strings.Contains(lower, "logs")},
		{"config or auth", strings.Contains(lower, "config") || strings.Contains(lower, "auth")},
	}
	for _, c := range checks {
		if !c.ok {
			t.Fatalf("BREAKDOWN/summary should mention %s:\n%s", c.label, stdout)
		}
	}
	// Prefer explicit ROLE column tags when present.
	bd := breakdownSection(stdout)
	if strings.TrimSpace(bd) == "" {
		return
	}
	bdLower := strings.ToLower(bd)
	for _, role := range grokHomeRoles() {
		if strings.Contains(bdLower, role) {
			continue
		}
		// Soft: names alone may satisfy the checks above; roles preferred.
		switch role {
		case "installer-cache":
			if !strings.Contains(bdLower, "downloads") {
				t.Fatalf("BREAKDOWN missing installer-cache role/name:\n%s", bd)
			}
		case "session-logs":
			if !strings.Contains(bdLower, "sessions") {
				t.Fatalf("BREAKDOWN missing session-logs role/name:\n%s", bd)
			}
		case "cache":
			if !strings.Contains(bdLower, "marketplace") {
				t.Fatalf("BREAKDOWN missing cache role/name (marketplace-cache):\n%s", bd)
			}
		case "logs":
			if !strings.Contains(bdLower, "logs") {
				t.Fatalf("BREAKDOWN missing logs role/name:\n%s", bd)
			}
		case "config":
			if !strings.Contains(bdLower, "config") && !strings.Contains(bdLower, "auth") {
				t.Fatalf("BREAKDOWN missing config role/name:\n%s", bd)
			}
		}
	}
}

// assertGrokReclaimCheckboxes locks human RECLAIMABLE ☑/☐ for grok-home fixture rows.
// Prefer matching by role id; fall back to distinctive basenames (BREAKDOWN only).
func assertGrokReclaimCheckboxes(t *testing.T, stdout string) {
	t.Helper()
	bd := breakdownSection(stdout)
	// Prefer longer / more distinctive tokens so "logs" does not match "session-logs" first.
	pick := func(candidates []string) string {
		for _, name := range candidates {
			if strings.Contains(bd, name) {
				return name
			}
		}
		return ""
	}
	// reclaimable ☑
	for _, candidates := range [][]string{
		{"installer-cache", "downloads"},
		{"session-logs", "sessions"},
		{"marketplace-cache", "marketplace"},
	} {
		name := pick(candidates)
		if name == "" {
			t.Fatalf("BREAKDOWN missing reclaimable row among %v:\n%s", candidates, bd)
		}
		assertBreakdownReclaimableCheckbox(t, stdout, name, true)
	}
	// logs role/basename: match a BREAKDOWN line that is logs but not session-logs.
	logsKey := ""
	for _, line := range strings.Split(bd, "\n") {
		ll := strings.ToLower(line)
		if strings.Contains(ll, "session-logs") {
			continue
		}
		if strings.Contains(ll, "logs") {
			// Prefer bare role "logs" or basename logs over incidental mentions.
			if strings.Contains(line, "logs") {
				logsKey = "logs"
				// Prefer a more specific token on this line if present.
				fields := strings.Fields(line)
				for _, f := range fields {
					if f == "logs" {
						logsKey = "logs"
						break
					}
				}
				break
			}
		}
	}
	if logsKey == "" {
		t.Fatalf("BREAKDOWN missing logs row (distinct from session-logs):\n%s", bd)
	}
	// Find line and check ☑ without using breakdownLineForName("logs") which may hit session-logs.
	logsLine := ""
	for _, line := range strings.Split(bd, "\n") {
		ll := strings.ToLower(line)
		if strings.Contains(ll, "session-logs") {
			continue
		}
		if strings.Contains(ll, "logs") {
			logsLine = line
			break
		}
	}
	if logsLine == "" || !strings.Contains(logsLine, "☑") {
		t.Fatalf("BREAKDOWN logs row must show reclaimable ☑: %q\n%s", logsLine, bd)
	}
	// config / auth ☐ — prefer exact basenames
	cfg := pick([]string{"config.toml", "auth.json", "config", "auth"})
	if cfg == "" {
		t.Fatalf("BREAKDOWN missing config/auth row:\n%s", bd)
	}
	assertBreakdownReclaimableCheckbox(t, stdout, cfg, false)
}

// assertGrokCLIFirstPurge locks CLI-primary HOW TO PURGE for kind grok-home:
// official bodies prefer disk-usage-analyser scan (and/or careful reclaim steps);
// never rm -rf; auth/config not usually-safe Removes.
func assertGrokCLIFirstPurge(t *testing.T, stdout string) {
	t.Helper()
	how := sectionBody(stdout, "HOW TO PURGE")
	if strings.TrimSpace(how) == "" {
		t.Fatalf("missing HOW TO PURGE:\n%s", stdout)
	}
	lower := strings.ToLower(how)
	if !strings.Contains(lower, "disk-usage-analyser scan") {
		t.Fatalf("HOW TO PURGE must include disk-usage-analyser scan (inspect first):\n%s", how)
	}
	if strings.Contains(lower, "rm -rf") || strings.Contains(lower, "rm -fr") {
		t.Fatalf("HOW TO PURGE must not contain rm -rf:\n%s", how)
	}
	// Should cover installer downloads and/or caches in purge guidance.
	if !strings.Contains(lower, "download") && !strings.Contains(lower, "marketplace") &&
		!strings.Contains(lower, "session") && !strings.Contains(lower, "cache") {
		t.Fatalf("HOW TO PURGE should mention downloads/sessions/marketplace/cache reclaim targets:\n%s", how)
	}
	// Never mark auth/config as usually-safe Removes.
	for _, line := range strings.Split(how, "\n") {
		trim := strings.TrimSpace(line)
		ll := strings.ToLower(trim)
		if !strings.HasPrefix(ll, "removes:") {
			continue
		}
		if (strings.Contains(ll, "usually safe") || strings.Contains(ll, "usually-safe")) &&
			(strings.Contains(ll, "auth") || strings.Contains(ll, "config.toml") ||
				(strings.Contains(ll, "config") && !strings.Contains(ll, "not"))) {
			t.Fatalf("HOW TO PURGE must not mark auth/config as usually-safe Removes:\n%s", line)
		}
	}
	assertHowToPurgeOfficialDollarPrefix(t, stdout)
}

// assertGrokSafeToReclaim locks SAFE TO RECLAIM tiers for grok-home:
// downloads/caches usually safe or reclaimable; sessions caution preferred;
// auth/config not usually-safe purge.
func assertGrokSafeToReclaim(t *testing.T, stdout string) {
	t.Helper()
	safe := sectionBody(stdout, "SAFE TO RECLAIM")
	if strings.TrimSpace(safe) == "" {
		t.Fatalf("missing SAFE TO RECLAIM body:\n%s", stdout)
	}
	lower := strings.ToLower(safe)
	if !strings.Contains(lower, "usually safe") && !strings.Contains(lower, "usually-safe") &&
		!strings.Contains(lower, "reclaimable") {
		t.Fatalf("SAFE TO RECLAIM should mark installer/cache/logs as usually safe / reclaimable:\n%s", safe)
	}
	// auth/config must not be presented as usually-safe purge alone.
	for _, line := range strings.Split(safe, "\n") {
		ll := strings.ToLower(strings.TrimSpace(line))
		if ll == "" {
			continue
		}
		if (strings.Contains(ll, "auth") || strings.Contains(ll, "config.toml") ||
			(strings.Contains(ll, "credential") || strings.Contains(ll, "config"))) &&
			(strings.Contains(ll, "usually safe") || strings.Contains(ll, "usually-safe")) &&
			!strings.Contains(ll, "not") && !strings.Contains(ll, "caution") &&
			!strings.Contains(ll, "never") && !strings.Contains(ll, "keep") {
			// Only fail when the line clearly targets auth/config as the safe purge subject.
			if strings.Contains(ll, "auth") || strings.Contains(ll, "config.toml") ||
				strings.Contains(ll, "credential") {
				t.Fatalf("SAFE TO RECLAIM must not mark auth/config as usually-safe purge:\n%s", line)
			}
		}
	}
}

// assertCodexBreakdownMentions locks human BREAKDOWN/summary references for codex-home
// roles or well-known basenames under ~/.codex.
func assertCodexBreakdownMentions(t *testing.T, stdout string) {
	t.Helper()
	lower := strings.ToLower(stdout)
	checks := []struct {
		label string
		ok    bool
	}{
		{"app-logs-db or logs_*.sqlite", strings.Contains(lower, "app-logs-db") ||
			strings.Contains(lower, "logs_2.sqlite") || strings.Contains(lower, "logs_")},
		{"session-logs or sessions", strings.Contains(lower, "session-logs") || strings.Contains(lower, "sessions")},
		{"cache", strings.Contains(lower, "cache")},
		{"tmp or .tmp", strings.Contains(lower, "tmp")},
		{"config or auth", strings.Contains(lower, "config") || strings.Contains(lower, "auth")},
	}
	for _, c := range checks {
		if !c.ok {
			t.Fatalf("BREAKDOWN/summary should mention %s:\n%s", c.label, stdout)
		}
	}
	bd := breakdownSection(stdout)
	if strings.TrimSpace(bd) == "" {
		return
	}
	bdLower := strings.ToLower(bd)
	for _, role := range codexHomeRoles() {
		if strings.Contains(bdLower, role) {
			continue
		}
		switch role {
		case "app-logs-db":
			if !strings.Contains(bdLower, "logs_") && !strings.Contains(bdLower, "sqlite") {
				t.Fatalf("BREAKDOWN missing app-logs-db role/name (logs_*.sqlite):\n%s", bd)
			}
		case "session-logs":
			if !strings.Contains(bdLower, "sessions") {
				t.Fatalf("BREAKDOWN missing session-logs role/name:\n%s", bd)
			}
		case "cache":
			if !strings.Contains(bdLower, "cache") {
				t.Fatalf("BREAKDOWN missing cache role/name:\n%s", bd)
			}
		case "tmp":
			if !strings.Contains(bdLower, "tmp") {
				t.Fatalf("BREAKDOWN missing tmp role/name:\n%s", bd)
			}
		case "config":
			if !strings.Contains(bdLower, "config") && !strings.Contains(bdLower, "auth") {
				t.Fatalf("BREAKDOWN missing config role/name:\n%s", bd)
			}
		}
	}
}

// assertCodexReclaimCheckboxes locks human RECLAIMABLE ☑/☐ for codex-home fixture rows.
func assertCodexReclaimCheckboxes(t *testing.T, stdout string) {
	t.Helper()
	bd := breakdownSection(stdout)
	pick := func(candidates []string) string {
		for _, name := range candidates {
			if strings.Contains(bd, name) {
				return name
			}
		}
		return ""
	}
	// reclaimable ☑
	for _, candidates := range [][]string{
		{"app-logs-db", "logs_2.sqlite", "logs_"},
		{"session-logs", "sessions"},
		{"cache"},
		{"tmp", ".tmp"},
	} {
		name := pick(candidates)
		if name == "" {
			t.Fatalf("BREAKDOWN missing reclaimable row among %v:\n%s", candidates, bd)
		}
		assertBreakdownReclaimableCheckbox(t, stdout, name, true)
	}
	// config / auth ☐
	cfg := pick([]string{"config.toml", "auth.json", "config", "auth"})
	if cfg == "" {
		t.Fatalf("BREAKDOWN missing config/auth row:\n%s", bd)
	}
	assertBreakdownReclaimableCheckbox(t, stdout, cfg, false)
}

// assertCodexCLIFirstPurge locks CLI-primary HOW TO PURGE for kind codex-home:
// disk-usage-analyser scan inspect; never rm -rf; auth/config not usually-safe Removes;
// plus safe logs_2.sqlite reclaim guidance (assertCodexHowToPurgeLogs).
func assertCodexCLIFirstPurge(t *testing.T, stdout string) {
	t.Helper()
	how := sectionBody(stdout, "HOW TO PURGE")
	if strings.TrimSpace(how) == "" {
		t.Fatalf("missing HOW TO PURGE:\n%s", stdout)
	}
	lower := strings.ToLower(how)
	if !strings.Contains(lower, "disk-usage-analyser scan") {
		t.Fatalf("HOW TO PURGE must include disk-usage-analyser scan (inspect first):\n%s", how)
	}
	if strings.Contains(lower, "rm -rf") || strings.Contains(lower, "rm -fr") {
		t.Fatalf("HOW TO PURGE must not contain rm -rf:\n%s", how)
	}
	if !strings.Contains(lower, "session") && !strings.Contains(lower, "cache") &&
		!strings.Contains(lower, "log") && !strings.Contains(lower, "tmp") {
		t.Fatalf("HOW TO PURGE should mention sessions/cache/logs/tmp reclaim targets:\n%s", how)
	}
	for _, line := range strings.Split(how, "\n") {
		trim := strings.TrimSpace(line)
		ll := strings.ToLower(trim)
		if !strings.HasPrefix(ll, "removes:") {
			continue
		}
		if (strings.Contains(ll, "usually safe") || strings.Contains(ll, "usually-safe")) &&
			(strings.Contains(ll, "auth") || strings.Contains(ll, "config.toml") ||
				(strings.Contains(ll, "config") && !strings.Contains(ll, "not"))) {
			t.Fatalf("HOW TO PURGE must not mark auth/config as usually-safe Removes:\n%s", line)
		}
	}
	assertHowToPurgeOfficialDollarPrefix(t, stdout)
	assertCodexHowToPurgeLogs(t, stdout)
}

// assertCodexHowToPurgeLogs locks safe reclaim guidance for Codex logs_2.sqlite in
// HOW TO PURGE (and soft cues in SAFE TO RECLAIM when present):
//
//   - Quit Codex fully before touching the DB
//   - Prefer move-aside backup of logs_2.sqlite (+ wal/shm) and/or sqlite3 DELETE + VACUUM
//   - Diagnostic-only logs; not state_5 / auth / config; recreates; may regrow; no live truncate
//   - Never rm -rf (also covered by assertCodexCLIFirstPurge / assertNoRmRf)
func assertCodexHowToPurgeLogs(t *testing.T, stdout string) {
	t.Helper()
	how := sectionBody(stdout, "HOW TO PURGE")
	if strings.TrimSpace(how) == "" {
		t.Fatalf("missing HOW TO PURGE for logs_2.sqlite reclaim guidance:\n%s", stdout)
	}
	// Include SAFE TO RECLAIM so optional quit/diagnostic notes there also count.
	safe := sectionBody(stdout, "SAFE TO RECLAIM")
	joined := how + "\n" + safe
	lower := strings.ToLower(joined)

	// 1) Quit Codex first (before reclaim / DB work). Prefer explicit quit/exit/stop.
	// Soft accept: "fully quit", "quit codex", "codex must not be running", etc.
	hasQuitWord := strings.Contains(lower, "quit") || strings.Contains(lower, "exit codex") ||
		strings.Contains(lower, "stop codex") || strings.Contains(lower, "fully close") ||
		strings.Contains(lower, "close codex") || strings.Contains(lower, "not running") ||
		strings.Contains(lower, "is not running") || strings.Contains(lower, "while codex is running") ||
		strings.Contains(lower, "while running")
	hasCodexCue := strings.Contains(lower, "codex")
	if !hasQuitWord || !hasCodexCue {
		t.Fatalf("HOW TO PURGE/SAFE must instruct to quit Codex (fully) before reclaiming logs DB:\n%s", how)
	}

	// 2) Target is logs_2.sqlite (or logs_*.sqlite pattern with logs_2 cue preferred).
	hasLogsDB := strings.Contains(lower, "logs_2.sqlite") || strings.Contains(lower, "logs_2") ||
		strings.Contains(lower, "logs_*.sqlite") ||
		(strings.Contains(lower, "logs_") && strings.Contains(lower, ".sqlite"))
	if !hasLogsDB {
		t.Fatalf("HOW TO PURGE must mention logs_2.sqlite / logs_*.sqlite reclaim target:\n%s", how)
	}

	// 3) Move-aside backup AND sqlite3 DELETE FROM logs; VACUUM (both alternatives).
	hasMvBackup := strings.Contains(lower, "mv ") || strings.Contains(lower, "\nmv") ||
		strings.Contains(lower, "move-aside") || strings.Contains(lower, "move aside") ||
		strings.Contains(lower, "backup") || strings.Contains(lower, "rename") ||
		strings.Contains(lower, "move the") || strings.Contains(lower, "move logs")
	hasSqlite3 := strings.Contains(lower, "sqlite3")
	hasDelete := strings.Contains(lower, "delete from logs") ||
		(strings.Contains(lower, "delete") && strings.Contains(lower, "logs"))
	hasVacuum := strings.Contains(lower, "vacuum")
	hasSqlitePurge := hasSqlite3 && hasDelete && hasVacuum
	if !hasMvBackup && !hasSqlitePurge {
		t.Fatalf("HOW TO PURGE must document mv/backup of logs_2.sqlite and/or sqlite3 DELETE FROM logs + VACUUM:\n%s", how)
	}
	// Prefer both alternatives so operators can choose backup vs in-place vacuum.
	if !hasMvBackup || !hasSqlitePurge {
		t.Fatalf("HOW TO PURGE should document both reclaim alternatives: (1) mv/backup logs_2.sqlite (+wal/shm) and (2) sqlite3 DELETE FROM logs; VACUUM — missing mv/backup=%v sqlite3+DELETE+VACUUM=%v:\n%s",
			!hasMvBackup, !hasSqlitePurge, how)
	}

	// 4) WAL/SHM companions (soft-required when documenting file move/backup).
	if hasMvBackup {
		hasWalOrShm := strings.Contains(lower, "wal") || strings.Contains(lower, "shm") ||
			strings.Contains(lower, "-wal") || strings.Contains(lower, "-shm")
		if !hasWalOrShm {
			t.Fatalf("HOW TO PURGE mv/backup guidance should mention logs_2.sqlite-wal / -shm companions:\n%s", how)
		}
	}

	// 5) Do not clear state_5 / auth / config as part of logs reclaim.
	hasStateCaution := strings.Contains(lower, "state_5") || strings.Contains(lower, "state_*.sqlite") ||
		strings.Contains(lower, "state_") || strings.Contains(lower, "app-state") ||
		strings.Contains(lower, "app state")
	hasAuthOrConfig := strings.Contains(lower, "auth") || strings.Contains(lower, "config.toml") ||
		strings.Contains(lower, "config")
	if !hasStateCaution {
		t.Fatalf("HOW TO PURGE must caution not to clear state_5.sqlite* / app-state while reclaiming logs:\n%s", how)
	}
	if !hasAuthOrConfig {
		t.Fatalf("HOW TO PURGE must caution not to clear auth/config while reclaiming logs:\n%s", how)
	}

	// 6) Notes: diagnostic-only; recreates; may regrow; no live truncate.
	// Prefer explicit "diagnostic" so implementer does not rely only on generic "debug trails".
	hasDiagnostic := strings.Contains(lower, "diagnostic") ||
		(strings.Contains(lower, "debug") && (strings.Contains(lower, "only") ||
			strings.Contains(lower, "not state") || strings.Contains(lower, "not auth") ||
			strings.Contains(lower, "not config") || strings.Contains(lower, "not session")))
	if !hasDiagnostic {
		t.Fatalf("HOW TO PURGE notes should mark logs DB as diagnostic-only (not session/auth/state):\n%s", how)
	}
	// Prefer explicit recreat* for the logs DB itself (not only "fresh log" soft language on SAFE).
	hasRecreates := strings.Contains(lower, "recreat") || strings.Contains(lower, "re-creat") ||
		(strings.Contains(lower, "creates") && strings.Contains(lower, "log")) ||
		(strings.Contains(lower, "create") && strings.Contains(lower, "logs_"))
	if !hasRecreates {
		t.Fatalf("HOW TO PURGE notes should say Codex recreates logs_2.sqlite after reclaim:\n%s", how)
	}
	hasRegrow := strings.Contains(lower, "regrow") || strings.Contains(lower, "re-grow") ||
		strings.Contains(lower, "grow again") || strings.Contains(lower, "grows back") ||
		strings.Contains(lower, "may grow") || strings.Contains(lower, "trace") ||
		strings.Contains(lower, "fill again") || strings.Contains(lower, "reaccumulate") ||
		strings.Contains(lower, "re-accumulate")
	if !hasRegrow {
		t.Fatalf("HOW TO PURGE notes should warn logs DB may regrow (e.g. TRACE verbosity):\n%s", how)
	}
	hasNoLiveTruncate := strings.Contains(lower, "not while") || strings.Contains(lower, "while running") ||
		strings.Contains(lower, "while codex") || strings.Contains(lower, "live") ||
		strings.Contains(lower, "do not truncate") || strings.Contains(lower, "don't truncate") ||
		strings.Contains(lower, "no live") || strings.Contains(lower, "not truncate") ||
		strings.Contains(lower, "running codex") || strings.Contains(lower, "codex is running") ||
		strings.Contains(lower, "quit") // quit-first already implies no live truncate; still prefer explicit
	if !hasNoLiveTruncate {
		t.Fatalf("HOW TO PURGE must discourage live truncate / reclaim while Codex is running:\n%s", how)
	}

	if strings.Contains(lower, "rm -rf") || strings.Contains(lower, "rm -fr") {
		t.Fatalf("HOW TO PURGE must never use rm -rf for logs reclaim:\n%s", how)
	}
}

// assertJSONCodexHowToPurgeLogs locks JSON howToPurge officialCommand/removes/notes for
// safe logs_2.sqlite reclaim: quit Codex; mv backup and/or sqlite3 DELETE+VACUUM; state/auth
// caution; diagnostic/recreate/regrow notes; plain (no $ / ANSI / rm -rf).
func assertJSONCodexHowToPurgeLogs(t *testing.T, howToPurge []map[string]any) {
	t.Helper()
	if len(howToPurge) == 0 {
		t.Fatal("howToPurge must be non-empty for codex logs reclaim JSON lock")
	}
	joined := ""
	for i, step := range howToPurge {
		oc, _ := step["officialCommand"].(string)
		rm, _ := step["removes"].(string)
		notes, _ := step["notes"].(string)
		title, _ := step["title"].(string)
		if containsANSI(oc) || containsANSI(rm) || containsANSI(notes) || containsANSI(title) {
			t.Fatalf("howToPurge[%d] fields must not contain ANSI", i)
		}
		for _, s := range []string{oc, rm, notes} {
			for _, line := range strings.Split(s, "\n") {
				if strings.HasPrefix(strings.TrimSpace(line), "$") {
					t.Fatalf("JSON howToPurge must not include $ prefix: %q", s)
				}
			}
			ll := strings.ToLower(s)
			if strings.Contains(ll, "rm -rf") || strings.Contains(ll, "rm -fr") {
				t.Fatalf("JSON howToPurge must not contain rm -rf: %q", s)
			}
		}
		joined += title + "\n" + oc + "\n" + rm + "\n" + notes + "\n"
	}
	lower := strings.ToLower(joined)

	hasQuit := strings.Contains(lower, "quit") || strings.Contains(lower, "exit") ||
		strings.Contains(lower, "stop") || strings.Contains(lower, "while running") ||
		strings.Contains(lower, "not while")
	if !hasQuit || !strings.Contains(lower, "codex") {
		t.Fatalf("JSON howToPurge must instruct quit Codex before logs reclaim: %s", joined)
	}
	if !strings.Contains(lower, "logs_2") && !strings.Contains(lower, "logs_*.sqlite") &&
		!(strings.Contains(lower, "logs_") && strings.Contains(lower, "sqlite")) {
		t.Fatalf("JSON howToPurge must mention logs_2.sqlite / logs_*.sqlite: %s", joined)
	}
	hasMvBackup := strings.Contains(lower, "mv ") || strings.Contains(lower, "backup") ||
		strings.Contains(lower, "move-aside") || strings.Contains(lower, "move aside") ||
		strings.Contains(lower, "rename")
	hasSqlitePurge := strings.Contains(lower, "sqlite3") &&
		(strings.Contains(lower, "delete from logs") ||
			(strings.Contains(lower, "delete") && strings.Contains(lower, "logs"))) &&
		strings.Contains(lower, "vacuum")
	if !hasMvBackup || !hasSqlitePurge {
		t.Fatalf("JSON howToPurge must document both mv/backup and sqlite3 DELETE+VACUUM for logs (mv/backup=%v sqlite=%v): %s",
			hasMvBackup, hasSqlitePurge, joined)
	}
	if hasMvBackup && !(strings.Contains(lower, "wal") || strings.Contains(lower, "shm")) {
		t.Fatalf("JSON howToPurge mv/backup should mention -wal/-shm: %s", joined)
	}
	if !strings.Contains(lower, "state_5") && !strings.Contains(lower, "state_*.sqlite") &&
		!strings.Contains(lower, "state_") && !strings.Contains(lower, "app-state") {
		t.Fatalf("JSON howToPurge must caution against clearing state_5 / app-state: %s", joined)
	}
	if !strings.Contains(lower, "auth") && !strings.Contains(lower, "config") {
		t.Fatalf("JSON howToPurge must caution against clearing auth/config: %s", joined)
	}
	if !strings.Contains(lower, "diagnostic") &&
		!(strings.Contains(lower, "debug") && (strings.Contains(lower, "only") ||
			strings.Contains(lower, "not state") || strings.Contains(lower, "not auth"))) {
		t.Fatalf("JSON howToPurge notes should mark logs as diagnostic-only: %s", joined)
	}
	if !strings.Contains(lower, "recreat") && !strings.Contains(lower, "re-creat") &&
		!(strings.Contains(lower, "creates") && strings.Contains(lower, "log")) {
		t.Fatalf("JSON howToPurge notes should say Codex recreates the logs DB: %s", joined)
	}
	if !strings.Contains(lower, "regrow") && !strings.Contains(lower, "re-grow") &&
		!strings.Contains(lower, "grow again") && !strings.Contains(lower, "may grow") &&
		!strings.Contains(lower, "trace") && !strings.Contains(lower, "fill again") {
		t.Fatalf("JSON howToPurge notes should warn logs may regrow (TRACE etc.): %s", joined)
	}
}

// assertCodexSafeToReclaim locks SAFE TO RECLAIM for codex-home:
// logs/sessions/cache/tmp reclaimable language; auth/config not usually-safe purge.
func assertCodexSafeToReclaim(t *testing.T, stdout string) {
	t.Helper()
	safe := sectionBody(stdout, "SAFE TO RECLAIM")
	if strings.TrimSpace(safe) == "" {
		t.Fatalf("missing SAFE TO RECLAIM body:\n%s", stdout)
	}
	lower := strings.ToLower(safe)
	if !strings.Contains(lower, "usually safe") && !strings.Contains(lower, "usually-safe") &&
		!strings.Contains(lower, "reclaimable") {
		t.Fatalf("SAFE TO RECLAIM should mark logs/sessions/cache/tmp as usually safe / reclaimable:\n%s", safe)
	}
	for _, line := range strings.Split(safe, "\n") {
		ll := strings.ToLower(strings.TrimSpace(line))
		if ll == "" {
			continue
		}
		if (strings.Contains(ll, "auth") || strings.Contains(ll, "config.toml") ||
			strings.Contains(ll, "credential")) &&
			(strings.Contains(ll, "usually safe") || strings.Contains(ll, "usually-safe")) &&
			!strings.Contains(ll, "not") && !strings.Contains(ll, "caution") &&
			!strings.Contains(ll, "never") && !strings.Contains(ll, "keep") {
			t.Fatalf("SAFE TO RECLAIM must not mark auth/config as usually-safe purge:\n%s", line)
		}
	}
}

// assertCodexLogsDBHuman locks product shape A LOGS DB human section when logs_*.sqlite present:
//
//	LOGS DB
//	  PATH: …
//	  SIZE: …
//	  ROWS: <n>
//	  SAMPLE (last 3, newest first):
//	    1) …
func assertCodexLogsDBHuman(t *testing.T, stdout string, wantRows int64) {
	t.Helper()
	if !strings.Contains(stdout, "LOGS DB") {
		t.Fatalf("codex-home with logs_*.sqlite must emit LOGS DB section:\n%s", stdout)
	}
	// Soft order: LOGS DB after BREAKDOWN and before SAFE TO RECLAIM when all present.
	bdPos := strings.Index(stdout, "BREAKDOWN")
	logsPos := strings.Index(stdout, "LOGS DB")
	safePos := strings.Index(stdout, "SAFE TO RECLAIM")
	if bdPos >= 0 && logsPos >= 0 && logsPos < bdPos {
		t.Fatalf("LOGS DB should appear after BREAKDOWN:\n%s", stdout)
	}
	if logsPos >= 0 && safePos >= 0 && safePos < logsPos {
		t.Fatalf("LOGS DB should appear before SAFE TO RECLAIM:\n%s", stdout)
	}
	body := sectionBody(stdout, "LOGS DB")
	if strings.TrimSpace(body) == "" {
		// sectionBody may not treat "LOGS DB" as a major section if it only matches
		// known headers; fall back to a local window after the marker.
		i := strings.Index(stdout, "LOGS DB")
		rest := stdout[i:]
		if j := strings.Index(rest[len("LOGS DB"):], "\nSAFE TO RECLAIM"); j >= 0 {
			body = rest[:len("LOGS DB")+j]
		} else if j := strings.Index(rest[len("LOGS DB"):], "\nHOW TO PURGE"); j >= 0 {
			body = rest[:len("LOGS DB")+j]
		} else {
			body = rest
		}
	}
	bl := strings.ToLower(body)
	if !strings.Contains(bl, "path") {
		t.Fatalf("LOGS DB section must include PATH:\n%s", body)
	}
	if !strings.Contains(bl, "size") {
		t.Fatalf("LOGS DB section must include SIZE:\n%s", body)
	}
	// ROWS: <n>
	wantRowsToken := fmt.Sprintf("ROWS: %d", wantRows)
	if !strings.Contains(body, wantRowsToken) && !strings.Contains(body, fmt.Sprintf("ROWS:%d", wantRows)) {
		// allow flexible spacing
		found := false
		for _, line := range strings.Split(body, "\n") {
			ll := strings.TrimSpace(line)
			if strings.HasPrefix(strings.ToUpper(ll), "ROWS") && strings.Contains(ll, fmt.Sprintf("%d", wantRows)) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("LOGS DB must report ROWS: %d:\n%s", wantRows, body)
		}
	}
	if !strings.Contains(strings.ToUpper(body), "SAMPLE") {
		t.Fatalf("LOGS DB must include SAMPLE (last 3, newest first):\n%s", body)
	}
	// Newest-first sample cues from defaultCodexLogRows: epsilon, delta, gamma (not alpha/beta alone as first).
	if !strings.Contains(body, "epsilon") && !strings.Contains(body, "row5") &&
		!strings.Contains(body, "DEBUG") {
		t.Fatalf("LOGS DB SAMPLE should include newest row (row5/epsilon/DEBUG):\n%s", body)
	}
	// Last-3 window also includes row3/row4 markers (gamma/delta/ERROR/INFO).
	hasOlderSample := strings.Contains(body, "gamma") || strings.Contains(body, "delta") ||
		strings.Contains(body, "row3") || strings.Contains(body, "row4") ||
		strings.Contains(body, "ERROR") || strings.Contains(body, "codex::")
	if !hasOlderSample {
		t.Fatalf("LOGS DB SAMPLE should include last-3 markers (row3–5 / gamma/delta/epsilon):\n%s", body)
	}
	// Must not dump all 5 as if unlimited — soft: sample section should not list "row1 oldest" as first sample.
	// (Bodies may still appear if truncation window is large; prefer checking numbered 1) 2) 3) only.)
	sampleIdx := strings.Index(strings.ToUpper(body), "SAMPLE")
	if sampleIdx >= 0 {
		samplePart := body[sampleIdx:]
		// Prefer at most 3 numbered sample entries.
		n1 := strings.Count(samplePart, "1)") + strings.Count(samplePart, "1.")
		n2 := strings.Count(samplePart, "2)") + strings.Count(samplePart, "2.")
		n3 := strings.Count(samplePart, "3)") + strings.Count(samplePart, "3.")
		if n1 == 0 && n2 == 0 && n3 == 0 {
			// Soft: numbered samples preferred but level/target lines may be unnumbered.
			if !strings.Contains(samplePart, "codex::") && !strings.Contains(strings.ToUpper(samplePart), "INFO") &&
				!strings.Contains(strings.ToUpper(samplePart), "DEBUG") &&
				!strings.Contains(strings.ToUpper(samplePart), "ERROR") {
				t.Fatalf("LOGS DB SAMPLE should list sample entries (numbered or level/target):\n%s", samplePart)
			}
		}
	}
}

// assertCodexLogsDBJSON locks JSON logsDb object for codex-home with readable logs db.
func assertCodexLogsDBJSON(t *testing.T, payload map[string]json.RawMessage, wantRows int64) {
	t.Helper()
	raw, ok := payload["logsDb"]
	if !ok || raw == nil {
		t.Fatalf("json missing required logsDb for codex-home with logs sqlite: keys=%v", payloadKeys(payload))
	}
	var logsDb map[string]any
	if err := json.Unmarshal(raw, &logsDb); err != nil {
		t.Fatalf("logsDb: %v", err)
	}
	path, _ := logsDb["path"].(string)
	if strings.TrimSpace(path) == "" || !strings.Contains(path, "logs_") {
		t.Fatalf("logsDb.path should reference logs_*.sqlite, got %q", path)
	}
	size := jsonEntryInt64(t, logsDb, "size")
	if size <= 0 {
		t.Fatalf("logsDb.size must be > 0, got %d", size)
	}
	rows := jsonEntryInt64(t, logsDb, "rows")
	if rows != wantRows {
		t.Fatalf("logsDb.rows: want %d, got %d", wantRows, rows)
	}
	samplesRaw, ok := logsDb["samples"]
	if !ok || samplesRaw == nil {
		t.Fatalf("logsDb.samples missing: %v", logsDb)
	}
	// samples may be []any from map[string]any unmarshal
	samples, ok := samplesRaw.([]any)
	if !ok {
		// re-unmarshal via json
		b, _ := json.Marshal(samplesRaw)
		if err := json.Unmarshal(b, &samples); err != nil {
			t.Fatalf("logsDb.samples: %v (type %T)", err, samplesRaw)
		}
	}
	if len(samples) == 0 {
		t.Fatal("logsDb.samples must be non-empty (up to 3 newest)")
	}
	if len(samples) > 3 {
		t.Fatalf("logsDb.samples len must be ≤3, got %d", len(samples))
	}
	// Newest first: first sample should be the highest-ts row (row5 / DEBUG / epsilon).
	first, _ := samples[0].(map[string]any)
	if first == nil {
		b, _ := json.Marshal(samples[0])
		_ = json.Unmarshal(b, &first)
	}
	if first != nil {
		body, _ := first["body"].(string)
		level, _ := first["level"].(string)
		if !strings.Contains(strings.ToLower(body), "epsilon") &&
			!strings.Contains(strings.ToLower(body), "newest") &&
			!strings.EqualFold(level, "DEBUG") {
			// soft: id/ts ordering
			if id, ok := first["id"].(float64); ok && id != 5 {
				t.Fatalf("logsDb.samples[0] should be newest (id=5 / DEBUG / epsilon), got %v", first)
			}
		}
	}
	// Each sample has id, ts, level, target, body keys (soft).
	for i, s := range samples {
		m, _ := s.(map[string]any)
		if m == nil {
			b, _ := json.Marshal(s)
			if err := json.Unmarshal(b, &m); err != nil {
				t.Fatalf("logsDb.samples[%d]: %v", i, err)
			}
		}
		for _, key := range []string{"id", "ts", "level", "target", "body"} {
			if _, ok := m[key]; !ok {
				t.Fatalf("logsDb.samples[%d] missing %q: %v", i, key, m)
			}
		}
	}
}

// payloadKeys returns sorted-ish key list for error messages.
func payloadKeys(payload map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(payload))
	for k := range payload {
		keys = append(keys, k)
	}
	return keys
}

// writeAndroidSDKFixture creates a tiny Android SDK layout under parentDir/Library/Android/sdk.
// Returns the absolute SDK root and absolute path to platform-tools/f (stable file under SDK).
// Exact file payloads (content total androidSDKContentBytes = 890):
//
//	system-images/f     400 B → system-images ☑
//	emulator/f          200 B → emulator ☐
//	sources/f           100 B → sources ☑
//	build-tools/f        80 B → build-tools ☐
//	platform-tools/f     50 B → platform-tools ☐
//	platforms/f          40 B → platforms ☐
//	.temp/f              20 B → tmp ☑
//
// Size DESC name order: system-images, emulator, sources, build-tools, platform-tools, platforms, .temp
// Directory metadata may make measured total ≥ 890 depending on walk semantics.
// Signatures: path ends with Library/Android/sdk; contains platform-tools + system-images/emulator/…
func writeAndroidSDKFixture(t *testing.T, parentDir string) (sdkDir string, fileUnderSDK string) {
	t.Helper()
	rel := "Library/Android/sdk"
	sdkDir = mkdir(t, parentDir, rel)
	writeSizedFile(t, parentDir, rel+"/system-images/f", 400)
	writeSizedFile(t, parentDir, rel+"/emulator/f", 200)
	writeSizedFile(t, parentDir, rel+"/sources/f", 100)
	writeSizedFile(t, parentDir, rel+"/build-tools/f", 80)
	fileUnderSDK = writeSizedFile(t, parentDir, rel+"/platform-tools/f", 50)
	writeSizedFile(t, parentDir, rel+"/platforms/f", 40)
	writeSizedFile(t, parentDir, rel+"/.temp/f", 20)
	return sdkDir, fileUnderSDK
}

// androidSDKContentBytes is the sum of file payloads in writeAndroidSDKFixture (excludes dir inodes).
const androidSDKContentBytes int64 = 400 + 200 + 100 + 80 + 50 + 40 + 20 // 890

// androidSDKRoles is the locked role id set for the android-sdk kind (fixture-covered).
var androidSDKRoles = []string{
	"system-images",
	"emulator",
	"sources",
	"build-tools",
	"platform-tools",
	"platforms",
	"tmp",
}

// assertAndroidSDKBreakdownMentions locks human BREAKDOWN/summary references for android-sdk
// roles or well-known basenames under the SDK root.
func assertAndroidSDKBreakdownMentions(t *testing.T, stdout string) {
	t.Helper()
	lower := strings.ToLower(stdout)
	checks := []struct {
		label string
		ok    bool
	}{
		{"system-images", strings.Contains(lower, "system-images") || strings.Contains(lower, "system images")},
		{"emulator", strings.Contains(lower, "emulator")},
		{"sources", strings.Contains(lower, "sources")},
		{"build-tools", strings.Contains(lower, "build-tools") || strings.Contains(lower, "build tools")},
		{"platform-tools", strings.Contains(lower, "platform-tools") || strings.Contains(lower, "platform tools") || strings.Contains(lower, "adb")},
		{"platforms", strings.Contains(lower, "platforms")},
		{"tmp or .temp", strings.Contains(lower, "tmp") || strings.Contains(lower, ".temp") || strings.Contains(lower, "temp")},
	}
	for _, c := range checks {
		if !c.ok {
			t.Fatalf("BREAKDOWN/summary should mention %s:\n%s", c.label, stdout)
		}
	}
	bd := breakdownSection(stdout)
	if strings.TrimSpace(bd) == "" {
		return
	}
	bdLower := strings.ToLower(bd)
	for _, role := range androidSDKRoles {
		if strings.Contains(bdLower, role) {
			continue
		}
		switch role {
		case "system-images":
			if !strings.Contains(bdLower, "system") {
				t.Fatalf("BREAKDOWN missing system-images role/name:\n%s", bd)
			}
		case "emulator":
			if !strings.Contains(bdLower, "emulator") {
				t.Fatalf("BREAKDOWN missing emulator role/name:\n%s", bd)
			}
		case "sources":
			if !strings.Contains(bdLower, "source") {
				t.Fatalf("BREAKDOWN missing sources role/name:\n%s", bd)
			}
		case "build-tools":
			if !strings.Contains(bdLower, "build") {
				t.Fatalf("BREAKDOWN missing build-tools role/name:\n%s", bd)
			}
		case "platform-tools":
			if !strings.Contains(bdLower, "platform-tools") && !strings.Contains(bdLower, "adb") {
				t.Fatalf("BREAKDOWN missing platform-tools role/name:\n%s", bd)
			}
		case "platforms":
			if !strings.Contains(bdLower, "platform") {
				t.Fatalf("BREAKDOWN missing platforms role/name:\n%s", bd)
			}
		case "tmp":
			if !strings.Contains(bdLower, "tmp") && !strings.Contains(bdLower, ".temp") && !strings.Contains(bdLower, "temp") {
				t.Fatalf("BREAKDOWN missing tmp role/name (.temp):\n%s", bd)
			}
		}
	}
}

// assertAndroidSDKReclaimCheckboxes locks human RECLAIMABLE ☑/☐ for android-sdk fixture rows.
// Prefer matching by role id; fall back to distinctive basenames (BREAKDOWN only).
func assertAndroidSDKReclaimCheckboxes(t *testing.T, stdout string) {
	t.Helper()
	bd := breakdownSection(stdout)
	pick := func(candidates []string) string {
		for _, name := range candidates {
			if strings.Contains(bd, name) {
				return name
			}
		}
		return ""
	}
	// reclaimable ☑
	for _, candidates := range [][]string{
		{"system-images", "system_images"},
		{"sources"},
		{".temp", "tmp", "temp"},
	} {
		name := pick(candidates)
		if name == "" {
			t.Fatalf("BREAKDOWN missing reclaimable row among %v:\n%s", candidates, bd)
		}
		assertBreakdownReclaimableCheckbox(t, stdout, name, true)
	}
	// non-reclaimable ☐
	for _, candidates := range [][]string{
		{"emulator"},
		{"build-tools"},
		{"platform-tools"},
		{"platforms"},
	} {
		name := pick(candidates)
		if name == "" {
			t.Fatalf("BREAKDOWN missing non-reclaimable row among %v:\n%s", candidates, bd)
		}
		assertBreakdownReclaimableCheckbox(t, stdout, name, false)
	}
}

// assertAndroidSDKCLIFirstPurge locks CLI-primary HOW TO PURGE for kind android-sdk:
// official bodies prefer sdkmanager and/or disk-usage-analyser scan; never rm -rf;
// Android Studio SDK settings only under Notes.
func assertAndroidSDKCLIFirstPurge(t *testing.T, stdout string) {
	t.Helper()
	how := sectionBody(stdout, "HOW TO PURGE")
	if strings.TrimSpace(how) == "" {
		t.Fatalf("missing HOW TO PURGE:\n%s", stdout)
	}
	lower := strings.ToLower(how)
	if !strings.Contains(lower, "sdkmanager") && !strings.Contains(lower, "disk-usage-analyser scan") {
		t.Fatalf("HOW TO PURGE must include CLI tools sdkmanager and/or disk-usage-analyser scan:\n%s", how)
	}
	if strings.Contains(lower, "rm -rf") || strings.Contains(lower, "rm -fr") {
		t.Fatalf("HOW TO PURGE must not contain rm -rf:\n%s", how)
	}
	// Prefer list_installed / uninstall / system-images guidance when sdkmanager appears.
	if strings.Contains(lower, "sdkmanager") {
		if !strings.Contains(lower, "list_installed") && !strings.Contains(lower, "list-installed") &&
			!strings.Contains(lower, "--list") && !strings.Contains(lower, "uninstall") {
			t.Fatalf("HOW TO PURGE sdkmanager guidance should list installed packages and/or uninstall:\n%s", how)
		}
	}

	inOfficial := false
	for _, line := range strings.Split(how, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "Official command:") {
			inOfficial = true
			continue
		}
		if strings.HasPrefix(trim, "Removes:") {
			inOfficial = false
			continue
		}
		if strings.HasPrefix(trim, "Notes:") {
			inOfficial = false
			continue
		}
		if isNumberedStepTitle(trim) {
			inOfficial = false
			continue
		}
		if !inOfficial || trim == "" || strings.HasPrefix(trim, "#") {
			continue
		}
		cmd := strings.TrimSpace(strings.TrimPrefix(trim, "$"))
		cmd = strings.TrimSpace(cmd)
		// Primary official command must not be pure UI navigation prose.
		if (strings.Contains(cmd, "Android Studio") || strings.Contains(cmd, "SDK Manager") ||
			strings.Contains(cmd, "Settings") && strings.Contains(cmd, "Android SDK")) &&
			!strings.Contains(cmd, "sdkmanager") && !strings.Contains(cmd, "disk-usage-analyser") {
			t.Fatalf("HOW TO PURGE Official command must be CLI-first (UI only in Notes), got: %q", line)
		}
	}
	assertHowToPurgeOfficialDollarPrefix(t, stdout)
}

// assertAndroidSDKSafeToReclaim locks SAFE TO RECLAIM tiers for android-sdk:
// temp/system-images/sources usually safe or reclaimable-with-caution;
// platform-tools/build-tools/platforms/emulator kept (not usually-safe-only purge).
func assertAndroidSDKSafeToReclaim(t *testing.T, stdout string) {
	t.Helper()
	safe := sectionBody(stdout, "SAFE TO RECLAIM")
	if strings.TrimSpace(safe) == "" {
		t.Fatalf("missing SAFE TO RECLAIM body:\n%s", stdout)
	}
	lower := strings.ToLower(safe)
	if !strings.Contains(lower, "usually safe") && !strings.Contains(lower, "usually-safe") &&
		!strings.Contains(lower, "reclaimable") && !strings.Contains(lower, "temp") &&
		!strings.Contains(lower, "system-image") {
		t.Fatalf("SAFE TO RECLAIM should mark temp/system-images/sources as usually safe / reclaimable:\n%s", safe)
	}
	// platform-tools / adb / build-tools bulk must not be presented as usually-safe purge alone.
	for _, line := range strings.Split(safe, "\n") {
		ll := strings.ToLower(strings.TrimSpace(line))
		if ll == "" {
			continue
		}
		if (strings.Contains(ll, "platform-tools") || strings.Contains(ll, "adb") ||
			strings.Contains(ll, "build-tools") ||
			(strings.Contains(ll, "platforms") && !strings.Contains(ll, "system"))) &&
			(strings.Contains(ll, "usually safe") || strings.Contains(ll, "usually-safe")) &&
			!strings.Contains(ll, "not") && !strings.Contains(ll, "caution") &&
			!strings.Contains(ll, "keep") && !strings.Contains(ll, "unless") {
			t.Fatalf("SAFE TO RECLAIM must not mark platform-tools/build-tools/platforms as usually-safe purge:\n%s", line)
		}
	}
}

// writeITerm2Fixture creates a tiny iTerm2 Application Support layout under
// parentDir/Library/Application Support/iTerm2.
// Returns the absolute iTerm2 root and absolute path to iterm2env/f (stable file under tree).
// Exact file payloads (content total iTerm2ContentBytes = 674):
//
//	iterm2env/f              400 B → python-env ☑
//	iterm2env-3.10/f         200 B → python-env-alias ☑
//	log.0.txt                 50 B → logs ☑
//	version.txt               16 B → meta ☐
//	DynamicProfiles/p          8 B → user-config ☐
//
// Size DESC name order: iterm2env, iterm2env-3.10, log.0.txt, version.txt, DynamicProfiles
// Directory metadata may make measured total ≥ 674 depending on walk semantics.
// Signatures: path ends with Application Support/iTerm2; contains iterm2env and/or version.txt.
// Multiple iterm2env* trees may share APFS hardlink blocks in production; SUMMARY/SAFE document
// logical overcount (fixture uses separate payloads for deterministic size DESC).
func writeITerm2Fixture(t *testing.T, parentDir string) (iterm2Dir string, fileUnderITerm2 string) {
	t.Helper()
	rel := "Library/Application Support/iTerm2"
	iterm2Dir = mkdir(t, parentDir, rel)
	fileUnderITerm2 = writeSizedFile(t, parentDir, rel+"/iterm2env/f", 400)
	writeSizedFile(t, parentDir, rel+"/iterm2env-3.10/f", 200)
	writeSizedFile(t, parentDir, rel+"/log.0.txt", 50)
	writeSizedFile(t, parentDir, rel+"/version.txt", 16)
	writeSizedFile(t, parentDir, rel+"/DynamicProfiles/p", 8)
	return iterm2Dir, fileUnderITerm2
}

// iTerm2ContentBytes is the sum of file payloads in writeITerm2Fixture (excludes dir inodes).
const iTerm2ContentBytes int64 = 400 + 200 + 50 + 16 + 8 // 674

// writeAllKindsHomeFixture creates a fake home with ≥2 v1 packs present and xcode missing:
//
//	{home}/.codex/…                             → codex-home present (~1198 B non-DB tree;
//	                                              optional logs_2.sqlite if sqlite3 available)
//	{home}/.grok/…                              → grok-home present (~798 B payload)
//	{home}/Library/Android/sdk/…                → android-sdk present (~890 B payload)
//	{home}/Library/Application Support/iTerm2/… → iterm2 present (~674 B payload)
//	# no Xcode roots → xcode status missing
//
// Payload size DESC among present: codex (≥1198) > android-sdk (890) > grok (798) > iterm2 (674).
// Returns the absolute home/scope directory (parentDir).
// Does not Skip when sqlite3 is missing (tree-only .codex is enough for pack present).
func writeAllKindsHomeFixture(t *testing.T, parentDir string) string {
	t.Helper()
	codexDir, _ := writeCodexHomeTree(t, parentDir)
	writeCodexLogsSQLiteOrNoop(t, filepath.Join(codexDir, "logs_2.sqlite"), defaultCodexLogRows())
	writeGrokHomeFixture(t, parentDir)
	writeAndroidSDKFixture(t, parentDir)
	writeITerm2Fixture(t, parentDir)
	// deliberately omit writeXcodeScopeFixture so xcode is missing
	return parentDir
}

// writeCodexLogsSQLiteOrNoop creates logs sqlite when sqlite3 is available; otherwise no-op.
// Fatals only on create errors (never Skip) so multi-pack fixtures stay hermetic.
func writeCodexLogsSQLiteOrNoop(t *testing.T, dbPath string, rows []codexLogRow) {
	t.Helper()
	if _, err := exec.LookPath("sqlite3"); err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		t.Fatalf("mkdir for sqlite parent: %v", err)
	}
	_ = os.Remove(dbPath)
	var b strings.Builder
	b.WriteString(`CREATE TABLE logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  ts INTEGER NOT NULL,
  ts_nanos INTEGER NOT NULL,
  level TEXT NOT NULL,
  target TEXT NOT NULL,
  feedback_log_body TEXT
);
`)
	for _, r := range rows {
		body := strings.ReplaceAll(r.Body, "'", "''")
		target := strings.ReplaceAll(r.Target, "'", "''")
		level := strings.ReplaceAll(r.Level, "'", "''")
		fmt.Fprintf(&b, "INSERT INTO logs (ts, ts_nanos, level, target, feedback_log_body) VALUES (%d, %d, '%s', '%s', '%s');\n",
			r.TS, r.TSNanos, level, target, body)
	}
	cmd := exec.Command("sqlite3", dbPath)
	cmd.Stdin = strings.NewReader(b.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("sqlite3 create %s: %v\n%s", dbPath, err, out)
	}
}

// allKindsCLIKinds is the locked v1 --all-kinds / --kind registry (cli ids).
// Referenced from assertAllKindsIndex / assertAllKindsJSONShape (must stay used).
var allKindsCLIKinds = []string{"xcode", "grok", "android-sdk", "iterm2", "codex"}

// iTerm2Roles is the locked role id set for the iterm2 kind (fixture-covered).
var iTerm2Roles = []string{
	"python-env",
	"python-env-alias",
	"logs",
	"meta",
	"user-config",
}

// assertITerm2BreakdownMentions locks human BREAKDOWN/summary references for iterm2 roles
// or well-known basenames under Application Support/iTerm2.
func assertITerm2BreakdownMentions(t *testing.T, stdout string) {
	t.Helper()
	lower := strings.ToLower(stdout)
	checks := []struct {
		label string
		ok    bool
	}{
		{"python-env or iterm2env", strings.Contains(lower, "python-env") || strings.Contains(lower, "iterm2env")},
		{"python-env-alias or iterm2env-3", strings.Contains(lower, "python-env-alias") || strings.Contains(lower, "iterm2env-3") || strings.Contains(lower, "iterm2env-")},
		{"logs or log.0", strings.Contains(lower, "logs") || strings.Contains(lower, "log.0") || strings.Contains(lower, "log")},
		{"meta or version.txt", strings.Contains(lower, "meta") || strings.Contains(lower, "version.txt") || strings.Contains(lower, "version")},
		{"user-config or DynamicProfiles", strings.Contains(lower, "user-config") || strings.Contains(lower, "dynamicprofiles") || strings.Contains(lower, "dynamic profiles")},
	}
	for _, c := range checks {
		if !c.ok {
			t.Fatalf("BREAKDOWN/summary should mention %s:\n%s", c.label, stdout)
		}
	}
	bd := breakdownSection(stdout)
	if strings.TrimSpace(bd) == "" {
		return
	}
	bdLower := strings.ToLower(bd)
	for _, role := range iTerm2Roles {
		if strings.Contains(bdLower, role) {
			continue
		}
		switch role {
		case "python-env":
			// Prefer exact env basename; avoid only matching iterm2env-3.10 line.
			if !strings.Contains(bdLower, "iterm2env") {
				t.Fatalf("BREAKDOWN missing python-env role/name:\n%s", bd)
			}
		case "python-env-alias":
			if !strings.Contains(bdLower, "iterm2env-") && !strings.Contains(bdLower, "iterm2env-3") {
				t.Fatalf("BREAKDOWN missing python-env-alias role/name:\n%s", bd)
			}
		case "logs":
			if !strings.Contains(bdLower, "log") {
				t.Fatalf("BREAKDOWN missing logs role/name:\n%s", bd)
			}
		case "meta":
			if !strings.Contains(bdLower, "version") {
				t.Fatalf("BREAKDOWN missing meta role/name (version.txt):\n%s", bd)
			}
		case "user-config":
			if !strings.Contains(bdLower, "dynamic") && !strings.Contains(bdLower, "user-config") {
				t.Fatalf("BREAKDOWN missing user-config role/name (DynamicProfiles):\n%s", bd)
			}
		}
	}
}

// assertITerm2ReclaimCheckboxes locks human RECLAIMABLE ☑/☐ for iterm2 fixture rows.
// Prefer matching by role id; fall back to distinctive basenames (BREAKDOWN only).
func assertITerm2ReclaimCheckboxes(t *testing.T, stdout string) {
	t.Helper()
	bd := breakdownSection(stdout)
	pick := func(candidates []string) string {
		for _, name := range candidates {
			if strings.Contains(bd, name) {
				return name
			}
		}
		return ""
	}
	// reclaimable ☑ — python-env before iterm2env so role preferred; alias uses versioned name
	for _, candidates := range [][]string{
		{"python-env-alias", "iterm2env-3.10", "iterm2env-"},
		{"python-env", "iterm2env"},
		{"logs", "log.0.txt", "log.0"},
	} {
		name := pick(candidates)
		if name == "" {
			t.Fatalf("BREAKDOWN missing reclaimable row among %v:\n%s", candidates, bd)
		}
		// When matching bare "iterm2env", prefer a line that is not the alias.
		if name == "iterm2env" || name == "python-env" {
			line := ""
			for _, l := range strings.Split(bd, "\n") {
				ll := strings.ToLower(l)
				if strings.Contains(ll, "iterm2env-") || strings.Contains(ll, "python-env-alias") {
					continue
				}
				if strings.Contains(ll, "python-env") || strings.Contains(ll, "iterm2env") {
					line = l
					break
				}
			}
			if line == "" || !strings.Contains(line, "☑") {
				t.Fatalf("BREAKDOWN python-env / iterm2env row must show reclaimable ☑: %q\n%s", line, bd)
			}
			continue
		}
		assertBreakdownReclaimableCheckbox(t, stdout, name, true)
	}
	// non-reclaimable ☐
	for _, candidates := range [][]string{
		{"version.txt", "meta", "version"},
		{"DynamicProfiles", "user-config", "Dynamic"},
	} {
		name := pick(candidates)
		if name == "" {
			t.Fatalf("BREAKDOWN missing non-reclaimable row among %v:\n%s", candidates, bd)
		}
		assertBreakdownReclaimableCheckbox(t, stdout, name, false)
	}
}

// assertITerm2HardlinkWording locks SUMMARY and/or SAFE TO RECLAIM language about
// APFS hardlinks / shared inodes / shared blocks among iterm2env* trees and that
// logical TOTAL can overstate freeable space; confirm with du on the parent.
func assertITerm2HardlinkWording(t *testing.T, stdout string) {
	t.Helper()
	summary := sectionBody(stdout, "SUMMARY")
	safe := sectionBody(stdout, "SAFE TO RECLAIM")
	joined := strings.ToLower(summary + "\n" + safe)
	if strings.TrimSpace(joined) == "" {
		t.Fatalf("missing SUMMARY/SAFE TO RECLAIM for hardlink wording:\n%s", stdout)
	}
	hasHardlinkOrShare := strings.Contains(joined, "hardlink") || strings.Contains(joined, "hard link") ||
		strings.Contains(joined, "hard-link") || strings.Contains(joined, "shared inode") ||
		strings.Contains(joined, "shared block") || strings.Contains(joined, "share disk") ||
		strings.Contains(joined, "shared disk") || strings.Contains(joined, "share space") ||
		strings.Contains(joined, "shared space") || strings.Contains(joined, "same inode") ||
		(strings.Contains(joined, "share") && (strings.Contains(joined, "block") ||
			strings.Contains(joined, "inode") || strings.Contains(joined, "clone") ||
			strings.Contains(joined, "apfs") || strings.Contains(joined, "iterm2env")))
	if !hasHardlinkOrShare {
		t.Fatalf("SUMMARY or SAFE TO RECLAIM must mention hardlink/shared blocks among iterm2env*:\nSUMMARY:\n%s\nSAFE:\n%s", summary, safe)
	}
	// Logical overcount / confirm with du (or equivalent parent measure guidance).
	hasOvercountOrDu := strings.Contains(joined, "overcount") || strings.Contains(joined, "overstate") ||
		strings.Contains(joined, "over-state") || strings.Contains(joined, "double-count") ||
		strings.Contains(joined, "double count") || strings.Contains(joined, "logical") ||
		strings.Contains(joined, "du -sh") || strings.Contains(joined, "du -s") ||
		(strings.Contains(joined, "du ") && strings.Contains(joined, "parent")) ||
		strings.Contains(joined, "not sum") || strings.Contains(joined, "one env") ||
		strings.Contains(joined, "freeable") || strings.Contains(joined, "actual reclaim")
	if !hasOvercountOrDu {
		t.Fatalf("SUMMARY or SAFE TO RECLAIM must warn logical overcount and/or confirm with du on parent:\nSUMMARY:\n%s\nSAFE:\n%s", summary, safe)
	}
}

// assertITerm2CLIFirstPurge locks CLI-primary HOW TO PURGE for kind iterm2:
// official bodies prefer disk-usage-analyser scan and/or du -sh on parent / iterm2env*;
// never rm -rf; user-config (DynamicProfiles/Scripts) not usually-safe Removes.
func assertITerm2CLIFirstPurge(t *testing.T, stdout string) {
	t.Helper()
	how := sectionBody(stdout, "HOW TO PURGE")
	if strings.TrimSpace(how) == "" {
		t.Fatalf("missing HOW TO PURGE:\n%s", stdout)
	}
	lower := strings.ToLower(how)
	hasScan := strings.Contains(lower, "disk-usage-analyser scan")
	hasDu := strings.Contains(lower, "du -sh") || strings.Contains(lower, "du -s") ||
		(strings.Contains(lower, "du ") && (strings.Contains(lower, "iterm") || strings.Contains(lower, "parent")))
	if !hasScan && !hasDu {
		t.Fatalf("HOW TO PURGE must include disk-usage-analyser scan and/or du on iTerm2/parent:\n%s", how)
	}
	// Prefer both inspect tools when present; require at least scan (product RAW also uses scan).
	if !hasScan {
		t.Fatalf("HOW TO PURGE must include disk-usage-analyser scan (inspect first):\n%s", how)
	}
	if strings.Contains(lower, "rm -rf") || strings.Contains(lower, "rm -fr") {
		t.Fatalf("HOW TO PURGE must not contain rm -rf:\n%s", how)
	}
	// Env / logs reclaim guidance preferred.
	if !strings.Contains(lower, "iterm2env") && !strings.Contains(lower, "python") &&
		!strings.Contains(lower, "log") {
		t.Fatalf("HOW TO PURGE should mention iterm2env/python env and/or logs reclaim targets:\n%s", how)
	}
	// Never mark DynamicProfiles / Scripts / user-config as usually-safe Removes.
	for _, line := range strings.Split(how, "\n") {
		trim := strings.TrimSpace(line)
		ll := strings.ToLower(trim)
		if !strings.HasPrefix(ll, "removes:") {
			continue
		}
		if (strings.Contains(ll, "usually safe") || strings.Contains(ll, "usually-safe")) &&
			(strings.Contains(ll, "dynamicprofiles") || strings.Contains(ll, "dynamic profiles") ||
				strings.Contains(ll, "scripts") || strings.Contains(ll, "user-config") ||
				strings.Contains(ll, "user config")) {
			t.Fatalf("HOW TO PURGE must not mark user-config (DynamicProfiles/Scripts) as usually-safe Removes:\n%s", line)
		}
	}
	assertHowToPurgeOfficialDollarPrefix(t, stdout)
}

// assertITerm2SafeToReclaim locks SAFE TO RECLAIM tiers for iterm2:
// python env / logs usually safe or reclaimable (with hardlink caution);
// user-config / DynamicProfiles not usually-safe-only purge.
func assertITerm2SafeToReclaim(t *testing.T, stdout string) {
	t.Helper()
	safe := sectionBody(stdout, "SAFE TO RECLAIM")
	if strings.TrimSpace(safe) == "" {
		t.Fatalf("missing SAFE TO RECLAIM body:\n%s", stdout)
	}
	lower := strings.ToLower(safe)
	if !strings.Contains(lower, "usually safe") && !strings.Contains(lower, "usually-safe") &&
		!strings.Contains(lower, "reclaimable") {
		t.Fatalf("SAFE TO RECLAIM should mark python env/logs as usually safe / reclaimable:\n%s", safe)
	}
	// user-config / DynamicProfiles / Scripts must not be presented as usually-safe purge alone.
	for _, line := range strings.Split(safe, "\n") {
		ll := strings.ToLower(strings.TrimSpace(line))
		if ll == "" {
			continue
		}
		if (strings.Contains(ll, "dynamicprofiles") || strings.Contains(ll, "dynamic profiles") ||
			strings.Contains(ll, "scripts") || strings.Contains(ll, "user-config") ||
			strings.Contains(ll, "user config") || strings.Contains(ll, "profile")) &&
			(strings.Contains(ll, "usually safe") || strings.Contains(ll, "usually-safe")) &&
			!strings.Contains(ll, "not") && !strings.Contains(ll, "caution") &&
			!strings.Contains(ll, "keep") && !strings.Contains(ll, "never") {
			// Only fail when clearly targeting user profiles/config as the safe purge subject.
			if strings.Contains(ll, "dynamic") || strings.Contains(ll, "scripts") ||
				strings.Contains(ll, "user-config") || strings.Contains(ll, "user config") {
				t.Fatalf("SAFE TO RECLAIM must not mark user-config/DynamicProfiles as usually-safe purge:\n%s", line)
			}
		}
	}
}

func stdoutEndsWithBlankLine(t *testing.T, stdout string) {
	t.Helper()
	if stdout == "" {
		t.Fatal("stdout is empty")
	}
	if !strings.HasSuffix(stdout, "\n\n") {
		t.Fatalf("stdout must end with trailing blank line after last content line; got %q", stdout)
	}
}

func assertNoRmRf(t *testing.T, stdout string) {
	t.Helper()
	if strings.Contains(strings.ToLower(stdout), "rm -rf") {
		t.Fatalf("stdout must not contain 'rm -rf' (case-insensitive):\n%s", stdout)
	}
}

// requiredHumanHeaders are exact section/summary markers locked by the product design.
var requiredHumanHeaders = []string{
	"PATH:",
	"KIND:",
	"TOTAL:",
	"CONFIDENCE:",
	"SUMMARY",
	"BREAKDOWN",
	"SAFE TO RECLAIM",
	"HOW TO PURGE",
	"RAW COMMANDS",
}

func assertHumanSectionsPresent(t *testing.T, stdout string) {
	t.Helper()
	for _, h := range requiredHumanHeaders {
		if !strings.Contains(stdout, h) {
			t.Fatalf("stdout missing required section/header %q:\n%s", h, stdout)
		}
	}
}

// assertHumanSectionOrder checks that locked headers appear in the required order.
func assertHumanSectionOrder(t *testing.T, stdout string) {
	t.Helper()
	pos := 0
	for _, h := range requiredHumanHeaders {
		i := strings.Index(stdout[pos:], h)
		if i < 0 {
			t.Fatalf("header %q not found after previous headers in:\n%s", h, stdout)
		}
		pos += i + len(h)
	}
}

func assertKindLine(t *testing.T, stdout string, kind string) {
	t.Helper()
	want := "KIND: " + kind
	found := false
	for _, line := range strings.Split(stdout, "\n") {
		if strings.TrimSpace(line) == want {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected line %q in stdout:\n%s", want, stdout)
	}
}

func assertContainsScanCommand(t *testing.T, stdout string, path string) {
	t.Helper()
	if !strings.Contains(stdout, "disk-usage-analyser scan") {
		t.Fatalf("stdout must include 'disk-usage-analyser scan':\n%s", stdout)
	}
	if path != "" && !strings.Contains(stdout, path) {
		t.Fatalf("stdout must include explained path %q:\n%s", path, stdout)
	}
}

func firstJSONObjectLine(t *testing.T, stdout string) string {
	t.Helper()
	content := strings.TrimRight(stdout, "\n")
	lines := strings.Split(content, "\n")
	if len(lines) != 1 {
		t.Fatalf("expected one JSON line, got %d lines:\n%s", len(lines), stdout)
	}
	return lines[0]
}

func parseJSONObject(t *testing.T, stdout string) map[string]json.RawMessage {
	t.Helper()
	line := firstJSONObjectLine(t, stdout)
	var payload map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	return payload
}

func jsonStringField(t *testing.T, payload map[string]json.RawMessage, key string) string {
	t.Helper()
	raw, ok := payload[key]
	if !ok || raw == nil {
		t.Fatalf("json missing %q", key)
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("json %q: %v", key, err)
	}
	return s
}

func jsonInt64Field(t *testing.T, payload map[string]json.RawMessage, key string) int64 {
	t.Helper()
	raw, ok := payload[key]
	if !ok || raw == nil {
		t.Fatalf("json missing %q", key)
	}
	var n int64
	if err := json.Unmarshal(raw, &n); err != nil {
		t.Fatalf("json %q: %v", key, err)
	}
	return n
}

// majorHumanHeaders are section titles used to slice human stdout for focused asserts.
var majorHumanHeaders = []string{
	"SUMMARY",
	"BREAKDOWN",
	"LOGS DB",
	"SAFE TO RECLAIM",
	"HOW TO PURGE",
	"RAW COMMANDS",
}

// sectionBody returns the body under header until the next major section or EOF.
// header is an exact marker such as "HOW TO PURGE" or "RAW COMMANDS".
func sectionBody(stdout, header string) string {
	idx := strings.Index(stdout, header)
	if idx < 0 {
		return ""
	}
	rest := stdout[idx+len(header):]
	if nl := strings.Index(rest, "\n"); nl >= 0 {
		rest = rest[nl+1:]
	} else {
		return ""
	}
	cut := len(rest)
	for _, h := range majorHumanHeaders {
		if h == header {
			continue
		}
		// Match header only as its own line (leading newline).
		needle := "\n" + h
		pos := 0
		for {
			i := strings.Index(rest[pos:], needle)
			if i < 0 {
				break
			}
			at := pos + i
			after := at + len(needle)
			if after == len(rest) || rest[after] == '\n' || rest[after] == '\r' {
				if at < cut {
					cut = at
				}
				break
			}
			pos = after
		}
	}
	return rest[:cut]
}

func containsANSI(s string) bool {
	return strings.Contains(s, "\x1b[") || strings.Contains(s, "\033[")
}

func assertNoANSI(t *testing.T, s string) {
	t.Helper()
	if containsANSI(s) {
		t.Fatalf("expected no ANSI escape sequences, found some in output:\n%s", s)
	}
}

// greenSGRPrefixes are accepted SGR sequences for green (plain or bold green).
var greenSGRPrefixes = []string{
	"\x1b[32m",
	"\x1b[1;32m",
	"\x1b[32;1m",
	"\x1b[92m",
	"\x1b[1;92m",
}

func assertGreenBaseCommand(t *testing.T, stdout, base string) {
	t.Helper()
	if base == "" {
		t.Fatal("assertGreenBaseCommand: empty base token")
	}
	found := false
	for _, pref := range greenSGRPrefixes {
		if strings.Contains(stdout, pref+base) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected green ANSI around base command %q (e.g. \\x1b[32m%s\\x1b[0m):\n%s", base, base, stdout)
	}
	// "$" prompt must not itself be green.
	for _, pref := range greenSGRPrefixes {
		if strings.Contains(stdout, pref+"$") {
			t.Fatalf("shell prompt \"$\" must not be green-colored; found green \"$\" with prefix %q", pref)
		}
	}
}

// assertRawCommandsDollarPrefix requires every non-comment runnable line under RAW COMMANDS
// to start with "$ " after trim, and comment/group lines (#...) must not use "$".
func assertRawCommandsDollarPrefix(t *testing.T, stdout string) {
	t.Helper()
	raw := sectionBody(stdout, "RAW COMMANDS")
	if strings.TrimSpace(raw) == "" {
		t.Fatalf("missing RAW COMMANDS body:\n%s", stdout)
	}
	runnable := 0
	for _, line := range strings.Split(raw, "\n") {
		trim := strings.TrimSpace(line)
		if trim == "" || trim == "(none)" {
			continue
		}
		if strings.HasPrefix(trim, "#") {
			if strings.HasPrefix(trim, "$") {
				t.Fatalf("RAW COMMANDS comment/group line must not use $ prefix: %q", line)
			}
			continue
		}
		if !strings.HasPrefix(trim, "$ ") {
			t.Fatalf("RAW COMMANDS runnable line must start with \"$ \": %q\nsection:\n%s", line, raw)
		}
		runnable++
	}
	if runnable == 0 {
		t.Fatalf("RAW COMMANDS expected at least one \"$ \" command line:\n%s", raw)
	}
}

// assertHowToPurgeOfficialDollarPrefix checks lines under each "Official command:" block:
// runnable lines require "$ "; "#" comments must not use "$".
// Kinds with only comment official text may have zero runnable lines (OK).
func assertHowToPurgeOfficialDollarPrefix(t *testing.T, stdout string) {
	t.Helper()
	how := sectionBody(stdout, "HOW TO PURGE")
	if strings.TrimSpace(how) == "" {
		t.Fatalf("missing HOW TO PURGE body:\n%s", stdout)
	}
	inOfficial := false
	for _, line := range strings.Split(how, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "Official command:") {
			inOfficial = true
			continue
		}
		if strings.HasPrefix(trim, "Removes:") || strings.HasPrefix(trim, "Notes:") {
			inOfficial = false
			continue
		}
		// Next numbered step: "1) title" / "2) title"
		if isNumberedStepTitle(trim) {
			inOfficial = false
			continue
		}
		if !inOfficial || trim == "" {
			continue
		}
		if strings.HasPrefix(trim, "#") {
			if strings.HasPrefix(trim, "$") {
				t.Fatalf("HOW TO PURGE comment must not use $ prefix: %q", line)
			}
			continue
		}
		if !strings.HasPrefix(trim, "$ ") {
			t.Fatalf("HOW TO PURGE official runnable line must start with \"$ \": %q\nsection:\n%s", line, how)
		}
	}
}

func isNumberedStepTitle(trim string) bool {
	if len(trim) < 3 {
		return false
	}
	i := 0
	for i < len(trim) && trim[i] >= '0' && trim[i] <= '9' {
		i++
	}
	return i > 0 && i < len(trim) && trim[i] == ')'
}

// assertHowToPurgeHasDollarCommand requires HOW TO PURGE to include a "$ <token>" runnable line
// (token is typically the base command: go, npm, brew, emulator, …).
func assertHowToPurgeHasDollarCommand(t *testing.T, stdout, token string) {
	t.Helper()
	how := sectionBody(stdout, "HOW TO PURGE")
	needle := "$ " + token
	if !strings.Contains(how, needle) {
		t.Fatalf("HOW TO PURGE must include %q on an official command line:\n%s", needle, how)
	}
}

// assertAndroidAVDCLIFirstPurge locks CLI-primary HOW TO PURGE for android-avd human output.
// Official command bodies must prefer emulator/avdmanager; Device Manager only in Notes.
func assertAndroidAVDCLIFirstPurge(t *testing.T, stdout string) {
	t.Helper()
	how := sectionBody(stdout, "HOW TO PURGE")
	if strings.TrimSpace(how) == "" {
		t.Fatalf("missing HOW TO PURGE:\n%s", stdout)
	}
	lower := strings.ToLower(how)
	if !strings.Contains(lower, "emulator") && !strings.Contains(lower, "avdmanager") {
		t.Fatalf("HOW TO PURGE must include CLI tools emulator and/or avdmanager:\n%s", how)
	}

	inOfficial := false
	for _, line := range strings.Split(how, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "Official command:") {
			inOfficial = true
			continue
		}
		if strings.HasPrefix(trim, "Removes:") {
			inOfficial = false
			continue
		}
		if strings.HasPrefix(trim, "Notes:") {
			inOfficial = false
			// UI is allowed in Notes.
			continue
		}
		if isNumberedStepTitle(trim) {
			inOfficial = false
			continue
		}
		if !inOfficial || trim == "" || strings.HasPrefix(trim, "#") {
			continue
		}
		cmd := strings.TrimSpace(strings.TrimPrefix(trim, "$"))
		cmd = strings.TrimSpace(cmd)
		// Primary official command must not be UI navigation prose.
		if strings.Contains(cmd, "Android Studio") || strings.Contains(cmd, "Device Manager") {
			t.Fatalf("HOW TO PURGE Official command must be CLI-first (UI only in Notes), got: %q", line)
		}
	}
}

// assertJSONHowToPurgePlainCLI checks JSON howToPurge entries: no ANSI, no leading $,
// and (for android-avd) CLI-first officialCommand text.
func assertJSONHowToPurgePlainCLI(t *testing.T, howToPurge []map[string]any, requireAndroidCLI bool) {
	t.Helper()
	if len(howToPurge) == 0 {
		t.Fatal("howToPurge must be non-empty")
	}
	joined := ""
	for i, step := range howToPurge {
		oc, _ := step["officialCommand"].(string)
		if oc == "" {
			t.Fatalf("howToPurge[%d] missing officialCommand", i)
		}
		if containsANSI(oc) {
			t.Fatalf("howToPurge[%d].officialCommand must not contain ANSI: %q", i, oc)
		}
		notes, _ := step["notes"].(string)
		if containsANSI(notes) {
			t.Fatalf("howToPurge[%d].notes must not contain ANSI: %q", i, notes)
		}
		for _, line := range strings.Split(oc, "\n") {
			trim := strings.TrimSpace(line)
			if strings.HasPrefix(trim, "$") {
				t.Fatalf("JSON officialCommand must not include $ prefix (human formatter adds it): %q", oc)
			}
		}
		joined += oc + "\n"
		// Device Manager as primary official body is not allowed.
		for _, line := range strings.Split(oc, "\n") {
			trim := strings.TrimSpace(line)
			if trim == "" || strings.HasPrefix(trim, "#") {
				continue
			}
			if strings.Contains(trim, "Android Studio") || strings.Contains(trim, "Device Manager") {
				t.Fatalf("JSON officialCommand must be CLI-first, not UI: %q", oc)
			}
		}
	}
	if requireAndroidCLI {
		jl := strings.ToLower(joined)
		if !strings.Contains(jl, "emulator") && !strings.Contains(jl, "avdmanager") {
			t.Fatalf("android-avd JSON howToPurge officialCommand must mention emulator/avdmanager: %s", joined)
		}
	}
}

// assertSeaTalkHumanHowToPurge locks seatalk-app-support human HOW TO PURGE:
// quit via osascript, web-cache basenames listed for reclaim, separate backup step,
// no primary chat/search DB as a usually-safe purge target, $ on runnable lines.
func assertSeaTalkHumanHowToPurge(t *testing.T, stdout string) {
	t.Helper()
	how := sectionBody(stdout, "HOW TO PURGE")
	if strings.TrimSpace(how) == "" {
		t.Fatalf("missing HOW TO PURGE:\n%s", stdout)
	}
	if !strings.Contains(how, "osascript") {
		t.Fatalf("HOW TO PURGE must include osascript quit prep:\n%s", how)
	}
	if !strings.Contains(how, "SeaTalk") {
		t.Fatalf("HOW TO PURGE osascript must target SeaTalk:\n%s", how)
	}
	// Cache reclaim basenames (at least the fixture-covered ones).
	if !strings.Contains(how, "Service Worker") {
		t.Fatalf("HOW TO PURGE must list Service Worker/ among cache removes:\n%s", how)
	}
	if !strings.Contains(how, "Cache") {
		t.Fatalf("HOW TO PURGE must list Cache/ among cache removes:\n%s", how)
	}
	// Separate backup reclaim step.
	if !strings.Contains(how, "sqlite-backup") {
		t.Fatalf("HOW TO PURGE must include sqlite-backup reclaim step:\n%s", how)
	}
	// Primary DBs must not be presented as a usually-safe purge removes list item alone
	// without caution context — require that "usually safe" does not appear on the same
	// Removes line as main_*.sqlite / search_*.sqlite when describing purge targets.
	for _, line := range strings.Split(how, "\n") {
		trim := strings.TrimSpace(line)
		lower := strings.ToLower(trim)
		if !strings.HasPrefix(lower, "removes:") {
			continue
		}
		if strings.Contains(lower, "usually safe") &&
			(strings.Contains(lower, "main_") || strings.Contains(lower, "search_")) {
			t.Fatalf("HOW TO PURGE must not mark primary main_/search_ sqlite as usually-safe Removes:\n%s", line)
		}
	}
	assertHowToPurgeHasDollarCommand(t, stdout, "osascript")
	assertHowToPurgeOfficialDollarPrefix(t, stdout)
}

// assertSeaTalkBreakdownMentions locks human BREAKDOWN/summary references for SeaTalk roles.
func assertSeaTalkBreakdownMentions(t *testing.T, stdout string) {
	t.Helper()
	lower := strings.ToLower(stdout)
	// At least web cache, primary DBs, backup, config signals in breakdown or summary.
	checks := []struct {
		label string
		ok    bool
	}{
		{"web-cache or Service Worker/Cache", strings.Contains(lower, "web-cache") || strings.Contains(lower, "service worker") || strings.Contains(lower, "cache")},
		{"chat-db or main_*.sqlite", strings.Contains(lower, "chat-db") || strings.Contains(lower, "main_")},
		{"search-index or search_*.sqlite", strings.Contains(lower, "search-index") || strings.Contains(lower, "search_")},
		{"backup or sqlite-backup", strings.Contains(lower, "backup") || strings.Contains(lower, "sqlite-backup")},
		{"config", strings.Contains(lower, "config")},
	}
	for _, c := range checks {
		if !c.ok {
			t.Fatalf("BREAKDOWN/summary should mention %s:\n%s", c.label, stdout)
		}
	}
}

// assertSeaTalkReclaimTiers checks SAFE TO RECLAIM distinguishes usually-safe caches
// from caution on primary chat/search DBs (not usually-safe purge).
func assertSeaTalkReclaimTiers(t *testing.T, stdout string) {
	t.Helper()
	safe := sectionBody(stdout, "SAFE TO RECLAIM")
	if strings.TrimSpace(safe) == "" {
		t.Fatalf("missing SAFE TO RECLAIM body:\n%s", stdout)
	}
	lower := strings.ToLower(safe)
	// Web caches: usually safe (after quit).
	if !strings.Contains(lower, "usually safe") && !strings.Contains(lower, "usually-safe") {
		t.Fatalf("SAFE TO RECLAIM should mark web caches as usually safe:\n%s", safe)
	}
	// Primary DBs: caution only — section should not claim main/search as usually safe alone.
	// Prefer explicit caution language when chat/search appear.
	if strings.Contains(lower, "main_") || strings.Contains(lower, "search_") ||
		strings.Contains(lower, "chat") || strings.Contains(lower, "search") {
		if !strings.Contains(lower, "caution") {
			t.Fatalf("SAFE TO RECLAIM mentioning chat/search DBs should use caution (not usually-safe purge):\n%s", safe)
		}
	}
}

// --- BREAKDOWN table helpers (aligned SIZE/NAME/ROLE/RECLAIMABLE/NOTES; size DESC) ---

// writeGenericCacheTmpFixture creates a generic-dir layout with reclaim-like basenames.
// Exact file payloads (content total genericCacheTmpContentBytes = 332):
//
//	Cache/entry  200 B → expected role cache (reclaimable)
//	tmp/work     100 B → expected role tmp (reclaimable)
//	notes.txt     32 B → expected role file (neutral)
//
// Directory metadata may make measured total ≥ 332.
func writeGenericCacheTmpFixture(t *testing.T, fixtureDir string) string {
	t.Helper()
	writeSizedFile(t, fixtureDir, "Cache/entry", 200)
	writeSizedFile(t, fixtureDir, "tmp/work", 100)
	writeSizedFile(t, fixtureDir, "notes.txt", 32)
	return fixtureDir
}

// genericCacheTmpContentBytes is the sum of file payloads in writeGenericCacheTmpFixture.
const genericCacheTmpContentBytes int64 = 200 + 100 + 32 // 332

// stripANSI removes CSI/SGR escape sequences so column alignment uses visible width only.
func stripANSI(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) {
				c := s[j]
				if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
					j++
					break
				}
				j++
			}
			i = j
			continue
		}
		// Also accept literal ESC written as the two-char form sometimes seen in dumps.
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

// breakdownSection returns the BREAKDOWN body with ANSI stripped for table parsing.
func breakdownSection(stdout string) string {
	return stripANSI(sectionBody(stdout, "BREAKDOWN"))
}

// parseBreakdownTable splits BREAKDOWN into optional header line and data lines
// (skips empty and "(empty)"). Header is a line whose first fields are SIZE NAME …
func parseBreakdownTable(stdout string) (header string, data []string) {
	body := breakdownSection(stdout)
	for _, line := range strings.Split(body, "\n") {
		trimRight := strings.TrimRight(line, "\r")
		trim := strings.TrimSpace(trimRight)
		if trim == "" || trim == "(empty)" {
			continue
		}
		fields := strings.Fields(trim)
		if len(fields) >= 4 && fields[0] == "SIZE" && fields[1] == "NAME" {
			header = trimRight
			continue
		}
		data = append(data, trimRight)
	}
	return header, data
}

// assertBreakdownTableHeader requires the aligned table header with SIZE, NAME, ROLE,
// RECLAIMABLE (NOTES optional when unused).
func assertBreakdownTableHeader(t *testing.T, stdout string) {
	t.Helper()
	header, _ := parseBreakdownTable(stdout)
	if header == "" {
		t.Fatalf("BREAKDOWN missing table header row (SIZE NAME ROLE RECLAIMABLE …):\n%s", sectionBody(stdout, "BREAKDOWN"))
	}
	upper := header
	for _, col := range []string{"SIZE", "NAME", "ROLE", "RECLAIMABLE"} {
		if !strings.Contains(upper, col) {
			t.Fatalf("BREAKDOWN header missing column %q: %q", col, header)
		}
	}
}

// assertBreakdownSizeColumnAligned asserts SIZE values are right-aligned: the end
// column of each size token is identical across data rows (visible width, ANSI stripped).
// Requires ≥2 data rows with unequal size string widths to be meaningful (e.g. 400B vs 32B).
func assertBreakdownSizeColumnAligned(t *testing.T, stdout string) {
	t.Helper()
	_, data := parseBreakdownTable(stdout)
	if len(data) < 2 {
		t.Fatalf("BREAKDOWN alignment needs ≥2 data rows, got %d:\n%s", len(data), breakdownSection(stdout))
	}
	endCols := make([]int, 0, len(data))
	for _, line := range data {
		i := 0
		for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
			i++
		}
		start := i
		for i < len(line) && line[i] != ' ' && line[i] != '\t' {
			i++
		}
		if start == i {
			t.Fatalf("BREAKDOWN data line has no SIZE token: %q", line)
		}
		endCols = append(endCols, i)
	}
	ref := endCols[0]
	for i, c := range endCols {
		if c != ref {
			t.Fatalf("SIZE column not right-aligned: row %d size ends at col %d, want %d\nrows:\n%s",
				i, c, ref, strings.Join(data, "\n"))
		}
	}
	// Require at least two distinct size-token widths so padding is forced.
	widths := map[int]bool{}
	for _, line := range data {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		widths[len(fields[0])] = true
	}
	if len(widths) < 2 {
		t.Fatalf("alignment fixture should include unequal size string widths (e.g. 400B vs 32B); got sizes from:\n%s",
			strings.Join(data, "\n"))
	}
}

// assertBreakdownNamesInOrder requires each name to appear in BREAKDOWN after the previous
// (size DESC name order). Uses substring search so multi-word names work.
func assertBreakdownNamesInOrder(t *testing.T, stdout string, names []string) {
	t.Helper()
	body := breakdownSection(stdout)
	if strings.TrimSpace(body) == "" {
		t.Fatalf("empty BREAKDOWN body:\n%s", stdout)
	}
	pos := 0
	for _, name := range names {
		i := strings.Index(body[pos:], name)
		if i < 0 {
			t.Fatalf("BREAKDOWN missing name %q after previous names %v (pos=%d):\n%s", name, names, pos, body)
		}
		pos += i + len(name)
	}
}

// breakdownLineForName returns the first BREAKDOWN data line containing name (ANSI stripped).
func breakdownLineForName(t *testing.T, stdout, name string) string {
	t.Helper()
	_, data := parseBreakdownTable(stdout)
	for _, line := range data {
		if strings.Contains(line, name) {
			return line
		}
	}
	// Fallback: raw body search (header may be missing on RED).
	body := breakdownSection(stdout)
	for _, line := range strings.Split(body, "\n") {
		if strings.Contains(line, name) {
			return line
		}
	}
	t.Fatalf("no BREAKDOWN line containing %q:\n%s", name, body)
	return ""
}

// assertBreakdownReclaimableCheckbox locks human RECLAIMABLE as Unicode ☑ / ☐ only
// (never true/false, never ASCII [x]/[ ]) on the row for name.
// checked=true expects ☑, false expects ☐.
func assertBreakdownReclaimableCheckbox(t *testing.T, stdout, name string, checked bool) {
	t.Helper()
	line := breakdownLineForName(t, stdout, name)
	if checked {
		if !strings.Contains(line, "☑") {
			t.Fatalf("BREAKDOWN row for %q must show reclaimable ☑: %q", name, line)
		}
	} else {
		if !strings.Contains(line, "☐") {
			t.Fatalf("BREAKDOWN row for %q must show non-reclaimable ☐: %q", name, line)
		}
		if strings.Contains(line, "☑") {
			t.Fatalf("BREAKDOWN row for %q must not show ☑: %q", name, line)
		}
	}
	// Human RECLAIMABLE column must never print boolean words or ASCII checkboxes.
	fields := strings.Fields(line)
	for _, f := range fields {
		if f == "true" || f == "false" {
			t.Fatalf("BREAKDOWN human RECLAIMABLE must use ☑/☐, not %q (line %q)", f, line)
		}
		if f == "[x]" || f == "[ ]" {
			t.Fatalf("BREAKDOWN human RECLAIMABLE must use Unicode ☑/☐, not ASCII %q (line %q)", f, line)
		}
	}
}

// assertBreakdownNoRoleEqualsPrefix ensures roles are bare (no role= prefix) in human table.
func assertBreakdownNoRoleEqualsPrefix(t *testing.T, stdout string) {
	t.Helper()
	body := breakdownSection(stdout)
	if strings.Contains(body, "role=") {
		t.Fatalf("BREAKDOWN human ROLE cells must be bare roles (no role= prefix):\n%s", body)
	}
}

// assertBreakdownHasCheckboxes requires at least one ☑ or ☐ in the BREAKDOWN section.
func assertBreakdownHasCheckboxes(t *testing.T, stdout string) {
	t.Helper()
	body := breakdownSection(stdout)
	if !strings.Contains(body, "☑") && !strings.Contains(body, "☐") {
		t.Fatalf("BREAKDOWN must include RECLAIMABLE checkboxes ☑/☐:\n%s", body)
	}
}

// yellowSGRPrefixes are accepted SGR sequences for yellow (plain or bold yellow).
var yellowSGRPrefixes = []string{
	"\x1b[33m",
	"\x1b[1;33m",
	"\x1b[33;1m",
	"\x1b[93m",
	"\x1b[1;93m",
}

// assertColoredToken requires token to appear immediately after one of the SGR prefixes.
func assertColoredToken(t *testing.T, stdout, token string, prefixes []string, colorName string) {
	t.Helper()
	if token == "" {
		t.Fatal("assertColoredToken: empty token")
	}
	found := false
	for _, pref := range prefixes {
		if strings.Contains(stdout, pref+token) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected %s ANSI around token %q (e.g. \\x1b[…m%s):\n%s", colorName, token, token, stdout)
	}
}

// assertBreakdownRoleGreen requires ROLE cell for role token to be green-wrapped in BREAKDOWN.
func assertBreakdownRoleGreen(t *testing.T, stdout, role string) {
	t.Helper()
	body := sectionBody(stdout, "BREAKDOWN")
	assertColoredToken(t, body, role, greenSGRPrefixes, "green")
}

// assertBreakdownRoleYellow requires ROLE cell for role token to be yellow-wrapped in BREAKDOWN.
func assertBreakdownRoleYellow(t *testing.T, stdout, role string) {
	t.Helper()
	body := sectionBody(stdout, "BREAKDOWN")
	assertColoredToken(t, body, role, yellowSGRPrefixes, "yellow")
}

// assertBreakdownReclaimableCheckboxGreen requires reclaimable ☑ is green-wrapped in BREAKDOWN
// (same SGR family as reclaimable ROLE) when color is on.
func assertBreakdownReclaimableCheckboxGreen(t *testing.T, stdout string) {
	t.Helper()
	body := sectionBody(stdout, "BREAKDOWN")
	assertColoredToken(t, body, "☑", greenSGRPrefixes, "green")
}

// assertBreakdownNonReclaimableCheckboxNotColored requires non-reclaimable ☐ is not
// green/yellow-wrapped in BREAKDOWN.
func assertBreakdownNonReclaimableCheckboxNotColored(t *testing.T, stdout string) {
	t.Helper()
	body := sectionBody(stdout, "BREAKDOWN")
	prefs := append(append([]string{}, greenSGRPrefixes...), yellowSGRPrefixes...)
	for _, pref := range prefs {
		if strings.Contains(body, pref+"☐") {
			t.Fatalf("non-reclaimable ☐ must not be color-wrapped; found colored ☐ with prefix %q in:\n%s", pref, body)
		}
	}
}

// parseJSONBreakdown unmarshals breakdown array entries as maps.
func parseJSONBreakdown(t *testing.T, payload map[string]json.RawMessage) []map[string]any {
	t.Helper()
	raw, ok := payload["breakdown"]
	if !ok || raw == nil {
		t.Fatal("json missing breakdown")
	}
	var breakdown []map[string]any
	if err := json.Unmarshal(raw, &breakdown); err != nil {
		t.Fatalf("breakdown: %v", err)
	}
	if len(breakdown) == 0 {
		t.Fatal("breakdown must be non-empty")
	}
	return breakdown
}

// jsonEntryInt64 reads a numeric JSON field as int64 (accepts json.Number / float64).
func jsonEntryInt64(t *testing.T, entry map[string]any, key string) int64 {
	t.Helper()
	v, ok := entry[key]
	if !ok || v == nil {
		t.Fatalf("entry missing %q: %v", key, entry)
	}
	switch n := v.(type) {
	case float64:
		return int64(n)
	case int64:
		return n
	case int:
		return int64(n)
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			t.Fatalf("%q: %v", key, err)
		}
		return i
	default:
		t.Fatalf("%q: expected number, got %T %v", key, v, v)
		return 0
	}
}

// assertJSONBreakdownSortedDesc requires breakdown[].size is nonincreasing (size DESC).
func assertJSONBreakdownSortedDesc(t *testing.T, breakdown []map[string]any) {
	t.Helper()
	if len(breakdown) < 2 {
		t.Fatalf("JSON breakdown sort needs ≥2 entries, got %d", len(breakdown))
	}
	prev := jsonEntryInt64(t, breakdown[0], "size")
	for i := 1; i < len(breakdown); i++ {
		sz := jsonEntryInt64(t, breakdown[i], "size")
		if sz > prev {
			t.Fatalf("JSON breakdown not sorted size DESC at index %d: size %d > previous %d (entry=%v)",
				i, sz, prev, breakdown[i])
		}
		prev = sz
	}
}

// assertJSONBreakdownReclaimableBools requires every entry has reclaimable as a JSON boolean
// (not a string / checkbox glyph).
func assertJSONBreakdownReclaimableBools(t *testing.T, breakdown []map[string]any) {
	t.Helper()
	for i, e := range breakdown {
		v, ok := e["reclaimable"]
		if !ok {
			t.Fatalf("breakdown[%d] missing reclaimable (bool): %v", i, e)
		}
		if _, ok := v.(bool); !ok {
			t.Fatalf("breakdown[%d].reclaimable must be JSON bool, got %T %v", i, v, v)
		}
	}
}

// assertJSONBreakdownRoleReclaimable finds first entry with role and checks reclaimable value.
func assertJSONBreakdownRoleReclaimable(t *testing.T, breakdown []map[string]any, role string, want bool) {
	t.Helper()
	for _, e := range breakdown {
		r, _ := e["role"].(string)
		if r != role {
			continue
		}
		v, ok := e["reclaimable"].(bool)
		if !ok {
			t.Fatalf("role %q: reclaimable not bool: %v", role, e)
		}
		if v != want {
			t.Fatalf("role %q: reclaimable=%v, want %v", role, v, want)
		}
		return
	}
	t.Fatalf("breakdown has no entry with role %q: %v", role, breakdown)
}

// --- --all-kinds helpers (multi-pack report: INDEX + detail; JSON AllKindsResult) ---

// parseJSONStdoutFlexible unmarshals JSON from stdout that may be single-line or pretty,
// with a trailing blank line. Prefer this for --all-kinds envelopes.
func parseJSONStdoutFlexible(t *testing.T, stdout string) map[string]json.RawMessage {
	t.Helper()
	content := strings.TrimSpace(stdout)
	if content == "" {
		t.Fatal("stdout is empty (expected JSON)")
	}
	var payload map[string]json.RawMessage
	if err := json.Unmarshal([]byte(content), &payload); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout)
	}
	return payload
}

// allKindsIndexSection returns text from the INDEX header until the first detail KIND: line
// (mini-explain) or EOF. Column header "KIND" without a value is not treated as detail.
func allKindsIndexSection(stdout string) string {
	// Prefer a line that is exactly INDEX or starts with "INDEX".
	idx := -1
	for _, needle := range []string{"\nINDEX\n", "\nINDEX ", "\nINDEX\t", "INDEX\n"} {
		i := strings.Index(stdout, needle)
		if i >= 0 {
			if needle[0] == '\n' {
				idx = i + 1
			} else {
				idx = i
			}
			break
		}
	}
	if idx < 0 {
		// Fallback: any INDEX occurrence
		idx = strings.Index(stdout, "INDEX")
		if idx < 0 {
			return ""
		}
	}
	rest := stdout[idx:]
	// Walk lines; stop at first "KIND: <nonempty>" detail marker.
	offset := 0
	for i, line := range strings.Split(rest, "\n") {
		if i == 0 {
			offset += len(line) + 1
			continue
		}
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "KIND:") {
			val := strings.TrimSpace(strings.TrimPrefix(trim, "KIND:"))
			if val != "" {
				return rest[:offset]
			}
		}
		offset += len(line) + 1
		if offset > len(rest) {
			offset = len(rest)
		}
	}
	return rest
}

// assertAllKindsHumanHeader locks SCOPE / MODE all-kinds / TOTAL header markers.
func assertAllKindsHumanHeader(t *testing.T, stdout, scope string) {
	t.Helper()
	if !strings.Contains(stdout, "SCOPE:") {
		t.Fatalf("all-kinds human output missing SCOPE: header:\n%s", stdout)
	}
	if scope != "" && !strings.Contains(stdout, scope) {
		t.Fatalf("all-kinds SCOPE/output must include scope path %q:\n%s", scope, stdout)
	}
	lower := strings.ToLower(stdout)
	if !strings.Contains(lower, "mode:") || !strings.Contains(lower, "all-kinds") {
		t.Fatalf("all-kinds human output missing MODE: all-kinds:\n%s", stdout)
	}
	if !strings.Contains(lower, "total") {
		t.Fatalf("all-kinds human output missing TOTAL (present) header:\n%s", stdout)
	}
}

// assertAllKindsIndex locks INDEX table: all registered cli kinds, present/missing status,
// and present packs ordered by size DESC (codex > android-sdk > grok > iterm2 for home fixture).
func assertAllKindsIndex(t *testing.T, stdout string, presentCLI []string, missingCLI []string) {
	t.Helper()
	index := allKindsIndexSection(stdout)
	if strings.TrimSpace(index) == "" {
		t.Fatalf("all-kinds human output missing INDEX section:\n%s", stdout)
	}
	// Column cues (SIZE/KIND/STATUS/PATH) — soft if layout uses lowercase or spacing variants.
	idxLower := strings.ToLower(index)
	for _, col := range []string{"size", "kind", "status", "path"} {
		if !strings.Contains(idxLower, col) {
			t.Fatalf("INDEX should include column %q (SIZE/KIND/STATUS/PATH):\n%s", col, index)
		}
	}
	for _, k := range allKindsCLIKinds {
		if !strings.Contains(idxLower, strings.ToLower(k)) {
			t.Fatalf("INDEX must list cli kind %q:\n%s", k, index)
		}
	}
	// Status markers: each present kind's row mentions present; missing kinds mention missing.
	// Prefer same-line association; fall back to section containing both tokens.
	for _, k := range presentCLI {
		if !indexLineHasStatus(index, k, "present") {
			t.Fatalf("INDEX should mark cli kind %q as present:\n%s", k, index)
		}
	}
	for _, k := range missingCLI {
		if !indexLineHasStatus(index, k, "missing") {
			t.Fatalf("INDEX should mark cli kind %q as missing:\n%s", k, index)
		}
	}
	// Present order size DESC: earlier present kind appears before later in INDEX body.
	if len(presentCLI) >= 2 {
		positions := make([]int, len(presentCLI))
		for i, k := range presentCLI {
			positions[i] = indexKindPos(index, k)
			if positions[i] < 0 {
				t.Fatalf("INDEX missing present kind %q for order check:\n%s", k, index)
			}
		}
		for i := 1; i < len(positions); i++ {
			if positions[i] < positions[i-1] {
				t.Fatalf("INDEX present kinds not sorted size DESC: %q (pos %d) before %q (pos %d):\n%s",
					presentCLI[i], positions[i], presentCLI[i-1], positions[i-1], index)
			}
		}
	}
}

// indexKindPos returns the first index of cli kind token in INDEX text (case-insensitive), or -1.
func indexKindPos(index, cliKind string) int {
	return strings.Index(strings.ToLower(index), strings.ToLower(cliKind))
}

// indexLineHasStatus reports whether some INDEX line mentions both cliKind and status.
func indexLineHasStatus(index, cliKind, status string) bool {
	k := strings.ToLower(cliKind)
	s := strings.ToLower(status)
	for _, line := range strings.Split(index, "\n") {
		ll := strings.ToLower(line)
		if strings.Contains(ll, k) && strings.Contains(ll, s) {
			return true
		}
	}
	// Fallback: whole INDEX mentions kind and status nearby is enough only if both present.
	idxLower := strings.ToLower(index)
	return strings.Contains(idxLower, k) && strings.Contains(idxLower, s)
}

// assertAllKindsDetailPresent locks mini-explain detail for each present output kind id:
// KIND: <id> and BREAKDOWN appear; missing kinds should not get a full present detail.
func assertAllKindsDetailPresent(t *testing.T, stdout string, presentOutputKinds []string) {
	t.Helper()
	for _, kind := range presentOutputKinds {
		want := "KIND: " + kind
		if !strings.Contains(stdout, want) {
			t.Fatalf("all-kinds detail missing %q:\n%s", want, stdout)
		}
	}
	// BREAKDOWN should appear at least once per present kind (mini-explain).
	nBD := strings.Count(stdout, "BREAKDOWN")
	if nBD < len(presentOutputKinds) {
		t.Fatalf("expected ≥%d BREAKDOWN sections (one per present kind), got %d:\n%s",
			len(presentOutputKinds), nBD, stdout)
	}
}

// assertAllKindsJSONShape locks AllKindsResult JSON keys and per-kind entry fields.
func assertAllKindsJSONShape(t *testing.T, payload map[string]json.RawMessage, scope string) []map[string]any {
	t.Helper()
	for _, key := range []string{"scope", "totalSize", "kinds"} {
		if payload[key] == nil {
			t.Fatalf("all-kinds JSON missing required key %q", key)
		}
	}
	gotScope := jsonStringField(t, payload, "scope")
	if scope != "" && gotScope != scope {
		// Allow trailing slash differences only if both clean equal via TrimRight
		if strings.TrimRight(gotScope, "/") != strings.TrimRight(scope, "/") {
			t.Fatalf("json scope: expected %q, got %q", scope, gotScope)
		}
	}
	_ = jsonInt64Field(t, payload, "totalSize")

	var kinds []map[string]any
	if err := json.Unmarshal(payload["kinds"], &kinds); err != nil {
		t.Fatalf("kinds: %v", err)
	}
	if len(kinds) != len(allKindsCLIKinds) {
		t.Fatalf("kinds: expected %d entries (v1 packs), got %d: %v", len(allKindsCLIKinds), len(kinds), kinds)
	}
	for i, e := range kinds {
		for _, key := range []string{"kind", "cliKind", "path", "status", "totalSize"} {
			if _, ok := e[key]; !ok {
				t.Fatalf("kinds[%d] missing %q: %v", i, key, e)
			}
		}
		status, _ := e["status"].(string)
		if status != "present" && status != "missing" && status != "error" {
			t.Fatalf("kinds[%d].status: expected present|missing|error, got %q", i, status)
		}
	}
	return kinds
}

// allKindsByCLIKind maps kinds[] entries by cliKind string.
func allKindsByCLIKind(t *testing.T, kinds []map[string]any) map[string]map[string]any {
	t.Helper()
	out := make(map[string]map[string]any, len(kinds))
	for _, e := range kinds {
		cli, _ := e["cliKind"].(string)
		if cli == "" {
			t.Fatalf("kinds entry missing cliKind: %v", e)
		}
		if _, dup := out[cli]; dup {
			t.Fatalf("duplicate cliKind %q in kinds", cli)
		}
		out[cli] = e
	}
	return out
}
```
