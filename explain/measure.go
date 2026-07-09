package explain

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"disk-usage-analyser/usagescan"
)

// buildExplanation measures sizes and fills kind-specific summary/breakdown/reclaim/raw.
func buildExplanation(targetPath string, info os.FileInfo, det detectResult) (Explanation, error) {
	measureRoot := det.ContextRoot
	if measureRoot == "" {
		measureRoot = targetPath
	}

	var totalSize int64
	var breakdown []Breakdown
	var err error

	switch det.Kind {
	case "android-avd":
		totalSize, breakdown, err = measureAndroidAVD(measureRoot)
	case "android-sdk":
		totalSize, breakdown, err = measureAndroidSDK(measureRoot)
	case "seatalk-app-support":
		totalSize, breakdown, err = measureSeaTalkAppSupport(measureRoot)
	case "grok-home":
		totalSize, breakdown, err = measureGrokHome(measureRoot)
	case "codex-home":
		totalSize, breakdown, err = measureCodexHome(measureRoot)
	case "iterm2":
		totalSize, breakdown, err = measureITerm2(measureRoot)
	case "go-build-cache", "npm-cache", "homebrew-cache", "generic-dir":
		totalSize, breakdown, err = measureDir(measureRoot, 3)
	case "generic-qcow2", "generic-file":
		totalSize, breakdown, err = measureFile(targetPath)
	default:
		if info.IsDir() {
			totalSize, breakdown, err = measureDir(measureRoot, 3)
		} else {
			totalSize, breakdown, err = measureFile(targetPath)
		}
	}
	if err != nil {
		return Explanation{}, err
	}

	// Prefer reporting the user-requested path in PATH/json.path.
	exp := Explanation{
		Path:       targetPath,
		Kind:       det.Kind,
		TotalSize:  totalSize,
		Confidence: det.Confidence,
		Breakdown:  finalizeBreakdown(breakdown),
	}
	exp.Summary = summaryFor(det.Kind, targetPath)
	exp.Reclaim = reclaimFor(det.Kind)
	exp.HowToPurge = howToPurgeFor(det.Kind, targetPath, measureRoot)
	exp.RawCommands = rawCommandsFor(targetPath, det.Kind)
	// Product shape A: always attach LOGS DB when logs_*.sqlite is readable.
	if det.Kind == "codex-home" {
		if logsDB, lerr := readCodexLogsDB(measureRoot); lerr == nil && logsDB != nil {
			exp.LogsDB = logsDB
		}
	}
	return exp, nil
}

// roleTier classifies a breakdown role for reclaim checkbox + ROLE cell color.
type roleTier int

const (
	tierNeutral roleTier = iota
	tierReclaimable
	tierCaution
)

// roleTierOf maps a role string to reclaim / caution / neutral.
func roleTierOf(role string) roleTier {
	switch role {
	case "cache", "web-cache", "tmp", "temp", "snapshot", "snapshots", "backup",
		"derived-data", "docs-cache", "device-support", "simulator",
		"installer-cache", "session-logs", "logs", "app-logs-db", "history",
		"system-images", "skins", "sources",
		"python-env", "python-env-alias":
		return tierReclaimable
	case "chat-db", "search-index", "user-data", "sdcard", "config", "session", "web-storage",
		"archives", "project-meta", "app-data", "runtime-db", "app-state-db",
		"build-tools", "platforms", "emulator", "platform-tools", "cmdline-tools", "licenses", "meta",
		"user-config", "helper-binary", "app-db", "state":
		return tierCaution
	default:
		return tierNeutral
	}
}

// isKnownKindPack reports whether kind is a registered forced pack/alias
// (v1: xcode, grok, android-sdk, iterm2, codex).
func isKnownKindPack(kind string) bool {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "xcode", "grok", "android-sdk", "iterm2", "codex":
		return true
	default:
		return false
	}
}

// kindPackRoot describes one relative root under a multi-root pack scope.
type kindPackRoot struct {
	Rel         string // relative path under scope (slash-separated)
	Role        string
	Notes       string
	Reclaimable bool // preferred reclaim signal; finalizeBreakdown still applies role tier
}

// xcodePackRoots is the ordered multi-root Xcode developer pack (frontend / tmp_analyse).
var xcodePackRoots = []kindPackRoot{
	{Rel: "Library/Developer/Xcode/DerivedData", Role: "derived-data", Notes: "rebuildable build products", Reclaimable: true},
	{Rel: "Library/Developer/CoreSimulator/Devices", Role: "simulator", Notes: "simulator devices (wipe caution)", Reclaimable: true},
	{Rel: "Library/Developer/Xcode/iOS DeviceSupport", Role: "device-support", Notes: "re-download on device reconnect", Reclaimable: true},
	{Rel: "Library/Developer/Xcode/Archives", Role: "archives", Notes: "signed builds — not usually safe", Reclaimable: false},
	{Rel: "Library/Developer/Xcode/DocumentationCache", Role: "docs-cache", Notes: "re-fetched documentation cache", Reclaimable: true},
}

// buildKindPackExplanation measures a forced multi-root pack under absolute scope.
func buildKindPackExplanation(kind, scope string) (Explanation, error) {
	kind = strings.ToLower(strings.TrimSpace(kind))
	var totalSize int64
	var breakdown []Breakdown
	var err error

	switch kind {
	case "xcode":
		totalSize, breakdown, err = measureKindPack(scope, xcodePackRoots)
	default:
		return Explanation{}, fmt.Errorf("unknown kind %q (supported: xcode, grok, android-sdk, iterm2, codex)", kind)
	}
	if err != nil {
		return Explanation{}, err
	}

	confidence := "high"
	if len(breakdown) == 0 {
		// Still success with empty pack under scope.
		confidence = "medium"
	}

	exp := Explanation{
		Path:       scope,
		Kind:       kind,
		TotalSize:  totalSize,
		Confidence: confidence,
		Breakdown:  finalizeBreakdown(breakdown),
	}
	exp.Summary = summaryFor(kind, scope)
	exp.Reclaim = reclaimFor(kind)
	exp.HowToPurge = howToPurgeFor(kind, scope, scope)
	exp.RawCommands = rawCommandsFor(scope, kind)
	return exp, nil
}

// buildAndroidSDKKindExplanation resolves --kind android-sdk under absolute scope.
// If scope is already an SDK root (path form or signatures), measure it; otherwise
// measure {scope}/Library/Android/sdk (must exist). Kind id is android-sdk.
func buildAndroidSDKKindExplanation(scope string) (Explanation, error) {
	target := scope
	if !isAndroidSDKRoot(scope) {
		target = filepath.Join(scope, filepath.FromSlash("Library/Android/sdk"))
	}
	info, err := os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			return Explanation{}, fmt.Errorf("path does not exist: %s", target)
		}
		return Explanation{}, fmt.Errorf("cannot access %s: %w", target, err)
	}
	if !info.IsDir() {
		return Explanation{}, fmt.Errorf("android-sdk must be a directory: %s", target)
	}
	det := detectResult{
		Kind:         "android-sdk",
		Confidence:   "high",
		ContextRoot:  target,
		TargetIsFile: false,
	}
	return buildExplanation(target, info, det)
}

// buildGrokKindExplanation resolves --kind grok under absolute scope and measures as grok-home.
// If scope basename is .grok, measure it; otherwise measure {scope}/.grok (must exist).
func buildGrokKindExplanation(scope string) (Explanation, error) {
	target := scope
	if filepath.Base(scope) != ".grok" {
		target = filepath.Join(scope, ".grok")
	}
	info, err := os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			return Explanation{}, fmt.Errorf("path does not exist: %s", target)
		}
		return Explanation{}, fmt.Errorf("cannot access %s: %w", target, err)
	}
	if !info.IsDir() {
		return Explanation{}, fmt.Errorf("grok home must be a directory: %s", target)
	}
	det := detectResult{
		Kind:         "grok-home",
		Confidence:   "high",
		ContextRoot:  target,
		TargetIsFile: false,
	}
	return buildExplanation(target, info, det)
}

