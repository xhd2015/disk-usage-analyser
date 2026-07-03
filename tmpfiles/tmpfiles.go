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
	OnRepo        func(repo scan_repo.Repo) error
	OnNamedHit    func(hit NamedHit, repo scan_repo.Repo) error
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

// ScanBinaries discovers git repos and classifies binaries, invoking onHit for each match.
func ScanBinaries(ctx context.Context, opts ScanOptions, homeDir string, onHit BinaryHitCallback) (BinaryScanSummary, error) {
	namesSet := makeNameSet(opts.Names)
	summary := BinaryScanSummary{}
	processRepo := func(repo scan_repo.Repo) error {
		summary.Repos++
		if err := scanRepoFiles(ctx, repo, homeDir, namesSet, opts.Verbose, opts.Stderr, func(hit BinaryHit) error {
			summary.Binaries++
			summary.TotalSize += hit.Size
			if onHit != nil {
				return onHit(hit, repo)
			}
			return nil
		}, func(hit NamedHit) error {
			summary.NamedHits++
			summary.TotalSize += hit.Size
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

	_, err := scan_repo.Scan(ctx, scan_repo.Options{
		Roots:    opts.Roots,
		MaxDepth: opts.MaxDepth,
		Verbose:  opts.Verbose,
		Stderr:   opts.Stderr,
		OnRepo:   processRepo,
	})
	if err != nil {
		return summary, err
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
	summary, err := ScanBinaries(ctx, scanOpts, homeDir, func(hit BinaryHit, _ scan_repo.Repo) error {
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

func scanRepoFiles(ctx context.Context, repo scan_repo.Repo, homeDir string, names map[string]bool, verbose bool, stderr io.Writer, onHit func(BinaryHit) error, onNamedHit func(NamedHit) error) error {
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
				if len(names) > 0 && names[d.Name()] {
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
						Name:      d.Name(),
						Size:      size,
						SizeHuman: FormatHumanSize(size),
						RepoPath:  DisplayPath(repo.Path, homeDir),
						RepoName:  repo.Name,
					}
					if err := onNamedHit(hit); err != nil {
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
