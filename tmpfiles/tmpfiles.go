package tmpfiles

import (
	"context"
	"debug/buildinfo"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/xhd2015/dot-pkgs/go-pkgs/file/detect"
	"github.com/xhd2015/dot-pkgs/go-pkgs/file/remotefs"
	"github.com/xhd2015/dot-pkgs/go-pkgs/git/scan_repo"
)

type BinaryHit struct {
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	SizeHuman string `json:"sizeHuman"`
	Kind      string `json:"kind"`
	TypeDesc  string `json:"typeDesc"`
	RepoPath  string `json:"repoPath"`
	RepoName  string `json:"repoName"`
}

type NamedHit struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	SizeHuman string `json:"sizeHuman"`
	RepoPath  string `json:"repoPath"`
	RepoName  string `json:"repoName"`
}

type ScanResult struct {
	Roots      []string
	Repos      int
	Binaries   []BinaryHit
	NamedHits  []NamedHit
	TotalSize  int64
	TotalHuman string
}

type CLIOptions struct {
	Stdout  io.Writer
	Stderr  io.Writer
	HomeDir string
}

type scanOptions struct {
	roots    []string
	maxDepth int
	json     bool
	verbose  bool
	names    []string
}

const help = `Usage: disk-usage-analyser tmp-files scan [OPTIONS]

Options:
  --go-binaries     Scan for Go, Mach-O, and ELF binaries
  --name NAME       Scan for entries whose basename matches NAME (repeatable)
  --root PATH       Root directory to scan (repeatable; default ~)
  --max-depth N     Repo discovery depth (0 = unlimited)
  --json            Emit one JSON object per hit
  -v, --verbose     Warn on stderr when skipping unreadable or remote-backed paths
  -h, --help        Show help
`

var ignoredDirBasenames = map[string]struct{}{
	".git":         {},
	"node_modules": {},
	"vendor":       {},
	".venv":        {},
	"__pycache__":  {},
	"dist":         {},
	"build":        {},
	"target":       {},
}

func RunCLI(ctx context.Context, args []string, opts CLIOptions) (*ScanResult, int, error) {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	if opts.Stderr == nil {
		opts.Stderr = os.Stderr
	}

	if len(args) == 0 {
		fmt.Fprint(stdout, help)
		return nil, 0, nil
	}
	if args[0] != "scan" {
		return nil, 2, fmt.Errorf("unknown tmp-files command: %s", args[0])
	}

	scanOpts, showHelp, err := parseScanArgs(args[1:], opts.HomeDir)
	if err != nil {
		return nil, 2, err
	}
	if showHelp {
		fmt.Fprint(stdout, help)
		return nil, 0, nil
	}

	result, err := runScan(ctx, scanOpts, opts.HomeDir, stdout, opts.Stderr)
	if err != nil {
		return result, 1, err
	}
	return result, 0, nil
}

func parseScanArgs(args []string, homeDir string) (scanOptions, bool, error) {
	var opts scanOptions
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if name, ok := extractNameValue(arg); ok {
			opts.names = append(opts.names, name)
			continue
		}
		switch arg {
		case "-h", "--help":
			return opts, true, nil
		case "--go-binaries":
		case "--name":
			if i+1 >= len(args) {
				return opts, false, fmt.Errorf("--name requires a value")
			}
			i++
			opts.names = append(opts.names, args[i])
		case "--json":
			opts.json = true
		case "-v", "--verbose":
			opts.verbose = true
		case "--root":
			if i+1 >= len(args) {
				return opts, false, fmt.Errorf("--root requires a path")
			}
			i++
			root, err := expandPath(args[i], homeDir)
			if err != nil {
				return opts, false, err
			}
			opts.roots = append(opts.roots, root)
		case "--max-depth":
			if i+1 >= len(args) {
				return opts, false, fmt.Errorf("--max-depth requires a value")
			}
			i++
			n, err := strconv.Atoi(args[i])
			if err != nil || n < 0 {
				return opts, false, fmt.Errorf("--max-depth must be a non-negative integer")
			}
			opts.maxDepth = n
		default:
			return opts, false, fmt.Errorf("unknown scan option: %s", arg)
		}
	}
	if len(opts.roots) == 0 {
		if homeDir == "" {
			var err error
			homeDir, err = os.UserHomeDir()
			if err != nil {
				return opts, false, err
			}
		}
		opts.roots = []string{homeDir}
	}
	return opts, false, nil
}

