package usagescan

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	lessflags "github.com/xhd2015/less-flags"
)

const help = `Usage: disk-usage-analyser scan [PATH]

Walk a directory recursively and print a size-annotated tree, or load a
saved TreeResult JSON with --inspect and query it offline (no re-walk).

Options:
  --inspect FILE    Load TreeResult JSON (FILE may be "-" for stdin)
  --json            Emit JSON (bare TreeResult for live capture; ViewResult otherwise)
  --min SIZE        Hide tree/match nodes below SIZE (default: live 1M, inspect 0; 0 = show all)
  --max-depth N     Max tree depth (default: live text 3, live --json capture 24, inspect 1; 0 = unlimited)
  --top N           Emit TOP N match ranking (default N=20 when match section is active)
  --at PATH         Focus the tree section on PATH (--at alone has no TOP section)
  --find SUBSTR     Match path/name containing SUBSTR (case-insensitive) → match section
  --suffix SUFFIX   Match name/path ending with SUFFIX → match section
  --include-root    Include scan root in global TOP / find rankings
  -h, --help        Show help
`

const defaultLiveMin = "1M"

func RunCLI(args []string, opts CLIOptions) (int, error) {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	var (
		jsonOut     bool
		minStr      string
		maxDepth    int
		top         int
		atPath      string
		find        string
		suffix      string
		includeRoot bool
		inspectFile string
	)

	maxDepth, maxDepthSet := intFlagFromArgs(args, "--max-depth")
	top, topSet := intFlagFromArgs(args, "--top")
	minSet := stringFlagPresent(args, "--min")
	inspectSet := stringFlagPresent(args, "--inspect")

	remain, err := lessflags.
		Bool("--json", &jsonOut).
		String("--min", &minStr).
		Int("--max-depth", &maxDepth).
		Int("--top", &top).
		String("--at", &atPath).
		String("--find", &find).
		String("--suffix", &suffix).
		Bool("--include-root", &includeRoot).
		String("--inspect", &inspectFile).
		HelpFunc("-h,--help", func() {
			fmt.Fprint(stdout, help)
		}).
		HelpNoExit().
		Parse(args)
	if err == lessflags.ErrHelp {
		return 0, nil
	}
	if err != nil {
		return 2, err
	}

	inspectMode := inspectSet || inspectFile != ""
	if inspectMode && inspectFile == "" {
		return 2, fmt.Errorf("--inspect requires a FILE argument (or \"-\")")
	}

	if inspectMode {
		if len(remain) > 0 {
			return 2, fmt.Errorf("--inspect does not accept a positional PATH (got %q)", remain[0])
		}
	} else if len(remain) > 1 {
		return 2, fmt.Errorf("unexpected extra argument: %s", remain[1])
	}

	wantMatches := topSet || find != "" || suffix != ""
	needView := inspectMode || wantMatches || atPath != ""

	// Defaults for min / max-depth
	var min int64
	if minSet {
		min, err = ParseCompactHumanSize(minStr)
		if err != nil {
			return 2, fmt.Errorf("invalid --min: %w", err)
		}
	} else if inspectMode {
		min = 0
	} else {
		min, err = ParseCompactHumanSize(defaultLiveMin)
		if err != nil {
			return 2, err
		}
	}

	if !maxDepthSet {
		if inspectMode {
			maxDepth = 1
		} else if jsonOut && !needView {
			// pure live JSON capture
			maxDepth = 24
		} else {
			// live text or live query view
			maxDepth = 3
		}
	}

	if atPath != "" {
		atPath = expandLeadingTilde(atPath)
	}

	// Phase 1: TreeSource
	var source TreeSource
	var sourceFile string
	if inspectMode {
		source = JSONTreeSource{Path: inspectFile}
		if inspectFile != "" && inspectFile != "-" {
			if abs, err := filepath.Abs(inspectFile); err == nil {
				sourceFile = abs
			} else {
				sourceFile = inspectFile
			}
		}
	} else {
		scanPath := ""
		if len(remain) == 1 {
			scanPath = remain[0]
		}
		if scanPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return 1, fmt.Errorf("cannot determine current directory: %w", err)
			}
			scanPath = cwd
		}
		scanMin, scanDepth := min, maxDepth
		if wantMatches {
			// Load full tree for match ranking; view prunes for the tree section.
			scanMin, scanDepth = 0, 0
		}
		source = LiveTreeSource{
			Path: scanPath,
			Opts: ScanOptions{Min: scanMin, MaxDepth: scanDepth},
		}
	}

	result, err := source.Load()
	if err != nil {
		return 2, err
	}

	// Pure live capture/text: tree already filtered by Scan; no view extras.
	if !needView {
		// Ensure metadata reflects CLI defaults/overrides for pure live.
		result.Min = min
		result.MaxDepth = maxDepth
		if jsonOut {
			return writeJSON(stdout, result)
		}
		if _, err := fmt.Fprint(stdout, formatTreeText(result)); err != nil {
			return 1, err
		}
		return 0, nil
	}

	// Phase 2: shared view
	viewOpts := ViewOptions{
		Min:         min,
		MaxDepth:    maxDepth,
		Top:         top,
		TopSet:      topSet,
		AtPath:      atPath,
		Find:        find,
		Suffix:      suffix,
		IncludeRoot: includeRoot,
		SourceFile:  sourceFile,
	}
	view := BuildView(result, viewOpts)

	if jsonOut {
		return writeViewJSON(stdout, view)
	}
	if _, err := fmt.Fprint(stdout, FormatViewText(view, viewOpts)); err != nil {
		return 1, err
	}
	return 0, nil
}

func writeJSON(stdout io.Writer, result TreeResult) (int, error) {
	payload, err := json.Marshal(result)
	if err != nil {
		return 1, fmt.Errorf("encode json: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, string(payload)); err != nil {
		return 1, err
	}
	if _, err := fmt.Fprintln(stdout); err != nil {
		return 1, err
	}
	return 0, nil
}

func writeViewJSON(stdout io.Writer, result ViewResult) (int, error) {
	payload, err := json.Marshal(result)
	if err != nil {
		return 1, fmt.Errorf("encode json: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, string(payload)); err != nil {
		return 1, err
	}
	if _, err := fmt.Fprintln(stdout); err != nil {
		return 1, err
	}
	return 0, nil
}

func intFlagFromArgs(args []string, name string) (int, bool) {
	prefix := name + "="
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == name {
			if i+1 >= len(args) {
				return 0, true
			}
			v, err := strconv.Atoi(args[i+1])
			if err == nil {
				return v, true
			}
			return 0, true
		}
		if strings.HasPrefix(arg, prefix) {
			v, err := strconv.Atoi(strings.TrimPrefix(arg, prefix))
			if err == nil {
				return v, true
			}
			return 0, true
		}
	}
	return 0, false
}

func stringFlagPresent(args []string, name string) bool {
	prefix := name + "="
	for _, arg := range args {
		if arg == name || strings.HasPrefix(arg, prefix) {
			return true
		}
	}
	return false
}
