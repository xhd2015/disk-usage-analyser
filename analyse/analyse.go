package analyse

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/xhd2015/dot-pkgs/go-pkgs/pathfmt"
)

type DirResult struct {
	Path          string `json:"path"`
	Size          int64  `json:"size"`
	SizeHuman     string `json:"sizeHuman"`
	SymlinkFiles  int    `json:"symlinkFiles"`
	SymlinkDirs   int    `json:"symlinkDirs"`
	HardlinkExtra       int    `json:"hardlinkExtra"`
	HardlinkSize        int64  `json:"hardlinkSize"`
	SharedHardlinkSize  int64  `json:"sharedHardlinkSize"`
	SharedHardlinkHuman string `json:"sharedHardlinkHuman"`
	SharedCloneSize     int64  `json:"sharedCloneSize"`
	SharedCloneHuman    string `json:"sharedCloneHuman"`
	PnpmSharedSize      int64  `json:"pnpmSharedSize"`
	PnpmSharedHuman     string `json:"pnpmSharedHuman"`
	BunSharedSize       int64  `json:"bunSharedSize"`
	BunSharedHuman      string `json:"bunSharedHuman"`
}

type Result struct {
	Root    string      `json:"root"`
	Rows    []DirResult `json:"rows"`
	Summary DirResult   `json:"summary"`
}

type CLIOptions struct {
	Stdout io.Writer
	Stderr io.Writer
}

const help = `Usage: disk-usage-analyser analyse [DIR]

Print per-immediate-child aligned table rows with columns:
  size, symlinks, hardlinks, hardlink_size, shared_hardlink, shared_clone, pnpm_shared, bun_shared, path

Options:
  --header          Accepted for backward compatibility (header is always printed)
  --json            Emit one JSON object instead of the table
  -h, --help        Show help
`

type inodeKey struct {
	dev uint64
	ino uint64
}

type subtreeMetrics struct {
	size          int64
	symlinkFiles  int
	symlinkDirs   int
	hardlinkExtra      int
	hardlinkSize       int64
	sharedHardlinkSize int64
	sharedCloneSize    int64
	pnpmSharedSize     int64
	bunSharedSize      int64
}

func Analyse(root string) (Result, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return Result{}, fmt.Errorf("invalid path %s: %w", root, err)
	}

	info, err := os.Lstat(absRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return Result{}, fmt.Errorf("path does not exist: %s", absRoot)
		}
		return Result{}, fmt.Errorf("cannot access %s: %w", absRoot, err)
	}
	if !info.IsDir() {
		return Result{}, fmt.Errorf("not a directory: %s", absRoot)
	}

	pnpmCtx := newPnpmSharedContext()
	bunCtx := newBunSharedContext()

	summaryMetrics, err := walkSubtree(absRoot, pnpmCtx, bunCtx)
	if err != nil {
		return Result{}, err
	}

	var rows []DirResult
	if err := collectSubdirRows(absRoot, pnpmCtx, bunCtx, &rows); err != nil {
		return Result{}, err
	}

	summary := summaryMetrics.toDirResult(absRoot)
	summary = alignSummarySharedWithRows(summary, rows)

	return Result{
		Root:    absRoot,
		Rows:    rows,
		Summary: summary,
	}, nil
}

// alignSummarySharedWithRows makes depth-1 output consistent: the summary row's
// shared_* columns are the sum of immediate-child rows. Summary size still comes
// from the root walk (clone-deduped unique bytes across children).
func alignSummarySharedWithRows(summary DirResult, rows []DirResult) DirResult {
	if len(rows) == 0 {
		return summary
	}
	var sharedHardlink, sharedClone, pnpmShared, bunShared int64
	for _, row := range rows {
		sharedHardlink += row.SharedHardlinkSize
		sharedClone += row.SharedCloneSize
		pnpmShared += row.PnpmSharedSize
		bunShared += row.BunSharedSize
	}
	summary.SharedHardlinkSize = sharedHardlink
	summary.SharedHardlinkHuman = formatHumanSize(sharedHardlink)
	summary.SharedCloneSize = sharedClone
	summary.SharedCloneHuman = formatHumanSize(sharedClone)
	summary.PnpmSharedSize = pnpmShared
	summary.PnpmSharedHuman = formatHumanSize(pnpmShared)
	summary.BunSharedSize = bunShared
	summary.BunSharedHuman = formatHumanSize(bunShared)
	return summary
}

