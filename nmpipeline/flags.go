package nmpipeline

import (
	"fmt"
	"io"
	"os"
	"runtime"

	lessflags "github.com/xhd2015/less-flags"

	"disk-usage-analyser/nminventory"
)

// Input lists node_modules sources from CLI args and optional records JSON.
type Input struct {
	RecordsFile string
	Paths       []string
}

func ParsePipelineFlags(args []string, helpText, usageLine string, stdout, stderr io.Writer) (nminventory.RunOptions, Input, int, error) {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	var (
		workers         = runtime.NumCPU()
		limit           int
		sizeThreshold   string
		dryRun          bool
		verbose         bool
		nodeModulesFile string
	)

	remainArgs, err := lessflags.Int("--workers", &workers).
		Int("--limit", &limit).
		String("--size-threshold", &sizeThreshold).
		String("--node-modules-records", &nodeModulesFile).
		Bool("--dry-run", &dryRun).
		Bool("--verbose", &verbose).
		HelpFunc("-h,--help", func() {
			fmt.Fprint(stdout, helpText)
		}).
		HelpNoExit().
		Parse(args)
	if err == lessflags.ErrHelp {
		return nminventory.RunOptions{}, Input{}, 0, nil
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return nminventory.RunOptions{}, Input{}, 1, nil
	}

	if len(remainArgs) == 0 && nodeModulesFile == "" {
		fmt.Fprintf(stderr, "%s\n", usageLine)
		return nminventory.RunOptions{}, Input{}, 2, nil
	}

	thresholdBytes := int64(0)
	if sizeThreshold != "" {
		thresholdBytes, err = nminventory.ParseCompactHumanSize(sizeThreshold)
		if err != nil {
			fmt.Fprintf(stderr, "invalid --size-threshold %q: %v\n", sizeThreshold, err)
			return nminventory.RunOptions{}, Input{}, 1, nil
		}
	}

	return nminventory.RunOptions{
		Workers:       workers,
		Limit:         limit,
		SizeThreshold: thresholdBytes,
		DryRun:        dryRun,
		Verbose:       verbose,
	}, Input{
		RecordsFile: nodeModulesFile,
		Paths:       remainArgs,
	}, -1, nil
}