// buildCodexKindExplanation resolves --kind codex under absolute scope and measures as codex-home.
// If scope basename is .codex, measure it; otherwise measure {scope}/.codex (must exist).
func buildCodexKindExplanation(scope string) (Explanation, error) {
	target := scope
	if filepath.Base(scope) != ".codex" {
		target = filepath.Join(scope, ".codex")
	}
	info, err := os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			return Explanation{}, fmt.Errorf("path does not exist: %s", target)
		}
		return Explanation{}, fmt.Errorf("cannot access %s: %w", target, err)
	}
	if !info.IsDir() {
		return Explanation{}, fmt.Errorf("codex home must be a directory: %s", target)
	}
	det := detectResult{
		Kind:         "codex-home",
		Confidence:   "high",
		ContextRoot:  target,
		TargetIsFile: false,
	}
	return buildExplanation(target, info, det)
}

// buildITerm2KindExplanation resolves --kind iterm2 under absolute scope.
// If scope is already an iTerm2 root (path form or signatures), measure it; otherwise
// measure {scope}/Library/Application Support/iTerm2 (must exist). Kind id is iterm2.
func buildITerm2KindExplanation(scope string) (Explanation, error) {
	target := scope
	if !isITerm2Root(scope) {
		target = filepath.Join(scope, filepath.FromSlash("Library/Application Support/iTerm2"))
	}
	info, err := os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			return Explanation{}, fmt.Errorf("path does not exist: %s", target)
		}
		return Explanation{}, fmt.Errorf("cannot access %s: %w", target, err)
	}
	if !info.IsDir() {
		return Explanation{}, fmt.Errorf("iterm2 must be a directory: %s", target)
	}
	det := detectResult{
		Kind:         "iterm2",
		Confidence:   "high",
		ContextRoot:  target,
		TargetIsFile: false,
	}
	return buildExplanation(target, info, det)
}

// measureKindPack walks each existing relative root under scope; omits missing roots.
func measureKindPack(scope string, roots []kindPackRoot) (int64, []Breakdown, error) {
	var total int64
	var breakdown []Breakdown
	for _, r := range roots {
		p := filepath.Join(scope, filepath.FromSlash(r.Rel))
		st, err := os.Stat(p)
		if err != nil {
			continue // omit missing roots
		}
		var size int64
		if st.IsDir() {
			size, err = pathSize(p)
			if err != nil {
				continue
			}
		} else {
			size = st.Size()
		}
		total += size
		breakdown = append(breakdown, Breakdown{
			Name:  filepath.Base(p),
			Path:  p,
			Size:  size,
			Role:  r.Role,
			Notes: r.Notes,
		})
	}
	return total, breakdown, nil
}

// isReclaimableRole reports whether role is the reclaimable tier (☑ / reclaimable:true).
func isReclaimableRole(role string) bool {
	return roleTierOf(role) == tierReclaimable
}

// genericBasenameRole remaps top-level child basenames that would otherwise be only
// directory/file into semantic reclaim roles (case-insensitive). Specialized kinds
// keep their own role assignment and should not call this for those roles.
func genericBasenameRole(name string, isDir bool) string {
	lower := strings.ToLower(name)
	switch lower {
	case "tmp", "temp", ".tmp":
		return "tmp"
	case "cache", "caches":
		return "cache"
	case "snapshots":
		return "snapshot"
	case "sqlite-backup", "idb-backup":
		return "backup"
	}
	// *Cache / *_cacache (e.g. GPUCache, Code Cache, _cacache)
	if strings.HasSuffix(lower, "cache") || strings.HasSuffix(lower, "_cacache") {
		return "cache"
	}
	if isDir {
		return "directory"
	}
	return "file"
}

// finalizeBreakdown sets Reclaimable from role tier and sorts size DESC, name ASC.
// Shared for human table and JSON.
func finalizeBreakdown(breakdown []Breakdown) []Breakdown {
	if len(breakdown) == 0 {
		return breakdown
	}
	for i := range breakdown {
		breakdown[i].Reclaimable = isReclaimableRole(breakdown[i].Role)
	}
	sort.SliceStable(breakdown, func(i, j int) bool {
		if breakdown[i].Size != breakdown[j].Size {
			return breakdown[i].Size > breakdown[j].Size
		}
		return breakdown[i].Name < breakdown[j].Name
	})
	return breakdown
}

func measureFile(path string) (int64, []Breakdown, error) {
	st, err := os.Stat(path)
	if err != nil {
		return 0, nil, err
	}
	size := st.Size()
	return size, []Breakdown{{
		Name: filepath.Base(path),
		Path: path,
		Size: size,
		Role: "file",
	}}, nil
}

// measureDir walks the directory (modest depth) and returns total content size
// plus top-level child breakdown entries.
func measureDir(root string, maxDepth int) (int64, []Breakdown, error) {
	// Use usagescan tree walk for directories when possible.
	result, err := usagescan.ScanTree(root, usagescan.ScanOptions{
		Min:      0,
		MaxDepth: maxDepth,
	})
	if err != nil {
		// Fall back to simple walk (e.g. edge cases).
		return walkSizeAndChildren(root)
	}

	var breakdown []Breakdown
	for _, child := range result.Tree.Children {
		role := genericBasenameRole(child.Name, child.IsDir)
		breakdown = append(breakdown, Breakdown{
			Name: child.Name,
			Path: child.Path,
			Size: child.Size,
			Role: role,
		})
	}
	if len(breakdown) == 0 {
		// Empty or leaf-only: include root itself.
		breakdown = append(breakdown, Breakdown{
			Name: filepath.Base(root),
			Path: root,
			Size: result.TotalSize,
			Role: "directory",
		})
	}
	return result.TotalSize, breakdown, nil
}

func walkSizeAndChildren(root string) (int64, []Breakdown, error) {
	var total int64
	entries, err := os.ReadDir(root)
	if err != nil {
		return 0, nil, err
	}
	var breakdown []Breakdown
	for _, e := range entries {
		p := filepath.Join(root, e.Name())
		size, err := pathSize(p)
		if err != nil {
			continue
		}
		total += size
		role := genericBasenameRole(e.Name(), e.IsDir())
		breakdown = append(breakdown, Breakdown{
			Name: e.Name(),
			Path: p,
			Size: size,
			Role: role,
		})
	}
	return total, breakdown, nil
}

func pathSize(path string) (int64, error) {
	var total int64
	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		total += info.Size()
		return nil
	})
	return total, err
}

// measureAndroidSDK walks top-level children of an Android SDK root and assigns roles.
func measureAndroidSDK(sdkDir string) (int64, []Breakdown, error) {
	entries, err := os.ReadDir(sdkDir)
	if err != nil {
		return 0, nil, err
	}
	var total int64
	var breakdown []Breakdown
	for _, e := range entries {
		p := filepath.Join(sdkDir, e.Name())
		size, err := pathSize(p)
		if err != nil {
			continue
		}
		if !e.IsDir() {
			if st, err := e.Info(); err == nil {
				size = st.Size()
			}
		}
		total += size
		role, note := androidSDKRoleFor(e.Name(), e.IsDir())
		breakdown = append(breakdown, Breakdown{
			Name:  e.Name(),
			Path:  p,
			Size:  size,
			Role:  role,
			Notes: note,
		})
	}
	if len(breakdown) == 0 {
		return measureDir(sdkDir, 3)
	}
	return total, breakdown, nil
}