func RunCLI(args []string, opts CLIOptions) (int, error) {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	var jsonOut bool
	var dir string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			fmt.Fprint(stdout, help)
			return 0, nil
		case "--header":
			// no-op: table header is always printed for backward compatibility
		case "--json":
			jsonOut = true
		default:
			if stringsHasPrefix(arg, "-") {
				return 2, fmt.Errorf("unknown option: %s", arg)
			}
			if dir != "" {
				return 2, fmt.Errorf("unexpected extra argument: %s", arg)
			}
			dir = arg
		}
	}

	if dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return 1, fmt.Errorf("cannot determine current directory: %w", err)
		}
		dir = cwd
	}

	result, err := Analyse(dir)
	if err != nil {
		return 2, err
	}

	if jsonOut {
		payload, err := json.Marshal(result)
		if err != nil {
			return 1, fmt.Errorf("encode json: %w", err)
		}
		if _, err := fmt.Fprintln(stdout, string(payload)); err != nil {
			return 1, err
		}
		return 0, nil
	}

	tableRows := make([]DirResult, 0, len(result.Rows)+1)
	tableRows = append(tableRows, result.Rows...)
	tableRows = append(tableRows, result.Summary)
	fmt.Fprintln(stdout, formatAnalyseTable(tableRows))
	return 0, nil
}

func collectSubdirRows(dir string, pnpmCtx *pnpmSharedContext, bunCtx *bunSharedContext, rows *[]DirResult) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	dirSymlinkTargets := make(map[string]struct{})
	for _, ent := range entries {
		if ent.Type()&fs.ModeSymlink == 0 {
			continue
		}
		linkPath := filepath.Join(dir, ent.Name())
		if target, ok := symlinkTargetDir(linkPath); ok {
			dirSymlinkTargets[target] = struct{}{}
		}
	}

	for _, ent := range entries {
		isSymlink := ent.Type()&fs.ModeSymlink != 0
		if !isSymlink && !ent.IsDir() {
			continue
		}

		subPath := filepath.Join(dir, ent.Name())
		metrics, err := walkSubtree(subPath, pnpmCtx, bunCtx)
		if err != nil {
			continue
		}
		if !isSymlink {
			if _, covered := dirSymlinkTargets[subPath]; covered {
				if metrics.symlinkFiles == 0 && metrics.symlinkDirs == 0 && metrics.hardlinkExtra == 0 {
					continue
				}
			}
		}
		*rows = append(*rows, metrics.toDirResult(subPath))
	}
	return nil
}