// ScanOptions configures a binary scan under git repositories.
type ScanOptions struct {
	Roots         []string
	MaxDepth      int
	Verbose       bool
	Stderr        io.Writer
	Names         []string
	OnRepo         func(repo scan_repo.Repo) error
	OnNamedPreview func(hit NamedHit, repo scan_repo.Repo) error
	OnNamedHit     func(hit NamedHit, repo scan_repo.Repo) error
}

// BinaryScanSummary holds aggregate binary scan results.
type BinaryScanSummary struct {
	Repos       int
	Binaries    int
	NamedHits   int
	TotalSize   int64
	TotalHuman  string
}

// BinaryHitCallback is invoked for each discovered binary.
type BinaryHitCallback func(hit BinaryHit, repo scan_repo.Repo) error

// Scan discovers git repos and walks each repo for binary and/or named-entry matches.
func Scan(ctx context.Context, opts ScanOptions, homeDir string, onHit BinaryHitCallback) (BinaryScanSummary, error) {
	namesSet := makeNameSet(opts.Names)
	summary := BinaryScanSummary{}
	seenRepos := make(map[string]struct{})
	var seenMu sync.Mutex
	var summaryMu sync.Mutex
	var sizingWg sync.WaitGroup
	trackAsyncSizing := opts.OnNamedPreview != nil

	processRepo := func(repo scan_repo.Repo) error {
		repoKey := strings.ToLower(filepath.Clean(repo.Path))
		seenMu.Lock()
		if _, ok := seenRepos[repoKey]; ok {
			seenMu.Unlock()
			return nil
		}
		seenRepos[repoKey] = struct{}{}
		seenMu.Unlock()

		summaryMu.Lock()
		summary.Repos++
		summaryMu.Unlock()

		var sizingWG *sync.WaitGroup
		if trackAsyncSizing {
			sizingWG = &sizingWg
		}
		if err := scanRepoFiles(ctx, repo, homeDir, namesSet, onHit != nil, opts.Verbose, opts.Stderr, sizingWG, func(hit NamedHit) error {
			if opts.OnNamedPreview == nil {
				return nil
			}
			return opts.OnNamedPreview(hit, repo)
		}, func(hit BinaryHit) error {
			summaryMu.Lock()
			summary.Binaries++
			summary.TotalSize += hit.Size
			summaryMu.Unlock()
			if onHit != nil {
				return onHit(hit, repo)
			}
			return nil
		}, func(hit NamedHit) error {
			summaryMu.Lock()
			summary.NamedHits++
			summary.TotalSize += hit.Size
			summaryMu.Unlock()
			if opts.OnNamedHit != nil {
				return opts.OnNamedHit(hit, repo)
			}
			return nil
		}); err != nil {
			return err
		}
		if opts.OnRepo != nil {
			return opts.OnRepo(repo)
		}
		return nil
	}

	scanOneRoot := func(roots []string, ignoreDirs []string) error {
		_, err := scan_repo.Scan(ctx, scan_repo.Options{
			Roots:      roots,
			MaxDepth:   opts.MaxDepth,
			IgnoreDirs: ignoreDirs,
			Verbose:    opts.Verbose,
			Stderr:     opts.Stderr,
			OnRepo:     processRepo,
		})
		return err
	}

	prefixIgnoreDirs := func(current string, prior []string) []string {
		currentClean := filepath.Clean(current)
		ignoreDirs := make([]string, 0, len(prior))
		for _, root := range prior {
			clean := filepath.Clean(root)
			rel, relErr := filepath.Rel(currentClean, clean)
			if relErr == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				ignoreDirs = append(ignoreDirs, clean)
			}
		}
		return ignoreDirs
	}

	var err error
	if len(opts.Roots) > 1 {
		homeRoot := opts.Roots[len(opts.Roots)-1]
		prefixRoots := opts.Roots[:len(opts.Roots)-1]
		scanned := make([]string, 0, len(prefixRoots))
		for _, root := range prefixRoots {
			if err = scanOneRoot([]string{root}, prefixIgnoreDirs(root, scanned)); err != nil {
				return summary, err
			}
			scanned = append(scanned, root)
		}
		err = scanOneRoot([]string{homeRoot}, prefixIgnoreDirs(homeRoot, prefixRoots))
	} else {
		err = scanOneRoot(opts.Roots, nil)
	}
	if err != nil {
		return summary, err
	}

	if trackAsyncSizing {
		done := make(chan struct{})
		go func() {
			sizingWg.Wait()
			close(done)
		}()
		select {
		case <-done:
		case <-ctx.Done():
			return summary, ctx.Err()
		}
	}

	summary.TotalHuman = FormatHumanSize(summary.TotalSize)
	return summary, nil
}