// androidSDKRoleFor maps a top-level Android SDK basename to a reclaim role.
func androidSDKRoleFor(name string, isDir bool) (role, note string) {
	switch name {
	case "system-images":
		return "system-images", "Emulator system images; list/uninstall via sdkmanager"
	case ".temp", ".downloadIntermediates":
		return "tmp", "Download temp / incomplete packages"
	case "skins":
		return "skins", "Emulator skins (optional)"
	case "sources":
		return "sources", "Platform sources (optional)"
	case "emulator":
		return "emulator", "Emulator package — keep for device runs"
	case "build-tools":
		return "build-tools", "Keep versions projects need"
	case "platforms":
		return "platforms", "Compile SDK platform levels"
	case "cmdline-tools":
		return "cmdline-tools", "sdkmanager and related CLI tools"
	case "platform-tools":
		return "platform-tools", "adb and platform tools — keep"
	case "licenses":
		return "licenses", "Accepted SDK licenses"
	case ".knownPackages":
		return "meta", "Package index metadata"
	}
	return genericBasenameRole(name, isDir), ""
}

// measureAndroidAVD assigns roles to known AVD artifacts and sums content.
func measureAndroidAVD(avdDir string) (int64, []Breakdown, error) {
	known := []struct {
		rel  string
		role string
		note string
	}{
		{"config.ini", "config", "AVD configuration"},
		{"userdata-qemu.img.qcow2", "user-data", "grown app/data overlay"},
		{"userdata-qemu.img", "user-data", "user data image"},
		{"sdcard.img", "sdcard", "emulated SD card"},
		{"sdcard.img.qcow2", "sdcard", "emulated SD card overlay"},
		{"snapshots", "snapshot", "emulator snapshots (often reclaimable)"},
		{"cache.img", "cache", "emulator cache image"},
		{"cache.img.qcow2", "cache", "emulator cache overlay"},
	}

	var total int64
	var breakdown []Breakdown
	seen := map[string]bool{}

	for _, k := range known {
		p := filepath.Join(avdDir, filepath.FromSlash(k.rel))
		st, err := os.Stat(p)
		if err != nil {
			continue
		}
		var size int64
		if st.IsDir() {
			size, _ = pathSize(p)
		} else {
			size = st.Size()
		}
		total += size
		seen[k.rel] = true
		// Prefer basename for display when nested under snapshots.
		name := filepath.Base(p)
		if k.role == "snapshot" {
			name = "snapshots"
		}
		breakdown = append(breakdown, Breakdown{
			Name:  name,
			Path:  p,
			Size:  size,
			Role:  k.role,
			Notes: k.note,
		})
	}

	// Also pick up any other top-level entries not already counted.
	entries, err := os.ReadDir(avdDir)
	if err != nil {
		return total, breakdown, nil
	}
	for _, e := range entries {
		if seen[e.Name()] {
			continue
		}
		// Skip if we already have userdata variants
		if strings.HasPrefix(e.Name(), "userdata-qemu.img") &&
			(seen["userdata-qemu.img.qcow2"] || seen["userdata-qemu.img"]) {
			continue
		}
		p := filepath.Join(avdDir, e.Name())
		size, err := pathSize(p)
		if err != nil {
			continue
		}
		// Avoid double-counting if pathSize of file vs size already used
		if !e.IsDir() {
			if st, err := e.Info(); err == nil {
				size = st.Size()
			}
		}
		total += size
		role := "other"
		lower := strings.ToLower(e.Name())
		if strings.Contains(lower, "snapshot") {
			role = "snapshot"
		} else if strings.Contains(lower, "userdata") {
			role = "user-data"
		} else if strings.Contains(lower, "sdcard") {
			role = "sdcard"
		}
		breakdown = append(breakdown, Breakdown{
			Name: e.Name(),
			Path: p,
			Size: size,
			Role: role,
		})
	}

	if len(breakdown) == 0 {
		return measureDir(avdDir, 3)
	}
	return total, breakdown, nil
}

// measureSeaTalkAppSupport walks top-level children of the SeaTalk Application
// Support directory and assigns semantic roles (web-cache, chat-db, …).
func measureSeaTalkAppSupport(seatalkDir string) (int64, []Breakdown, error) {
	entries, err := os.ReadDir(seatalkDir)
	if err != nil {
		return 0, nil, err
	}
	var total int64
	var breakdown []Breakdown
	for _, e := range entries {
		p := filepath.Join(seatalkDir, e.Name())
		size, err := pathSize(p)
		if err != nil {
			continue
		}
		if !e.IsDir() {
			if st, err := e.Info(); err == nil {
				size = st.Size()
			}
		}
		total += size
		role, note := seatalkRoleFor(e.Name())
		breakdown = append(breakdown, Breakdown{
			Name:  e.Name(),
			Path:  p,
			Size:  size,
			Role:  role,
			Notes: note,
		})
	}
	if len(breakdown) == 0 {
		return measureDir(seatalkDir, 3)
	}
	return total, breakdown, nil
}

// seatalkRoleFor maps a top-level SeaTalk Application Support basename to a role.
func seatalkRoleFor(name string) (role, note string) {
	// Exact / prefix matches first.
	switch name {
	case "Service Worker", "Cache", "Code Cache", "GPUCache":
		return "web-cache", "Chromium-style web cache (rebuildable after quit)"
	case "blob_storage", "IndexedDB":
		return "web-storage", "Web storage (blobs / IndexedDB)"
	case "sqlite-backup", "idb-backup":
		return "backup", "Local DB backup directory"
	case "Cookies", "Local Storage", "Session Storage":
		return "session", "Browser session/cookie storage"
	case "Preferences", "config.json":
		return "config", "App configuration"
	case "Crashpad", "sentry":
		return "diagnostics", "Crash/diagnostics data"
	}
	if strings.HasPrefix(name, "Dawn") {
		return "web-cache", "Chromium Dawn GPU cache (rebuildable after quit)"
	}
	if strings.HasPrefix(name, "Singleton") {
		return "runtime", "Process singleton lock/runtime"
	}
	// main_*.sqlite (+ -wal / -shm)
	if isSeaTalkNamedSQLite(name, "main_") {
		return "chat-db", "Primary chat database — do not purge as usually-safe"
	}
	// search_*.sqlite (+ -wal / -shm)
	if isSeaTalkNamedSQLite(name, "search_") {
		return "search-index", "Search index database — do not purge as usually-safe"
	}
	return "other", ""
}

// isSeaTalkNamedSQLite reports whether name is prefix*.sqlite or its -wal/-shm siblings.
func isSeaTalkNamedSQLite(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	// Strip -wal / -shm suffixes used by SQLite.
	base := name
	for _, suf := range []string{"-wal", "-shm"} {
		if strings.HasSuffix(base, suf) {
			base = strings.TrimSuffix(base, suf)
			break
		}
	}
	return strings.HasSuffix(base, ".sqlite")
}

// measureCodexHome walks top-level children of the Codex CLI home (.codex) and assigns roles.
func measureCodexHome(codexDir string) (int64, []Breakdown, error) {
	entries, err := os.ReadDir(codexDir)
	if err != nil {
		return 0, nil, err
	}
	var total int64
	var breakdown []Breakdown
	for _, e := range entries {
		p := filepath.Join(codexDir, e.Name())
		size, err := pathSize(p)
		if err != nil {
			continue
		}
		if !e.IsDir() {
			if st, err := e.Info(); err == nil {
				size = st.Size()
			}
		}
		total += size
		role, note := codexRoleFor(e.Name(), e.IsDir())
		breakdown = append(breakdown, Breakdown{
			Name:  e.Name(),
			Path:  p,
			Size:  size,
			Role:  role,
			Notes: note,
		})
	}
	if len(breakdown) == 0 {
		return measureDir(codexDir, 3)
	}
	return total, breakdown, nil
}

