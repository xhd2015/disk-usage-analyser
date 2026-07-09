package explain

import (
	"os"
	"path/filepath"
	"strings"
)

// detectKind classifies path into a reclaim kind. First high-confidence match wins.
// absPath must exist; info is the Stat of absPath.
func detectKind(absPath string, info os.FileInfo) detectResult {
	isFile := !info.IsDir()

	// Prefer android-avd parent context for files (and dirs) inside *.avd.
	if avdRoot := findAndroidAVDRoot(absPath, info); avdRoot != "" {
		return detectResult{
			Kind:         "android-avd",
			Confidence:   "high",
			ContextRoot:  avdRoot,
			TargetIsFile: isFile,
		}
	}

	base := filepath.Base(absPath)

	// go-build-cache: path ends with go-build (GOCACHE-like).
	if base == "go-build" || strings.HasSuffix(absPath, string(filepath.Separator)+"go-build") {
		return detectResult{
			Kind:         "go-build-cache",
			Confidence:   "high",
			ContextRoot:  absPath,
			TargetIsFile: isFile,
		}
	}

	// npm-cache: .npm with _cacache, or path under .npm/_cacache.
	if isNPMCache(absPath, info) {
		root := absPath
		if base == "_cacache" {
			// Prefer parent .npm when present.
			parent := filepath.Dir(absPath)
			if filepath.Base(parent) == ".npm" {
				root = parent
			}
		}
		return detectResult{
			Kind:         "npm-cache",
			Confidence:   "high",
			ContextRoot:  root,
			TargetIsFile: isFile,
		}
	}

	// homebrew-cache: Caches/Homebrew layout.
	if isHomebrewCache(absPath) {
		return detectResult{
			Kind:         "homebrew-cache",
			Confidence:   "high",
			ContextRoot:  absPath,
			TargetIsFile: isFile,
		}
	}

	// seatalk-app-support: …/Application Support/SeaTalk (dir or under tree).
	if seatalkRoot := findSeaTalkAppSupportRoot(absPath, info); seatalkRoot != "" {
		return detectResult{
			Kind:         "seatalk-app-support",
			Confidence:   "high",
			ContextRoot:  seatalkRoot,
			TargetIsFile: isFile,
		}
	}

	// grok-home: …/.grok (dir or under tree) with Grok CLI signatures.
	if grokRoot := findGrokHomeRoot(absPath, info); grokRoot != "" {
		return detectResult{
			Kind:         "grok-home",
			Confidence:   "high",
			ContextRoot:  grokRoot,
			TargetIsFile: isFile,
		}
	}

	// codex-home: …/.codex (dir or under tree) with Codex CLI signatures.
	if codexRoot := findCodexHomeRoot(absPath, info); codexRoot != "" {
		return detectResult{
			Kind:         "codex-home",
			Confidence:   "high",
			ContextRoot:  codexRoot,
			TargetIsFile: isFile,
		}
	}

	// android-sdk: …/Android/sdk or signature layout (dir or under tree).
	if sdkRoot := findAndroidSDKRoot(absPath, info); sdkRoot != "" {
		return detectResult{
			Kind:         "android-sdk",
			Confidence:   "high",
			ContextRoot:  sdkRoot,
			TargetIsFile: isFile,
		}
	}

	// iterm2: …/Application Support/iTerm2 (dir or under tree) with signatures.
	if iterm2Root := findITerm2Root(absPath, info); iterm2Root != "" {
		return detectResult{
			Kind:         "iterm2",
			Confidence:   "high",
			ContextRoot:  iterm2Root,
			TargetIsFile: isFile,
		}
	}

	// generic-qcow2: lone disk image file (not under AVD — already handled above).
	if isFile {
		lower := strings.ToLower(base)
		if strings.HasSuffix(lower, ".qcow2") {
			return detectResult{
				Kind:         "generic-qcow2",
				Confidence:   "medium",
				ContextRoot:  absPath,
				TargetIsFile: true,
			}
		}
	}

	if isFile {
		return detectResult{
			Kind:         "generic-file",
			Confidence:   "low",
			ContextRoot:  absPath,
			TargetIsFile: true,
		}
	}
	return detectResult{
		Kind:         "generic-dir",
		Confidence:   "low",
		ContextRoot:  absPath,
		TargetIsFile: false,
	}
}