func runScan(ctx context.Context, opts scanOptions, homeDir string, stdout, stderr io.Writer) (*ScanResult, error) {
	scanOpts := ScanOptions{
		Roots:    opts.roots,
		MaxDepth: opts.maxDepth,
		Verbose:  opts.verbose,
		Stderr:   stderr,
		Names:    opts.names,
	}
	var result ScanResult
	result.Roots = opts.roots
	scanOpts.OnNamedHit = func(hit NamedHit, _ scan_repo.Repo) error {
		result.NamedHits = append(result.NamedHits, hit)
		result.TotalSize += hit.Size
		if opts.json {
			if err := writeJSONNamedHit(stdout, hit); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(stdout, "%s  name:%-5s  %s  (repo: %s)\n", hit.SizeHuman, hit.Name, hit.Path, hit.RepoPath); err != nil {
				return err
			}
		}
		return flush(stdout)
	}
	summary, err := Scan(ctx, scanOpts, homeDir, func(hit BinaryHit, _ scan_repo.Repo) error {
		result.Binaries = append(result.Binaries, hit)
		result.TotalSize += hit.Size
		if opts.json {
			if err := writeJSONHit(stdout, hit); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(stdout, "%s  %-5s  %s  (repo: %s)\n", hit.SizeHuman, hit.Kind, hit.Path, hit.RepoPath); err != nil {
				return err
			}
		}
		return flush(stdout)
	})
	if err != nil {
		return &result, err
	}
	result.Repos = summary.Repos
	result.TotalHuman = summary.TotalHuman
	if len(opts.names) > 0 {
		if _, err := fmt.Fprintf(stdout, "Found %d binaries, %d named entries, total %s\n", len(result.Binaries), len(result.NamedHits), result.TotalHuman); err != nil {
			return &result, err
		}
	} else {
		if _, err := fmt.Fprintf(stdout, "Found %d binaries, total %s\n", len(result.Binaries), result.TotalHuman); err != nil {
			return &result, err
		}
	}
	return &result, flush(stdout)
}

func emitNamedDir(path, name string, repo scan_repo.Repo, homeDir string, names map[string]bool, sizingWG *sync.WaitGroup, onNamedPreview func(NamedHit) error, onNamedHit func(NamedHit) error) error {
	if onNamedPreview != nil {
		preview := NamedHit{
			Path:      DisplayPath(path, homeDir),
			Name:      name,
			Size:      0,
			SizeHuman: "0 B",
			RepoPath:  DisplayPath(repo.Path, homeDir),
			RepoName:  repo.Name,
		}
		if err := onNamedPreview(preview); err != nil {
			return err
		}
		if sizingWG != nil {
			sizingWG.Add(1)
			go func() {
				defer sizingWG.Done()
				_ = emitNamedDirSized(path, name, repo, homeDir, names, onNamedHit)
			}()
		} else {
			go emitNamedDirSized(path, name, repo, homeDir, names, onNamedHit)
		}
		return nil
	}
	return emitNamedDirSized(path, name, repo, homeDir, names, onNamedHit)
}

func emitNamedDirSized(path, name string, repo scan_repo.Repo, homeDir string, names map[string]bool, onNamedHit func(NamedHit) error) error {
	size, innerHits, err := computeNamedDirSize(path, names)
	if err != nil {
		return err
	}
	for _, ih := range innerHits {
		ih.Path = DisplayPath(ih.Path, homeDir)
		ih.RepoPath = DisplayPath(repo.Path, homeDir)
		ih.RepoName = repo.Name
		if err := onNamedHit(ih); err != nil {
			return err
		}
	}
	hit := NamedHit{
		Path:      DisplayPath(path, homeDir),
		Name:      name,
		Size:      size,
		SizeHuman: FormatHumanSize(size),
		RepoPath:  DisplayPath(repo.Path, homeDir),
		RepoName:  repo.Name,
	}
	return onNamedHit(hit)
}

