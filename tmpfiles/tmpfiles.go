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

type ScanResult struct {
	Roots      []string
	Repos      int
	Binaries   []BinaryHit
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
}

const help = `Usage: disk-usage-analyser tmp-files scan [OPTIONS]

Options:
  --go-binaries     Scan for Go, Mach-O, and ELF binaries
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
		switch arg {
		case "-h", "--help":
			return opts, true, nil
		case "--go-binaries":
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
	Roots    []string
	MaxDepth int
	Verbose  bool
	Stderr   io.Writer
	OnRepo   func(repo scan_repo.Repo) error
}

// BinaryScanSummary holds aggregate binary scan results.
type BinaryScanSummary struct {
	Repos       int
	Binaries    int
	TotalSize   int64
	TotalHuman  string
}

// BinaryHitCallback is invoked for each discovered binary.
type BinaryHitCallback func(hit BinaryHit, repo scan_repo.Repo) error

// ScanBinaries discovers git repos and classifies binaries, invoking onHit for each match.
func ScanBinaries(ctx context.Context, opts ScanOptions, homeDir string, onHit BinaryHitCallback) (BinaryScanSummary, error) {
	summary := BinaryScanSummary{}
	processRepo := func(repo scan_repo.Repo) error {
		summary.Repos++
		if err := scanRepoFiles(ctx, repo, homeDir, opts.Verbose, opts.Stderr, func(hit BinaryHit) error {
			summary.Binaries++
			summary.TotalSize += hit.Size
			if onHit != nil {
				return onHit(hit, repo)
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
	}
	var result ScanResult
	result.Roots = opts.roots
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
	if _, err := fmt.Fprintf(stdout, "Found %d binaries, total %s\n", len(result.Binaries), result.TotalHuman); err != nil {
		return &result, err
	}
	return &result, flush(stdout)
}

func scanRepoFiles(ctx context.Context, repo scan_repo.Repo, homeDir string, verbose bool, stderr io.Writer, onHit func(BinaryHit) error) error {
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