func walkSubtree(rootPath string, pnpmCtx *pnpmSharedContext, bunCtx *bunSharedContext) (subtreeMetrics, error) {
	var metrics subtreeMetrics

	rootInfo, err := os.Lstat(rootPath)
	if err != nil {
		return metrics, err
	}
	if rootInfo.Mode()&os.ModeSymlink != 0 {
		if isDirSymlink(rootPath) {
			metrics.symlinkDirs++
			if target, ok := symlinkTargetDir(rootPath); ok {
				files, dirs := countSymlinksOnly(target)
				metrics.symlinkFiles += files
				metrics.symlinkDirs += dirs
			}
		} else {
			metrics.symlinkFiles++
		}
		return metrics, nil
	}

	if filepath.Base(rootPath) == "node_modules" {
		metrics.pnpmSharedSize += pnpmCtx.pnpmSharedForNodeModules(rootPath)
		metrics.bunSharedSize += bunCtx.bunSharedForNodeModules(rootPath)
	} else if filepath.Base(filepath.Dir(rootPath)) == "node_modules" {
		metrics.pnpmSharedSize += pnpmCtx.pnpmSharedForNodeModules(rootPath)
		metrics.bunSharedSize += bunCtx.bunSharedForNodeModules(rootPath)
	}

	seenInodes := make(map[inodeKey]struct{})
	countedHardlinks := make(map[inodeKey]struct{})
	cloneTracker := newCloneGroupTracker()

	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if path == rootPath {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		mode := info.Mode()
		if mode&os.ModeSymlink != 0 {
			if isDirSymlink(path) {
				metrics.symlinkDirs++
				if target, ok := symlinkTargetDir(path); ok {
					files, dirs := countSymlinksOnly(target)
					metrics.symlinkFiles += files
					metrics.symlinkDirs += dirs
				}
			} else {
				metrics.symlinkFiles++
			}
			return nil
		}

		if d.IsDir() {
			if d.Name() == "node_modules" && shouldScanOuterNodeModules(rootPath, path) {
				metrics.pnpmSharedSize += pnpmCtx.pnpmSharedForNodeModules(path)
				metrics.bunSharedSize += bunCtx.bunSharedForNodeModules(path)
			}
			return nil
		}

		if !mode.IsRegular() {
			return nil
		}

		key, ok := inodeKeyFrom(info)
		if !ok {
			metrics.size += info.Size()
			return nil
		}

		if _, seen := seenInodes[key]; !seen {
			seenInodes[key] = struct{}{}
			metrics.size += cloneTracker.CountSize(path, info.Size())
		}

		nlink, ok := nlinkFrom(info)
		if ok && nlink > 1 {
			if _, counted := countedHardlinks[key]; !counted {
				countedHardlinks[key] = struct{}{}
				metrics.hardlinkExtra += int(nlink) - 1
				metrics.hardlinkSize += info.Size()
				metrics.sharedHardlinkSize += info.Size() * int64(nlink)
			}
			return nil
		}

		cloneTracker.Add(path, key, info.Size())
		return nil
	})
	if err == nil {
		metrics.sharedCloneSize = cloneTracker.TotalSharedCloneSize()
	}
	return metrics, err
}

// shouldScanOuterNodeModules returns true when nodeModulesPath is the
// outermost node_modules relative to walkSubtree rootPath. Nested pnpm layouts
// (e.g. node_modules/.pnpm/pkg@1/node_modules) must not trigger a separate scan;
// when rootPath itself is node_modules the upfront scan already covers the tree.
func shouldScanOuterNodeModules(rootPath, nodeModulesPath string) bool {
	if filepath.Base(rootPath) == "node_modules" {
		return false
	}
	if filepath.Base(filepath.Dir(rootPath)) == "node_modules" {
		return false
	}
	rel, err := filepath.Rel(rootPath, nodeModulesPath)
	if err != nil {
		return false
	}
	rel = filepath.ToSlash(rel)
	count := 0
	for _, part := range strings.Split(rel, "/") {
		if part == "node_modules" {
			count++
		}
	}
	return count == 1
}

func resolveSymlinkTarget(linkPath string) (string, error) {
	target, err := os.Readlink(linkPath)
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(linkPath), target)
	}
	return target, nil
}

func isDirSymlink(linkPath string) bool {
	target, err := resolveSymlinkTarget(linkPath)
	if err != nil {
		return false
	}
	info, err := os.Stat(target)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func symlinkTargetDir(linkPath string) (string, bool) {
	target, err := resolveSymlinkTarget(linkPath)
	if err != nil {
		return "", false
	}
	info, err := os.Stat(target)
	if err != nil || !info.IsDir() {
		return "", false
	}
	return target, true
}

func countSymlinksOnly(root string) (files, dirs int) {
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || path == root {
			return nil
		}
		if d.Type()&fs.ModeSymlink != 0 {
			if isDirSymlink(path) {
				dirs++
			} else {
				files++
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		return nil
	})
	return files, dirs
}

func inodeKeyFrom(info os.FileInfo) (inodeKey, bool) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return inodeKey{}, false
	}
	return inodeKey{dev: uint64(stat.Dev), ino: stat.Ino}, true
}