// codexRoleFor maps a top-level .codex basename to a reclaim role.
func codexRoleFor(name string, isDir bool) (role, note string) {
	if isCodexLogsDBName(name) {
		return "app-logs-db", "App logs database (caution when reclaiming)"
	}
	// state_*.sqlite, goals_*.sqlite, memories_*.sqlite (+ wal/shm)
	if isCodexStateDBName(name) {
		return "app-state-db", "App state database — not usually-safe"
	}
	// cloud-*-cache.json
	if strings.HasPrefix(name, "cloud-") && strings.HasSuffix(name, "-cache.json") {
		return "cache", "Cloud metadata cache"
	}
	switch name {
	case "sessions":
		return "session-logs", "Session transcripts (caution before purge)"
	case ".tmp":
		return "tmp", "Transient temporary files"
	case "cache":
		return "cache", "Local cache"
	case "models_cache.json":
		return "cache", "Model metadata cache"
	case "shell_snapshots":
		return "snapshots", "Shell snapshots (caution)"
	case "history.jsonl", "session_index.jsonl":
		return "history", "History index (caution before purge)"
	case "config.toml", "auth.json", "hooks.json", "installation_id":
		return "config", "Credentials/settings — not usually-safe"
	case "skills", "plugins", "vendor_imports", "prompts", "rules":
		return "app-data", "Installed app data"
	}
	return genericBasenameRole(name, isDir), ""
}

// isCodexStateDBName reports state_/goals_/memories_*.sqlite (+ wal/shm).
func isCodexStateDBName(name string) bool {
	base := name
	for _, suf := range []string{"-wal", "-shm"} {
		if strings.HasSuffix(base, suf) {
			base = strings.TrimSuffix(base, suf)
			break
		}
	}
	if !strings.HasSuffix(base, ".sqlite") {
		return false
	}
	return strings.HasPrefix(base, "state_") ||
		strings.HasPrefix(base, "goals_") ||
		strings.HasPrefix(base, "memories_")
}

// measureGrokHome walks top-level children of the Grok CLI home (.grok) and assigns roles.
func measureGrokHome(grokDir string) (int64, []Breakdown, error) {
	entries, err := os.ReadDir(grokDir)
	if err != nil {
		return 0, nil, err
	}
	var total int64
	var breakdown []Breakdown
	for _, e := range entries {
		p := filepath.Join(grokDir, e.Name())
		size, err := pathSize(p)
		if err != nil {
			continue
		}
		if !e.IsDir() {
			if st, err := e.Info(); err == nil {
				size = st.Size()
			}
		}
		total += size
		role, note := grokRoleFor(e.Name(), e.IsDir())
		breakdown = append(breakdown, Breakdown{
			Name:  e.Name(),
			Path:  p,
			Size:  size,
			Role:  role,
			Notes: note,
		})
	}
	if len(breakdown) == 0 {
		return measureDir(grokDir, 3)
	}
	return total, breakdown, nil
}

// grokRoleFor maps a top-level .grok basename to a reclaim role.
func grokRoleFor(name string, isDir bool) (role, note string) {
	switch name {
	case "downloads":
		return "installer-cache", "CLI packages and installer leftovers (*.tmp)"
	case "sessions":
		return "session-logs", "Agent transcripts (caution before purge)"
	case "marketplace-cache":
		return "cache", "Plugin marketplace cache"
	case "logs":
		return "logs", "App logs"
	case "models_cache.json":
		return "cache", "Model metadata cache"
	case "upload_queue":
		return "tmp", "Transient upload queue"
	case "projects":
		return "project-meta", "Project state"
	case "skills", "vendor", "bundled", "docs", "completions":
		return "app-data", "Installed app data"
	case "config.toml", "auth.json", "trusted_folders.toml", "agent_id":
		return "config", "Credentials/settings — not usually-safe"
	case "worktrees.db", "active_sessions.json":
		return "runtime-db", "Runtime indexes"
	}
	if strings.HasSuffix(name, ".lock") {
		return "config", "Lock/config file"
	}
	return genericBasenameRole(name, isDir), ""
}

// measureITerm2 walks top-level children of iTerm2 Application Support and assigns roles.
func measureITerm2(iterm2Dir string) (int64, []Breakdown, error) {
	entries, err := os.ReadDir(iterm2Dir)
	if err != nil {
		return 0, nil, err
	}
	var total int64
	var breakdown []Breakdown
	for _, e := range entries {
		p := filepath.Join(iterm2Dir, e.Name())
		size, err := pathSize(p)
		if err != nil {
			continue
		}
		if !e.IsDir() {
			if st, err := e.Info(); err == nil {
				size = st.Size()
			}
		}
		total += size
		role, note := iTerm2RoleFor(e.Name(), e.IsDir())
		breakdown = append(breakdown, Breakdown{
			Name:  e.Name(),
			Path:  p,
			Size:  size,
			Role:  role,
			Notes: note,
		})
	}
	if len(breakdown) == 0 {
		return measureDir(iterm2Dir, 3)
	}
	return total, breakdown, nil
}

// iTerm2RoleFor maps a top-level iTerm2 Application Support basename to a reclaim role.
func iTerm2RoleFor(name string, isDir bool) (role, note string) {
	switch name {
	case "iterm2env":
		return "python-env", "Primary bundled Python environment"
	case "SavedState":
		return "state", "Window/session state"
	case "Scripts", "DynamicProfiles":
		return "user-config", "User scripts/profiles — not usually-safe bulk purge"
	case "version.txt":
		return "meta", "App version metadata"
	case "parsers", "private":
		return "runtime", "Runtime support data"
	}
	// Versioned env aliases: iterm2env-*
	if strings.HasPrefix(name, "iterm2env-") {
		return "python-env-alias", "Versioned Python tree; often hardlinked with iterm2env"
	}
	// Helper binaries: iTermServer-*
	if strings.HasPrefix(name, "iTermServer-") {
		return "helper-binary", "iTerm2 helper binary"
	}
	// App DB: chatdb.sqlite*
	if strings.HasPrefix(name, "chatdb.sqlite") {
		return "app-db", "App chat/state database"
	}
	// Logs: log*.txt or name starts with log and ends with .txt
	if strings.HasSuffix(name, ".txt") && strings.HasPrefix(name, "log") {
		return "logs", "App logs (usually reclaimable after quit)"
	}
	// Runtime: sockets, locks (basename contains or exact-ish)
	lower := strings.ToLower(name)
	if strings.Contains(lower, "socket") || strings.Contains(lower, "lock") {
		return "runtime", "Runtime sockets/locks"
	}
	return genericBasenameRole(name, isDir), ""
}