// findITerm2Root returns the iTerm2 Application Support directory if path is
// that dir or a file/dir under it (…/Application Support/iTerm2/…).
func findITerm2Root(absPath string, info os.FileInfo) string {
	cur := absPath
	if !info.IsDir() {
		cur = filepath.Dir(absPath)
	}

	for depth := 0; depth < 16; depth++ {
		if isITerm2Root(cur) {
			return cur
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			break
		}
		cur = parent
	}
	return ""
}

// isITerm2Root reports whether dir is an iTerm2 Application Support root by
// path form (basename iTerm2 under Application Support) or signatures.
func isITerm2Root(dir string) bool {
	if dir == "" {
		return false
	}
	// Path form: …/Application Support/iTerm2
	if filepath.Base(dir) == "iTerm2" && filepath.Base(filepath.Dir(dir)) == "Application Support" {
		return true
	}
	sep := string(filepath.Separator)
	if strings.HasSuffix(dir, sep+"Application Support"+sep+"iTerm2") {
		return true
	}
	// Dir named iTerm2 with signatures (iterm2env / version.txt / iTermServer-*).
	if filepath.Base(dir) == "iTerm2" && hasITerm2Signature(dir) {
		return true
	}
	return false
}

// hasITerm2Signature reports iterm2env, version.txt, or iTermServer-* under dir.
func hasITerm2Signature(dir string) bool {
	if dirExists(filepath.Join(dir, "iterm2env")) {
		return true
	}
	if fileExists(filepath.Join(dir, "version.txt")) {
		return true
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "iTermServer-") {
			return true
		}
	}
	return false
}

// findAndroidSDKRoot returns the Android SDK root if path is that dir or a
// file/dir under it. SDK root: path ends with Android/sdk (basename sdk under
// Android), or directory signatures (platform-tools + one of platforms /
// system-images / build-tools / cmdline-tools / emulator).
func findAndroidSDKRoot(absPath string, info os.FileInfo) string {
	cur := absPath
	if !info.IsDir() {
		cur = filepath.Dir(absPath)
	}

	for depth := 0; depth < 16; depth++ {
		if isAndroidSDKRoot(cur) {
			return cur
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			break
		}
		cur = parent
	}
	return ""
}

// isAndroidSDKRoot reports whether dir is an Android SDK root by path form or signatures.
func isAndroidSDKRoot(dir string) bool {
	if dir == "" {
		return false
	}
	// Path form: …/Android/sdk
	if filepath.Base(dir) == "sdk" && filepath.Base(filepath.Dir(dir)) == "Android" {
		return true
	}
	sep := string(filepath.Separator)
	if strings.HasSuffix(dir, sep+"Android"+sep+"sdk") {
		return true
	}
	return hasAndroidSDKSignature(dir)
}

// hasAndroidSDKSignature reports platform-tools plus at least one other SDK marker.
func hasAndroidSDKSignature(dir string) bool {
	if !dirExists(filepath.Join(dir, "platform-tools")) {
		return false
	}
	for _, name := range []string{"platforms", "system-images", "build-tools", "cmdline-tools", "emulator"} {
		if dirExists(filepath.Join(dir, name)) {
			return true
		}
	}
	return false
}

// findGrokHomeRoot returns the .grok directory if path is that dir or a file/dir
// under it and the layout has at least one Grok CLI signature.
func findGrokHomeRoot(absPath string, info os.FileInfo) string {
	cur := absPath
	if !info.IsDir() {
		cur = filepath.Dir(absPath)
	}

	for depth := 0; depth < 12; depth++ {
		if filepath.Base(cur) == ".grok" && hasGrokHomeSignature(cur) {
			return cur
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			break
		}
		cur = parent
	}
	return ""
}

// hasGrokHomeSignature reports whether dir looks like a Grok CLI home
// (config.toml, auth.json, sessions/, or downloads/).
func hasGrokHomeSignature(dir string) bool {
	if fileExists(filepath.Join(dir, "config.toml")) {
		return true
	}
	if fileExists(filepath.Join(dir, "auth.json")) {
		return true
	}
	if dirExists(filepath.Join(dir, "sessions")) {
		return true
	}
	if dirExists(filepath.Join(dir, "downloads")) {
		return true
	}
	return false
}

// findCodexHomeRoot returns the .codex directory if path is that dir or a file/dir
// under it and the layout has at least one Codex CLI signature.
func findCodexHomeRoot(absPath string, info os.FileInfo) string {
	cur := absPath
	if !info.IsDir() {
		cur = filepath.Dir(absPath)
	}

	for depth := 0; depth < 12; depth++ {
		if filepath.Base(cur) == ".codex" && hasCodexHomeSignature(cur) {
			return cur
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			break
		}
		cur = parent
	}
	return ""
}

