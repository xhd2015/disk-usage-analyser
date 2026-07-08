package usagescan

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	lessflags "github.com/xhd2015/less-flags"
)

const help = `Usage: disk-usage-analyser scan [PATH]

Walk a directory recursively and print a size-annotated tree.

Options:
  --json            Emit one JSON object with nested tree
  --threshold SIZE  Hide nodes below SIZE (default: 1M; 0 = show all)
  --max-depth N     Max tree depth (default: 3 text, 24 with --json; 0 = unlimited)
  -h, --help        Show help
`

const defaultThreshold = "1M"

func RunCLI(args []string, opts CLIOptions) (int, error) {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	var jsonOut bool
	var thresholdStr string
	var maxDepth int
	var scanPath string

	maxDepth, maxDepthSet := maxDepthFromArgs(args)

	remain, err := lessflags.
		Bool("--json", &jsonOut).
		String("--threshold", &thresholdStr).
		Int("--max-depth", &maxDepth).
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

	if len(remain) > 1 {
		return 2, fmt.Errorf("unexpected extra argument: %s", remain[1])
	}
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

	if thresholdStr == "" {
		thresholdStr = defaultThreshold
	}
	threshold, err := ParseCompactHumanSize(thresholdStr)
	if err != nil {
		return 2, err
	}

	if !maxDepthSet {
		if jsonOut {
			maxDepth = 24
		} else {
			maxDepth = 3
		}
	}

	result, err := Scan(scanPath, ScanOptions{
		Threshold: threshold,
		MaxDepth:  maxDepth,
	})
	if err != nil {
		return 2, err
	}

	if jsonOut {
		return writeJSON(stdout, result)
	}

	if _, err := fmt.Fprint(stdout, formatTreeText(result)); err != nil {
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

func maxDepthFromArgs(args []string) (int, bool) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--max-depth" {
			if i+1 >= len(args) {
				continue
			}
			v, err := strconv.Atoi(args[i+1])
			if err == nil {
				return v, true
			}
		}
		if strings.HasPrefix(arg, "--max-depth=") {
			v, err := strconv.Atoi(strings.TrimPrefix(arg, "--max-depth="))
			if err == nil {
				return v, true
			}
		}
	}
	return 0, false
}