func summaryFor(kind, path string) []string {
	base := filepath.Base(path)
	switch kind {
	case "android-avd":
		return []string{
			fmt.Sprintf("Android emulator AVD directory (%s)", base),
			"Contains userdata, optional sdcard, and snapshots that often dominate size.",
		}
	case "android-sdk":
		return []string{
			"Android SDK root (system-images, platforms, build-tools, platform-tools, emulator, sources, temp).",
			"system-images, sources, skins, and download temp are reclaimable with care; keep platform-tools (adb), build-tools, platforms, and emulator bulk unless you know they are unused.",
		}
	case "seatalk-app-support":
		return []string{
			"SeaTalk Application Support directory (Chromium-style caches + local chat/search SQLite).",
			"Web caches (Service Worker, Cache, …) are usually rebuildable; primary main_*.sqlite / search_*.sqlite hold chat and search data — reclaim only with caution.",
			"Local backups live under sqlite-backup / idb-backup and are conditional reclaim targets.",
		}
	case "grok-home":
		return []string{
			"Grok / xAI CLI home directory (.grok) with installer downloads, session logs, marketplace cache, and config.",
			"downloads (installer-cache), marketplace-cache, logs, and models_cache are usually reclaimable; sessions are caution; auth/config are not usually-safe purge targets.",
		}
	case "codex-home":
		return []string{
			"OpenAI Codex CLI home directory (.codex) with app logs DB (logs_*.sqlite), sessions, cache, tmp, and config.",
			"app-logs-db, sessions (session-logs), cache, and .tmp are usually reclaimable with caution; auth/config are not usually-safe purge targets.",
			"When logs_*.sqlite is present, a logs database preview shows row count and the last 3 samples (newest first).",
		}
	case "iterm2":
		return []string{
			"iTerm2 Application Support directory (bundled Python envs, logs, profiles, helper binaries).",
			"iterm2env (python-env), iterm2env-* aliases (python-env-alias), and log*.txt are usually reclaimable; DynamicProfiles/Scripts (user-config) are not usually-safe bulk purge.",
			"Multiple iterm2env* trees may share APFS hardlink / shared-inode blocks; the logical TOTAL can overstate freeable space. Confirm with du -sh on the parent — expect roughly one env of reclaim, not the sum of aliases.",
		}
	case "go-build-cache":
		return []string{
			"Go build cache (GOCACHE-style directory).",
			"Safe to clear via go clean -cache; next builds recompile cold.",
		}
	case "npm-cache":
		return []string{
			"npm package cache (.npm / _cacache).",
			"Safe to clear via npm cache clean --force; installs re-download packages.",
		}
	case "homebrew-cache":
		return []string{
			"Homebrew download/bottle cache.",
			"Safe to clear via brew cleanup; bottles re-download when needed.",
		}
	case "generic-qcow2":
		return []string{
			fmt.Sprintf("QEMU/KVM disk image file (%s).", base),
			"Size is the image file itself; reclaim only if the VM is disposable.",
		}
	case "generic-dir":
		return []string{
			fmt.Sprintf("Generic directory (%s); no specialized reclaim kind matched.", base),
			"Review breakdown before reclaiming space.",
		}
	case "generic-file":
		return []string{
			fmt.Sprintf("Generic file (%s); no specialized reclaim kind matched.", base),
		}
	case "xcode":
		return []string{
			"Xcode developer multi-root pack (DerivedData, simulators, DeviceSupport, Archives, DocumentationCache).",
			"Measures existing roots under the scope home-like directory; missing roots are omitted.",
			"DerivedData, docs cache, and device support are usually rebuildable; Archives hold signed builds; simulator wipe removes devices.",
		}
	default:
		return []string{fmt.Sprintf("Path classified as %s.", kind)}
	}
}

func reclaimFor(kind string) []Reclaim {
	switch kind {
	case "android-avd":
		return []Reclaim{
			{
				Title:         "Snapshots only",
				SafeToReclaim: true,
				Detail:        "Usually safe to reclaim snapshots; next boot is colder while apps and user data stay intact. Prefer Android Studio Device Manager or emulator -wipe-data only when you accept data loss.",
			},
			{
				Title:         "Full AVD wipe",
				SafeToReclaim: false,
				Detail:        "Wiping userdata removes installed apps and settings for this AVD. Confirm the emulator is stopped and the AVD is disposable before reclaiming userdata or the whole AVD directory.",
			},
		}
	case "android-sdk":
		return []Reclaim{
			{
				Title:         "Download temp",
				SafeToReclaim: true,
				Detail:        "Usually safe: reclaim .temp / .downloadIntermediates incomplete download leftovers. Does not remove installed packages.",
			},
			{
				Title:         "Unused system-images / sources / skins",
				SafeToReclaim: true,
				Detail:        "Usually safe with caution: unused system-images, optional sources, and skins free the most space. Prefer sdkmanager --uninstall for packages you no longer need; confirm no AVDs still reference those images.",
			},
			{
				Title:         "Keep platform-tools / build-tools / platforms / emulator",
				SafeToReclaim: false,
				Detail:        "Caution: do not purge platform-tools (adb), build-tools, platforms bulk, or the emulator package as usually-safe. Keep versions your projects still compile and run against.",
			},
		}
	case "seatalk-app-support":
		return []Reclaim{
			{
				Title:         "Web caches",
				SafeToReclaim: true,
				Detail:        "Usually safe after quitting SeaTalk: Service Worker, Cache, Code Cache, GPUCache, and Dawn* caches are rebuildable on next launch. Does not touch chat or search databases.",
			},
			{
				Title:         "Local DB backups",
				SafeToReclaim: false,
				Detail:        "Conditional: sqlite-backup and idb-backup may free space if you no longer need those backup copies. Confirm you do not rely on them for recovery before reclaiming.",
			},
			{
				Title:         "Primary chat/search databases",
				SafeToReclaim: false,
				Detail:        "Caution: main_*.sqlite (chat-db) and search_*.sqlite (search-index) hold live message and search data. Never treat them as usually-safe purge targets.",
			},
		}
	case "grok-home":
		return []Reclaim{
			{
				Title:         "Installer downloads and caches",
				SafeToReclaim: true,
				Detail:        "Usually safe / reclaimable: obsolete packages and *.tmp under downloads/, marketplace-cache, logs, models_cache.json, and upload_queue. Next use may re-download installers or rebuild caches.",
			},
			{
				Title:         "Session logs",
				SafeToReclaim: false,
				Detail:        "Caution: sessions/ (session-logs) holds agent transcripts. Space is reclaimable after review, but purging may remove conversation history.",
			},
			{
				Title:         "Config and credentials",
				SafeToReclaim: false,
				Detail:        "Never usually-safe: keep auth.json, config.toml, trusted_folders.toml, and agent_id. These hold credentials and settings — not purge targets.",
			},
		}
	case "codex-home":
		return []Reclaim{
			{
				Title:         "Logs DB, sessions, cache, and tmp",
				SafeToReclaim: true,
				Detail:        "Usually safe / reclaimable with caution: logs_*.sqlite (app-logs-db), sessions/ (session-logs), cache/, and .tmp/. For app-logs-db: quit Codex fully first, then move-aside/backup logs_2.sqlite (+ -wal/-shm) or sqlite3 DELETE FROM logs; VACUUM — diagnostic-only; do not clear state_5.sqlite*, auth, or config. Next use may rebuild caches or recreate a fresh log DB. Review sessions before bulk purge.",
			},
			{
				Title:         "Config and credentials",
				SafeToReclaim: false,
				Detail:        "Never usually-safe: keep auth.json, config.toml, hooks.json, and installation_id. These hold credentials and settings — not purge targets.",
			},
			{
				Title:         "App state databases",
				SafeToReclaim: false,
				Detail:        "Caution: state_*.sqlite (including state_5.sqlite*), goals_*.sqlite, and memories_*.sqlite are app-state-db — not usually-safe bulk purge targets. Do not clear state_5 when reclaiming diagnostic logs.",
			},
		}
	case "iterm2":
		return []Reclaim{
			{
				Title:         "Python env / logs",
				SafeToReclaim: true,
				Detail:        "Usually safe / reclaimable after quitting iTerm2: iterm2env and iterm2env-* (python-env / python-env-alias) and log*.txt. APFS hardlinks may share disk blocks among iterm2env* trees — logical sizes can overstate freeable space; confirm with du -sh on the parent and expect roughly one env of reclaim, not the sum of aliases.",
			},
			{
				Title:         "User config / profiles",
				SafeToReclaim: false,
				Detail:        "Not usually-safe: keep DynamicProfiles, Scripts, and other user-config. Do not bulk-delete profiles as a space reclaim step.",
			},
			{
				Title:         "Helper binary / app DB / state",
				SafeToReclaim: false,
				Detail:        "Caution: iTermServer-* helper binaries, chatdb.sqlite*, SavedState, and version.txt are runtime/app state — not usually-safe purge targets.",
			},
		}
	case "go-build-cache":
		return []Reclaim{
			{
				Title:         "Clear Go build cache",
				SafeToReclaim: true,
				Detail:        "Usually safe: run go clean -cache. Rebuilds will recompile packages cold.",
			},
		}
	case "npm-cache":
		return []Reclaim{
			{
				Title:         "Clear npm cache",
				SafeToReclaim: true,
				Detail:        "Usually safe: run npm cache clean --force. Subsequent installs re-fetch packages.",
			},
		}
	case "homebrew-cache":
		return []Reclaim{
			{
				Title:         "Homebrew cleanup",
				SafeToReclaim: true,
				Detail:        "Usually safe: run brew cleanup. Cached bottles download again when needed.",
			},
		}
	case "generic-qcow2":
		return []Reclaim{
			{
				Title:         "Delete disk image",
				SafeToReclaim: false,
				Detail:        "Only reclaim if the VM/disk is disposable. Back up first if unsure; deleting the qcow2 destroys the virtual disk contents.",
			},
		}
	case "generic-dir":
		return []Reclaim{
			{
				Title:         "Review before reclaim",
				SafeToReclaim: false,
				Detail:        "No specialized reclaim playbook matched. Inspect breakdown and use project-specific tooling before deleting contents.",
			},
		}
	case "generic-file":
		return []Reclaim{
			{
				Title:         "Review before reclaim",
				SafeToReclaim: false,
				Detail:        "Single file with no specialized kind. Confirm it is disposable before deleting.",
			},
		}
	case "xcode":
		return []Reclaim{
			{
				Title:         "DerivedData / docs / device support",
				SafeToReclaim: true,
				Detail:        "Usually safe: DerivedData rebuilds on next build; DocumentationCache re-fetches; iOS DeviceSupport re-downloads when devices reconnect.",
			},
			{
				Title:         "Simulators",
				SafeToReclaim: true,
				Detail:        "Reclaimable but caution: erasing simulators wipes device contents. Prefer xcrun simctl erase / delete after confirming no needed simulator state.",
			},
			{
				Title:         "Archives",
				SafeToReclaim: false,
				Detail:        "Not usually safe: Xcode Archives are signed builds for distribution. Keep unless you intentionally discard old App Store / TestFlight archives.",
			},
		}
	default:
		return []Reclaim{
			{
				Title:         "Manual review",
				SafeToReclaim: false,
				Detail:        "Inspect the path and reclaim only what you recognize as disposable cache or rebuildable artifacts.",
			},
		}
	}
}

