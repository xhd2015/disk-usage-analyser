package nminventory

import (
	"fmt"
	"io"
	"os"
	"runtime"

	lessflags "github.com/xhd2015/less-flags"
)

type FlagConfig struct {
	HelpText    string
	UsageLine   string
	Stdout      io.Writer
	Stderr      io.Writer
	DefaultHelp func(io.Writer)
}

func ParseFlags(args []string, cfg FlagConfig) (RunOptions, string, int, error) {
	stdout := cfg.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := cfg.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	var (
		workers       = runtime.NumCPU()
		limit         int
		sizeThreshold string
		dryRun        bool
		verbose       bool
	)

	remainArgs, err := lessflags.Int("--workers", &workers).
		Int("--limit", &limit).
		String("--size-threshold", &sizeThreshold).
		Bool("--dry-run", &dryRun).
		Bool("--verbose", &verbose).
		HelpFunc("-h,--help", func() {
			if cfg.DefaultHelp != nil {
				cfg.DefaultHelp(stdout)
			} else {
				fmt.Fprint(stdout, cfg.HelpText)
			}
		}).
		HelpNoExit().
		Parse(args)
	if err == lessflags.ErrHelp {
		return RunOptions{}, "", 0, nil
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return RunOptions{}, "", 1, nil
	}
	if len(remainArgs) != 1 {
		fmt.Fprintf(stderr, "%s\n", cfg.UsageLine)
		return RunOptions{}, "", 2, nil
	}

	thresholdBytes := int64(0)
	if sizeThreshold != "" {
		thresholdBytes, err = ParseCompactHumanSize(sizeThreshold)
		if err != nil {
			fmt.Fprintf(stderr, "invalid --size-threshold %q: %v\n", sizeThreshold, err)
			return RunOptions{}, "", 1, nil
		}
	}

	return RunOptions{
		Workers:       workers,
		Limit:         limit,
		SizeThreshold: thresholdBytes,
		DryRun:        dryRun,
		Verbose:       verbose,
	}, remainArgs[0], -1, nil
}