func scanRepoFiles(ctx context.Context, repo scan_repo.Repo, homeDir string, names map[string]bool, classifyBinaries bool, verbose bool, stderr io.Writer, sizingWG *sync.WaitGroup, onNamedPreview func(NamedHit) error, onHit func(BinaryHit) error, onNamedHit func(NamedHit) error) error {
	handledNamedDirs := make(map[string]struct{})
	if len(names) > 0 {
		for name := range names {
			dirPath := filepath.Join(repo.Path, name)
			info, err := os.Stat(dirPath)
			if err != nil || !info.IsDir() {
				continue
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			if err := emitNamedDir(dirPath, name, repo, homeDir, names, sizingWG, onNamedPreview, onNamedHit); err != nil {
				return err
			}
			handledNamedDirs[dirPath] = struct{}{}
		}
	}

	return filepath.WalkDir(repo.Path, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if os.IsPermission(walkErr) {
				maybeWarnSkip(verbose, stderr, path, walkErr)
				return filepath.SkipDir
			}
			return walkErr
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if d.IsDir() {
			if path != repo.Path {
				if _, done := handledNamedDirs[path]; done {
					return filepath.SkipDir
				}
				if len(names) > 0 && names[d.Name()] {
					if err := emitNamedDir(path, d.Name(), repo, homeDir, names, sizingWG, onNamedPreview, onNamedHit); err != nil {
						return err
					}
					return filepath.SkipDir
				}
				if _, ok := ignoredDirBasenames[d.Name()]; ok {
					return filepath.SkipDir
				}
				if remote, err := remotefs.IsRemoteBackedPath(path); err != nil {
					maybeWarnSkip(verbose, stderr, path, err)
					return filepath.SkipDir
				} else if remote {
					maybeWarnSkipRemote(verbose, stderr, path)
					return filepath.SkipDir
				}
			}
			return nil
		}

		if !classifyBinaries && (len(names) == 0 || !names[d.Name()]) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			if os.IsPermission(err) {
				return nil
			}
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		if len(names) > 0 && names[d.Name()] {
			hit := NamedHit{
				Path:      DisplayPath(path, homeDir),
				Name:      d.Name(),
				Size:      info.Size(),
				SizeHuman: FormatHumanSize(info.Size()),
				RepoPath:  DisplayPath(repo.Path, homeDir),
				RepoName:  repo.Name,
			}
			return onNamedHit(hit)
		}

		if !classifyBinaries {
			return nil
		}
		hit, ok := ClassifyFile(path, info.Size(), repo, homeDir)
		if !ok {
			return nil
		}
		return onHit(hit)
	})
}

// ClassifyFile returns a BinaryHit when path is a Go/Mach-O/ELF binary.
func ClassifyFile(path string, size int64, repo scan_repo.Repo, homeDir string) (BinaryHit, bool) {
	typeDesc, isBinary, detectErr := detect.DetectFileType(path)
	if detectErr != nil {
		return BinaryHit{}, false
	}

	kind := ""
	if _, err := buildinfo.ReadFile(path); err == nil {
		kind = "go"
	} else if isBinary {
		switch {
		case strings.HasPrefix(typeDesc, "Mach-O"):
			kind = "macho"
		case strings.HasPrefix(typeDesc, "ELF"):
			kind = "elf"
		}
	}
	if kind == "" {
		return BinaryHit{}, false
	}
	if typeDesc == "" {
		typeDesc = kind
	}
	return BinaryHit{
		Path:      DisplayPath(path, homeDir),
		Size:      size,
		SizeHuman: FormatHumanSize(size),
		Kind:      kind,
		TypeDesc:  typeDesc,
		RepoPath:  DisplayPath(repo.Path, homeDir),
		RepoName:  repo.Name,
	}, true
}

func writeJSONHit(w io.Writer, hit BinaryHit) error {
	data, err := json.Marshal(hit)
	if err != nil {
		return err
	}
	if _, err := w.Write(append(data, '\n')); err != nil {
		return err
	}
	return nil
}

func maybeWarnSkip(verbose bool, stderr io.Writer, path string, err error) {
	if !verbose || stderr == nil {
		return
	}
	fmt.Fprintf(stderr, "\nwarning: skipping\n%s: %v", path, err)
}

func maybeWarnSkipRemote(verbose bool, stderr io.Writer, path string) {
	if !verbose || stderr == nil {
		return
	}
	fmt.Fprintf(stderr, "\nwarning: skipping remote-backed filesystem\n%s", path)
}

func flush(w io.Writer) error {
	if f, ok := w.(interface{ Flush() error }); ok {
		return f.Flush()
	}
	return nil
}

func expandPath(path string, homeDir string) (string, error) {
	if path == "~" || strings.HasPrefix(path, "~/") {
		if homeDir == "" {
			var err error
			homeDir, err = os.UserHomeDir()
			if err != nil {
				return "", err
			}
		}
		if path == "~" {
			return homeDir, nil
		}
		return filepath.Join(homeDir, path[2:]), nil
	}
	return path, nil
}