// resolveAVDId returns AvdId from config.ini when present, else fallback (folder base without .avd).
func resolveAVDId(avdDir, fallback string) string {
	data, err := os.ReadFile(filepath.Join(avdDir, "config.ini"))
	if err != nil {
		return fallback
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "AvdId=") {
			continue
		}
		id := strings.TrimSpace(strings.TrimPrefix(line, "AvdId="))
		if id != "" {
			return id
		}
	}
	return fallback
}

func howToPurgeFor(kind, targetPath, measureRoot string) []PurgeStep {
	// Prefer measuring root for AVD (the .avd dir) even when user pointed at a file inside.
	root := measureRoot
	if root == "" {
		root = targetPath
	}
	avdName := strings.TrimSuffix(filepath.Base(root), ".avd")
	avdID := resolveAVDId(root, avdName)
	switch kind {
	case "android-avd":
		return []PurgeStep{
			{
				Title: "Cold boot / drop snapshots",
				// No single official CLI for snapshots-only across SDK versions; comment recipe only.
				OfficialCommand: "# After stopping the emulator, remove snapshot data under the AVD (snapshots/).\n# There is no single official one-liner for all SDK versions; wipe-data (below) also invalidates snapshots.",
				Removes:         "Under the AVD: snapshots/ (including default_boot/ram.bin and other quick-boot snapshot data). Does not remove userdata or sdcard.",
				Notes:           "Stop the emulator first. Next launch is a cold boot (slower) but apps/data remain. UI (optional): Device Manager → Cold Boot Now.",
			},
			{
				Title:           "Wipe AVD user data (official)",
				OfficialCommand: fmt.Sprintf("emulator -avd %s -wipe-data\n# (requires emulator on PATH or under $ANDROID_HOME/emulator; AVD must not be running)", avdID),
				Removes:         "userdata-qemu.img, userdata-qemu.img.qcow2, and related user-data overlays; typically resets apps/settings on this AVD. Snapshots may also be invalidated.",
				Notes:           "Keeps the AVD definition (config.ini) and usually the shared SDK system image. Loses installed apps and app data on this device only. UI (optional): Device Manager → Wipe Data.",
			},
			{
				Title:           "Delete entire AVD (official)",
				OfficialCommand: fmt.Sprintf("avdmanager delete avd -n %s\n# (avdmanager under $ANDROID_HOME/cmdline-tools/*/bin; name must match AvdId / list avd)", avdID),
				Removes:         "The whole AVD data directory (config.ini, userdata*, sdcard.img*, snapshots/, cache*, encryptionkey*, locks, etc.) and the registry .ini under ~/.android/avd/.",
				Notes:           "Does not remove the shared Android SDK system-images. Recreate later with avdmanager create avd. UI (optional): Device Manager → Delete.",
			},
		}
	case "android-sdk":
		systemImages := filepath.Join(root, "system-images")
		return []PurgeStep{
			{
				Title: "List installed SDK packages (official)",
				OfficialCommand: "sdkmanager --list_installed\n# or: $ANDROID_HOME/cmdline-tools/latest/bin/sdkmanager --list_installed",
				Removes:         "Nothing; lists installed packages so you can choose unused system-images, platforms, or build-tools versions.",
				Notes:           "sdkmanager lives under cmdline-tools/*/bin when not on PATH. UI (optional): Android Studio → Settings → Languages & Frameworks → Android SDK.",
			},
			{
				Title: "Uninstall unused packages (official)",
				OfficialCommand: "sdkmanager --uninstall \"system-images;android-30;google_apis;x86_64\"\n# replace package path with unused entries from --list_installed (system-images, sources, …)",
				Removes:         "Selected packages only (e.g. unused system-images, optional sources). Does not remove platform-tools (adb) or the emulator package unless you uninstall those packages explicitly.",
				Notes:           "Confirm no AVDs still require the system image. Prefer sdkmanager over bulk deletes. UI (optional): Android Studio SDK Manager → uncheck packages → Apply.",
			},
			{
				Title: "Inspect SDK / system-images sizes",
				OfficialCommand: fmt.Sprintf(
					"disk-usage-analyser scan %s --max-depth 2 --min 1M\ndisk-usage-analyser scan %s --max-depth 2 --min 1M",
					shellQuote(root), shellQuote(systemImages),
				),
				Removes: "Nothing; measurement only so you can confirm which children dominate before uninstalling packages.",
				Notes:   "Do not bulk-delete the SDK tree; use sdkmanager --uninstall for reclaim.",
			},
		}
	case "seatalk-app-support":
		return []PurgeStep{
			{
				Title:           "Quit SeaTalk",
				OfficialCommand: `osascript -e 'quit app "SeaTalk"'`,
				Removes:         "Nothing on disk; prepares safe reclaim by closing open DB/cache handles.",
				Notes:           "Always quit SeaTalk before reclaiming caches or backups. There is no official SeaTalk cache CLI.",
			},
			{
				Title: "Reclaim web caches only",
				OfficialCommand: fmt.Sprintf(
					"osascript -e 'quit app \"SeaTalk\"'\ndisk-usage-analyser scan %s --max-depth 2 --min 1M",
					shellQuote(root),
				),
				Removes: "Service Worker/, Cache/, Code Cache/, GPUCache/, DawnCache/, DawnGraphiteCache/, DawnWebGPUCache/ — rebuilt on next launch. Does not remove main_*.sqlite, search_*.sqlite, or Cookies.",
				Notes:   "No official seatalk cache CLI; inspect with scan then remove only listed cache basenames after quit. UI (optional): none.",
			},
			{
				Title: "Reclaim local DB backups only",
				OfficialCommand: fmt.Sprintf(
					"disk-usage-analyser scan %s --max-depth 2 --min 1M\n# Focus on sqlite-backup/ and idb-backup/ under the SeaTalk Application Support dir",
					shellQuote(root),
				),
				Removes: "sqlite-backup/ and idb-backup/ only — not live main_*.sqlite or search_*.sqlite.",
				Notes:   "Conditional reclaim: keep backups if you may need recovery. Never purge primary chat/search DBs as usually-safe.",
			},
		}
	case "grok-home":
		downloads := filepath.Join(root, "downloads")
		sessions := filepath.Join(root, "sessions")
		return []PurgeStep{
			{
				Title: "Inspect Grok home",
				OfficialCommand: fmt.Sprintf(
					"disk-usage-analyser scan %s --max-depth 2 --min 1M\ndisk-usage-analyser scan %s --max-depth 2 --min 1M\ndisk-usage-analyser scan %s --max-depth 2 --min 1M",
					shellQuote(root), shellQuote(downloads), shellQuote(sessions),
				),
				Removes: "Nothing; measurement only so you can confirm which children dominate size before reclaim.",
				Notes:   "Grok CLI home under .grok. Prefer inspect-first reclaim; never bulk-delete the tree.",
			},
			{
				Title: "Reclaim installer downloads / leftovers",
				OfficialCommand: fmt.Sprintf(
					"disk-usage-analyser scan %s --max-depth 2 --min 1M\n# Keep the current CLI version if still needed; drop obsolete packages and *.tmp under downloads/",
					shellQuote(downloads),
				),
				Removes: "Obsolete packages and *.tmp under downloads/ (installer-cache). Does not remove auth.json or config.toml.",
				Notes:   "Confirm no in-progress install. There is no official grok cache CLI.",
			},
			{
				Title: "Optional old session logs",
				OfficialCommand: fmt.Sprintf(
					"disk-usage-analyser scan %s --max-depth 2 --min 1M\n# Remove only old agent transcripts you no longer need under sessions/",
					shellQuote(sessions),
				),
				Removes: "Optional: old transcripts under sessions/ (session-logs). Caution — may remove conversation history. Does not remove auth/config.",
				Notes:   "Prefer selective reclaim of old sessions only.",
			},
			{
				Title: "Reclaim caches and logs",
				OfficialCommand: fmt.Sprintf(
					"disk-usage-analyser scan %s --max-depth 2 --min 1M\n# Focus on marketplace-cache/, logs/, models_cache.json, upload_queue/",
					shellQuote(root),
				),
				Removes: "marketplace-cache/, logs/, models_cache.json, upload_queue/. Does not remove auth.json, config.toml, or projects/.",
				Notes:   "Never treat auth/config as usually-safe purge targets.",
			},
		}
	case "codex-home":
		sessions := filepath.Join(root, "sessions")
		cache := filepath.Join(root, "cache")
		tmpDir := filepath.Join(root, ".tmp")
		logsDB := filepath.Join(root, "logs_2.sqlite")
		logsWal := filepath.Join(root, "logs_2.sqlite-wal")
		logsShm := filepath.Join(root, "logs_2.sqlite-shm")
		return []PurgeStep{
			{
				Title: "Inspect Codex home",
				OfficialCommand: fmt.Sprintf(
					"disk-usage-analyser scan %s --max-depth 2 --min 1M\ndisk-usage-analyser scan %s --max-depth 2 --min 1M\ndisk-usage-analyser scan %s --max-depth 2 --min 1M",
					shellQuote(root), shellQuote(sessions), shellQuote(cache),
				),
				Removes: "Nothing; measurement only so you can confirm which children dominate size before reclaim (logs_2.sqlite, sessions, cache, .tmp).",
				Notes:   "Codex CLI home under .codex. Prefer inspect-first reclaim; never bulk-delete the tree or auth/config.",
			},
			{
				Title: "Clear diagnostic logs DB",
				OfficialCommand: fmt.Sprintf(
					"# Quit Codex fully first (Codex must not be running — do not truncate while running)\nmv %s %s.bak\nmv %s %s.bak\nmv %s %s.bak\n# Alternative after quit (in-place reclaim):\nsqlite3 %s \"DELETE FROM logs; VACUUM;\"",
					shellQuote(logsDB), shellQuote(logsDB),
					shellQuote(logsWal), shellQuote(logsWal),
					shellQuote(logsShm), shellQuote(logsShm),
					shellQuote(logsDB),
				),
				Removes: "logs_2.sqlite and companions logs_2.sqlite-wal / logs_2.sqlite-shm (app-logs-db) via move-aside backup, or row content via DELETE FROM logs + VACUUM. Does not remove state_5.sqlite*, auth.json, or config.toml.",
				Notes:   "Quit Codex fully first. logs_2.sqlite is diagnostic-only (not state_5 / auth / config / session history). Codex recreates the logs DB on next run; it may regrow under TRACE verbosity. No live truncate while Codex is running. Prefer mv backup (+ -wal/-shm) or sqlite3 DELETE FROM logs; VACUUM — never clear state_5.sqlite*, auth, or config as part of logs reclaim.",
			},
			{
				Title: "Optional sessions / cache / tmp",
				OfficialCommand: fmt.Sprintf(
					"disk-usage-analyser scan %s --max-depth 2 --min 1M\ndisk-usage-analyser scan %s --max-depth 2 --min 1M\ndisk-usage-analyser scan %s --max-depth 2 --min 1M\n# After reviewing session history, reclaim old transcripts under sessions/, cache/, and .tmp/ selectively",
					shellQuote(sessions), shellQuote(cache), shellQuote(tmpDir),
				),
				Removes: "Optional reclaim targets: old transcripts under sessions/ (session-logs), cache/, and .tmp/. Does not remove auth.json, config.toml, or state_5.sqlite*.",
				Notes:   "Caution on sessions — may remove conversation history. Never treat auth/config/state_5 as usually-safe.",
			},
		}
	case "iterm2":
		parent := filepath.Dir(root)
		envPrimary := filepath.Join(root, "iterm2env")
		return []PurgeStep{
			{
				Title: "Inspect iTerm2 Application Support sizes",
				OfficialCommand: fmt.Sprintf(
					"du -sh %s\ndu -sh %s\ndisk-usage-analyser scan %s --max-depth 2 --min 1M\ndu -sh %s/iterm2env*",
					shellQuote(parent), shellQuote(root), shellQuote(root), shellQuote(root),
				),
				Removes: "Nothing; measurement only so you can confirm which children dominate size and whether hardlinked iterm2env* trees share blocks (logical total can overstate freeable space).",
				Notes:   "Prefer du -sh on the parent to gauge actual reclaim. Optional: ls -li under iterm2env* to inspect shared inodes.",
			},
			{
				Title: "Optional logs after quit",
				OfficialCommand: fmt.Sprintf(
					"disk-usage-analyser scan %s --max-depth 1 --min 1M\n# After quitting iTerm2, remove only log*.txt you no longer need",
					shellQuote(root),
				),
				Removes: "Optional: log*.txt app logs after iTerm2 has quit. Does not remove DynamicProfiles, Scripts, or iterm2env*.",
				Notes:   "Quit iTerm2 first so log handles are closed.",
			},
			{
				Title: "Optional all iterm2env* after quit",
				OfficialCommand: fmt.Sprintf(
					"du -sh %s %s/iterm2env*\ndisk-usage-analyser scan %s --max-depth 2 --min 1M\n# After quitting iTerm2, remove all iterm2env* trees together (not a single alias only) so hardlinked blocks actually free",
					shellQuote(envPrimary), shellQuote(root), shellQuote(root),
				),
				Removes: "Optional: all iterm2env and iterm2env-* python env trees together after quit. Because of APFS hardlinks, removing one alias alone may free little. Does not remove DynamicProfiles, Scripts, or other user-config.",
				Notes:   "Never treat DynamicProfiles/Scripts as usually-safe bulk purge. Prefer inspect-first reclaim.",
			},
		}
	case "go-build-cache":
		return []PurgeStep{
			{
				Title:           "Clear Go build cache (official)",
				OfficialCommand: "go clean -cache",
				Removes:         "Contents of GOCACHE (this go-build directory): compiled package objects and build IDs. Does not remove module source under GOMODCACHE (go/pkg/mod).",
				Notes:           "Next builds recompile cold. Optional: go env GOCACHE to confirm the path.",
			},
		}
	case "npm-cache":
		return []PurgeStep{
			{
				Title:           "Clear npm cache (official)",
				OfficialCommand: "npm cache clean --force",
				Removes:         "npm cache contents under the configured cache (typically ~/.npm/_cacache and related index files). Does not remove project node_modules/.",
				Notes:           "npm cache verify can audit first. Re-installs re-download packages.",
			},
		}
	case "homebrew-cache":
		return []PurgeStep{
			{
				Title:           "Homebrew cleanup (official)",
				OfficialCommand: "brew cleanup\n# preview only:\nbrew cleanup -n",
				Removes:         "Old formulae versions and cached bottle/downloads under the Homebrew cache (e.g. Library/Caches/Homebrew). Does not uninstall currently linked formulae.",
				Notes:           "brew cleanup -s also removes the latest download cache more aggressively when you want extra space.",
			},
		}
	case "generic-qcow2":
		return []PurgeStep{
			{
				Title:           "Remove disk image via owning tool",
				OfficialCommand: "# Prefer the VM manager that created the image (Android Studio / AVD Manager, virt-manager, UTM, etc.).\n# There is no single OS-wide official purge for arbitrary .qcow2 files.",
				Removes:         "Only this .qcow2 file if you delete it through that tool; destroying the image removes the virtual disk contents permanently.",
				Notes:           "If this path is under an Android AVD, use explain on the parent *.avd directory for official AVD wipe/delete commands instead.",
			},
		}
	case "generic-dir":
		return []PurgeStep{
			{
				Title:           "No single official purge",
				OfficialCommand: "# Identify the owning app/tool, then use its cleanup or uninstall flow.\n# Inspect first:\ndisk-usage-analyser scan " + shellQuote(root) + " --max-depth 2 --min 1M",
				Removes:         "Depends on the tool. Do not bulk-delete until you know which children are caches vs source data.",
				Notes:           "generic-dir matched; no specialized reclaim profile.",
			},
		}
	case "generic-file":
		return []PurgeStep{
			{
				Title:           "No single official purge",
				OfficialCommand: "# Confirm the file owner/app, then use that app's delete/export UI or package manager uninstall.\nls -lah " + shellQuote(targetPath),
				Removes:         "Only this file if you remove it intentionally; no automated purge recipe for unclassified files.",
			},
		}
	case "xcode":
		derived := filepath.Join(root, filepath.FromSlash("Library/Developer/Xcode/DerivedData"))
		simDevices := filepath.Join(root, filepath.FromSlash("Library/Developer/CoreSimulator/Devices"))
		return []PurgeStep{
			{
				Title: "Inspect Xcode developer roots",
				OfficialCommand: fmt.Sprintf(
					"disk-usage-analyser scan %s --max-depth 2 --min 1M\ndisk-usage-analyser scan %s --max-depth 2 --min 1M",
					shellQuote(derived), shellQuote(simDevices),
				),
				Removes: "Nothing; measurement only so you can confirm which roots dominate size before reclaim.",
				Notes:   "UI (optional): Xcode → Settings → Locations to confirm Derived Data path.",
			},
			{
				Title:           "Erase all simulators (official)",
				OfficialCommand: "xcrun simctl erase all\n# list devices first:\nxcrun simctl list devices",
				Removes:         "Contents of all iOS/watchOS/tvOS simulators under CoreSimulator/Devices (apps, data). Does not remove Xcode.app or DeviceSupport symbols.",
				Notes:           "Stop running simulators first. UI (optional): Xcode → Window → Devices and Simulators → delete/erase devices.",
			},
			{
				Title: "Reclaim DerivedData / DeviceSupport / docs cache",
				OfficialCommand: fmt.Sprintf(
					"disk-usage-analyser scan %s --max-depth 2 --min 1M\n# Prefer Xcode Locations / Storage UI, or remove only confirmed rebuildable dirs after quitting Xcode",
					shellQuote(derived),
				),
				Removes: "DerivedData build products, optional iOS DeviceSupport symbols, DocumentationCache — rebuilt/re-downloaded as needed. Does not remove Archives.",
				Notes:   "UI (optional): Xcode → Settings → Locations → Derived Data arrow; never bulk-delete Archives as usually-safe.",
			},
		}
	default:
		return []PurgeStep{
			{
				Title:           "Manual review",
				OfficialCommand: "disk-usage-analyser scan " + shellQuote(root) + " --max-depth 2 --min 1M",
				Removes:         "Unknown; inspect breakdown before reclaiming anything.",
			},
		}
	}
}