func nlinkFrom(info os.FileInfo) (uint64, bool) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, false
	}
	return uint64(stat.Nlink), true
}

func (m subtreeMetrics) toDirResult(path string) DirResult {
	return DirResult{
		Path:          path,
		Size:          m.size,
		SizeHuman:     formatHumanSize(m.size),
		SymlinkFiles:  m.symlinkFiles,
		SymlinkDirs:   m.symlinkDirs,
		HardlinkExtra:       m.hardlinkExtra,
		HardlinkSize:        m.hardlinkSize,
		SharedHardlinkSize:  m.sharedHardlinkSize,
		SharedHardlinkHuman: formatHumanSize(m.sharedHardlinkSize),
		SharedCloneSize:     m.sharedCloneSize,
		SharedCloneHuman:    formatHumanSize(m.sharedCloneSize),
		PnpmSharedSize:      m.pnpmSharedSize,
		PnpmSharedHuman:     formatHumanSize(m.pnpmSharedSize),
		BunSharedSize:       m.bunSharedSize,
		BunSharedHuman:      formatHumanSize(m.bunSharedSize),
	}
}

const tableColGap = "  "

var analyseTableHeaders = []string{
	"size", "symlinks", "hardlinks", "hardlink_size", "shared_hardlink", "shared_clone", "pnpm_shared", "bun_shared", "path",
}

func formatAnalyseTable(rows []DirResult) string {
	cellRows := make([][]string, 0, len(rows)+1)
	cellRows = append(cellRows, analyseTableHeaders)
	for _, row := range rows {
		cellRows = append(cellRows, analyseTableCells(row))
	}

	widths := make([]int, len(analyseTableHeaders))
	for _, cells := range cellRows {
		for i, cell := range cells {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	lines := make([]string, len(cellRows))
	for rowIdx, cells := range cellRows {
		parts := make([]string, len(cells))
		for i, cell := range cells {
			if i == len(cells)-1 {
				parts[i] = cell
			} else {
				parts[i] = padLeft(cell, widths[i])
			}
		}
		lines[rowIdx] = strings.Join(parts, tableColGap)
	}
	return strings.Join(lines, "\n")
}

func analyseTableCells(row DirResult) []string {
	return []string{
		row.SizeHuman,
		symlinkLabel(row.SymlinkFiles, row.SymlinkDirs),
		strconv.Itoa(row.HardlinkExtra),
		formatHumanSize(row.HardlinkSize),
		formatHumanSize(row.SharedHardlinkSize),
		formatHumanSize(row.SharedCloneSize),
		formatHumanSize(row.PnpmSharedSize),
		formatHumanSize(row.BunSharedSize),
		pathfmt.Short(row.Path),
	}
}

func padLeft(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-len(s)) + s
}

func symlinkLabel(files, dirs int) string {
	return fmt.Sprintf("%df+%dd", files, dirs)
}

func formatHumanSize(size int64) string {
	if size <= 1024 {
		return fmt.Sprintf("%d B", size)
	}
	units := []string{"K", "M", "G", "T", "P"}
	value := float64(size)
	for _, unit := range units {
		value /= 1024
		if value < 1024 {
			if value == float64(int64(value)) {
				return fmt.Sprintf("%.0f%s", value, unit)
			}
			return fmt.Sprintf("%.1f%s", value, unit)
		}
	}
	return fmt.Sprintf("%.1fE", value/1024)
}

func stringsHasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}