// hasCodexHomeSignature reports whether dir looks like a Codex CLI home
// (config.toml, auth.json, logs_*.sqlite, sessions/, or history.jsonl).
func hasCodexHomeSignature(dir string) bool {
	if fileExists(filepath.Join(dir, "config.toml")) {
		return true
	}
	if fileExists(filepath.Join(dir, "auth.json")) {
		return true
	}
	if dirExists(filepath.Join(dir, "sessions")) {
		return true
	}
	if fileExists(filepath.Join(dir, "history.jsonl")) {
		return true
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if isCodexLogsDBName(e.Name()) {
			return true
		}
	}
	return false
}

// isCodexLogsDBName reports whether name is logs_*.sqlite or its -wal/-shm siblings.
func isCodexLogsDBName(name string) bool {
	base := name
	for _, suf := range []string{"-wal", "-shm"} {
		if strings.HasSuffix(base, suf) {
			base = strings.TrimSuffix(base, suf)
			break
		}
	}
	return strings.HasPrefix(base, "logs_") && strings.HasSuffix(base, ".sqlite")
}

// findSeaTalkAppSupportRoot returns the SeaTalk Application Support directory if
// path is that dir or a file/dir under it (…/Application Support/SeaTalk/…).
func findSeaTalkAppSupportRoot(absPath string, info os.FileInfo) string {
	cur := absPath
	if !info.IsDir() {
		cur = filepath.Dir(absPath)
	}

	for depth := 0; depth < 12; depth++ {
		base := filepath.Base(cur)
		parent := filepath.Dir(cur)
		if parent == cur {
			break
		}
		// Exact: basename SeaTalk under Application Support.
		if base == "SeaTalk" && filepath.Base(parent) == "Application Support" {
			return cur
		}
		// Path segment form: …/Application Support/SeaTalk/…
		sep := string(filepath.Separator)
		marker := sep + "Application Support" + sep + "SeaTalk"
		if strings.HasSuffix(cur, marker) || strings.Contains(cur, marker+sep) {
			// Walk up until we hit the SeaTalk root under Application Support.
			if base == "SeaTalk" && filepath.Base(parent) == "Application Support" {
				return cur
			}
		}
		cur = parent
	}
	return ""
}

// findAndroidAVDRoot returns the *.avd directory if path is an AVD dir or inside one.
func findAndroidAVDRoot(absPath string, info os.FileInfo) string {
	cur := absPath
	// Include the path itself when it is a directory (e.g. MediumPhone.avd).
	if !info.IsDir() {
		cur = filepath.Dir(absPath)
	}

	for depth := 0; depth < 8; depth++ {
		base := filepath.Base(cur)
		if strings.HasSuffix(strings.ToLower(base), ".avd") {
			return cur
		}
		if looksLikeAVDDir(cur) {
			return cur
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			break
		}
		cur = parent
	}
	return ""
}

func looksLikeAVDDir(dir string) bool {
	// config.ini + (userdata-qemu.img* / sdcard.img / snapshots)
	hasConfig := fileExists(filepath.Join(dir, "config.ini"))
	if !hasConfig {
		return false
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, "userdata-qemu.img") ||
			name == "sdcard.img" ||
			name == "snapshots" {
			return true
		}
	}
	return false
}

func isNPMCache(absPath string, info os.FileInfo) bool {
	base := filepath.Base(absPath)
	if base == ".npm" && info.IsDir() {
		// Signature: _cacache under .npm
		if dirExists(filepath.Join(absPath, "_cacache")) {
			return true
		}
	}
	if base == "_cacache" {
		parent := filepath.Base(filepath.Dir(absPath))
		if parent == ".npm" {
			return true
		}
	}
	// Path contains .npm/_cacache
	sep := string(filepath.Separator)
	if strings.Contains(absPath, sep+".npm"+sep+"_cacache") ||
		strings.HasSuffix(absPath, sep+".npm"+sep+"_cacache") {
		return true
	}
	return false
}

func isHomebrewCache(absPath string) bool {
	base := filepath.Base(absPath)
	if base == "Homebrew" {
		// Prefer Caches/Homebrew
		parent := filepath.Base(filepath.Dir(absPath))
		if parent == "Caches" {
			return true
		}
	}
	sep := string(filepath.Separator)
	if strings.Contains(absPath, sep+"Caches"+sep+"Homebrew") {
		return true
	}
	return false
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

func dirExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}