// shellQuote wraps path for display in suggested commands (simple single-quote).
func shellQuote(path string) string {
	if path == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(path, "'", `'"'"'`) + "'"
}

func rawCommandsFor(path, kind string) []RawCommand {
	cmds := []RawCommand{
		{
			Group:   "disk-usage-analyser",
			Command: fmt.Sprintf("disk-usage-analyser scan %s --max-depth 2 --min 1M", path),
		},
		{
			Group:   "system",
			Command: fmt.Sprintf("du -sh %s", path),
		},
	}
	switch kind {
	case "android-avd":
		cmds = append(cmds, RawCommand{
			Group:   "system",
			Command: fmt.Sprintf("du -sh %s/*", path),
		})
		cmds = append(cmds, RawCommand{
			Group:   "android",
			Command: "avdmanager list avd",
		})
	case "android-sdk":
		cmds = append(cmds,
			RawCommand{Group: "android", Command: "sdkmanager --list_installed"},
			RawCommand{Group: "android", Command: "emulator -list-avds"},
			RawCommand{Group: "system", Command: fmt.Sprintf("du -sh %s/*", path)},
		)
	case "go-build-cache":
		cmds = append(cmds,
			RawCommand{Group: "go", Command: "go env GOCACHE"},
			RawCommand{Group: "go", Command: "go clean -cache"},
		)
	case "npm-cache":
		cmds = append(cmds, RawCommand{
			Group:   "npm",
			Command: "npm cache clean --force",
		})
	case "homebrew-cache":
		cmds = append(cmds, RawCommand{
			Group:   "homebrew",
			Command: "brew cleanup -n",
		})
	case "xcode":
		cmds = append(cmds,
			RawCommand{Group: "xcode", Command: "xcrun simctl list devices"},
			RawCommand{Group: "xcode", Command: "xcrun simctl erase all"},
			RawCommand{Group: "system", Command: fmt.Sprintf("du -sh %s/Library/Developer", path)},
		)
	case "grok-home":
		cmds = append(cmds,
			RawCommand{Group: "system", Command: fmt.Sprintf("du -sh %s/*", path)},
			RawCommand{Group: "grok", Command: "grok --version"},
		)
	case "codex-home":
		cmds = append(cmds,
			RawCommand{Group: "system", Command: fmt.Sprintf("du -sh %s/*", path)},
			RawCommand{Group: "system", Command: fmt.Sprintf("ls -lah %s/logs_*.sqlite 2>/dev/null", path)},
		)
	case "iterm2":
		parent := filepath.Dir(path)
		cmds = append(cmds,
			RawCommand{Group: "system", Command: fmt.Sprintf("du -sh %s", parent)},
			RawCommand{Group: "system", Command: fmt.Sprintf("du -sh %s/*", path)},
			RawCommand{Group: "system", Command: fmt.Sprintf("ls -li %s/iterm2env*/ 2>/dev/null | head", path)},
		)
	}
	return cmds
}