// PrioritizedScanRoots returns home subdirectories likely to contain git repos first,
// then the full home directory, so named scans surface early hits sooner.
func PrioritizedScanRoots(homeDir string) []string {
	home := filepath.Clean(homeDir)
	candidates := []string{
		filepath.Join(home, "Projects", "gopath"),
		filepath.Join(home, "Projects"),
		filepath.Join(home, "go", "src"),
		filepath.Join(home, "work"),
		filepath.Join(home, ".wrk", "worktrees"),
		filepath.Join(home, "code"),
		filepath.Join(home, "dev"),
		filepath.Join(home, "src"),
		filepath.Join(home, "workspace"),
		filepath.Join(home, "repos"),
	}
	roots := make([]string, 0, len(candidates)+1)
	seen := make(map[string]struct{}, len(candidates)+1)
	for _, candidate := range candidates {
		abs, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		abs = filepath.Clean(abs)
		key := strings.ToLower(abs)
		if _, ok := seen[key]; ok {
			continue
		}
		info, err := os.Stat(abs)
		if err != nil || !info.IsDir() {
			continue
		}
		seen[key] = struct{}{}
		roots = append(roots, abs)
	}
	homeAbs, err := filepath.Abs(home)
	if err != nil {
		if len(roots) == 0 {
			return []string{home}
		}
		return roots
	}
	homeAbs = filepath.Clean(homeAbs)
	if _, ok := seen[strings.ToLower(homeAbs)]; !ok {
		roots = append(roots, homeAbs)
	}
	if len(roots) == 0 {
		return []string{homeAbs}
	}
	return roots
}

// DisplayPath converts an absolute path to a ~/ prefixed slash path for display.
func DisplayPath(path string, homeDir string) string {
	clean := filepath.Clean(path)
	if homeDir == "" {
		return clean
	}
	home := filepath.Clean(homeDir)
	if clean == home {
		return "~"
	}
	rel, err := filepath.Rel(home, clean)
	if err == nil && rel != "." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != ".." {
		return filepath.ToSlash(filepath.Join("~", rel))
	}
	return filepath.ToSlash(clean)
}

func FormatHumanSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	value := float64(size)
	for _, unit := range units {
		value /= 1024
		if value < 1024 {
			return fmt.Sprintf("%.1f %s", value, unit)
		}
	}
	return fmt.Sprintf("%.1f EB", value/1024)
}

func extractNameValue(arg string) (string, bool) {
	if strings.HasPrefix(arg, "--name=") {
		val := arg[len("--name="):]
		if val == "" {
			return "", false
		}
		return val, true
	}
	return "", false
}

func makeNameSet(names []string) map[string]bool {
	if len(names) == 0 {
		return nil
	}
	set := make(map[string]bool, len(names))
	for _, n := range names {
		set[n] = true
	}
	return set
}

func computeNamedDirSize(path string, names map[string]bool) (size int64, innerHits []NamedHit, err error) {
	var total int64
	var hits []NamedHit
	err = filepath.WalkDir(path, func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if os.IsPermission(walkErr) {
				return filepath.SkipDir
			}
			return walkErr
		}
		if p == path {
			return nil
		}
		if d.IsDir() && names != nil && names[d.Name()] {
			innerSize, innerInner, innerErr := computeNamedDirSize(p, names)
			if innerErr != nil {
				return innerErr
			}
			hits = append(hits, innerInner...)
			hits = append(hits, NamedHit{Path: p, Name: d.Name(), Size: innerSize, SizeHuman: FormatHumanSize(innerSize)})
			return filepath.SkipDir
		}
		if !d.Type().IsRegular() {
			return nil
		}
		info, infoErr := d.Info()
		if infoErr != nil {
			return nil
		}
		total += info.Size()
		return nil
	})
	return total, hits, err
}

func computeDirSize(path string, names map[string]bool) (int64, error) {
	size, _, err := computeNamedDirSize(path, names)
	return size, err
}

// DirSize returns the recursive byte size of path. When names is non-nil, nested
// directories whose basename is in names are sized separately and excluded from
// the parent total (same rules as node_modules scan sizing).
func DirSize(path string, names map[string]bool) (int64, error) {
	return computeDirSize(path, names)
}

func writeJSONNamedHit(w io.Writer, hit NamedHit) error {
	type namedHitJSON struct {
		Type      string `json:"type"`
		Path      string `json:"path"`
		Name      string `json:"name"`
		Size      int64  `json:"size"`
		SizeHuman string `json:"sizeHuman"`
		RepoPath  string `json:"repoPath"`
		RepoName  string `json:"repoName"`
	}
	data, err := json.Marshal(namedHitJSON{
		Type:      "named",
		Path:      hit.Path,
		Name:      hit.Name,
		Size:      hit.Size,
		SizeHuman: hit.SizeHuman,
		RepoPath:  hit.RepoPath,
		RepoName:  hit.RepoName,
	})
	if err != nil {
		return err
	}
	if _, err := w.Write(append(data, '\n')); err != nil {
		return err
	}
	return nil
}
