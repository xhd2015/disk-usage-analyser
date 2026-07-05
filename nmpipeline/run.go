package nmpipeline

import (
	"fmt"
	"io"
	"os"
	"sync"

	"disk-usage-analyser/nmcacheshared"
	"disk-usage-analyser/nminventory"
	"disk-usage-analyser/nmmigrate"
)

const help = `
Usage: node-modules-migration-report [options] [--node-modules-records FILE] [PATH...]

Run before cache-shared analyse, migrate eligible node_modules to pnpm via corepack,
then post analyse, and print an aligned summary table.

Provide at least one node_modules PATH and/or --node-modules-records FILE.

Options:
  --node-modules-records FILE
                         inventory JSON (e.g. node_modules_on_the_system.json)
  --workers N            concurrent workers (default: NumCPU)
  --limit N              process only first N eligible entries after filters (0 = all)
  --size-threshold SIZE  skip entries with total_size below SIZE (e.g. 10M, 10MB)
  --dry-run              measure before only; skip migrate and keep after=before
  --verbose              progress logs to stderr
  -h, --help             show help
`

const usageLine = "usage: node-modules-migration-report [options] [--node-modules-records FILE] PATH..."

type CLIOptions struct {
	Stdout     io.Writer
	Stderr     io.Writer
	Runner     nmmigrate.CommandRunner
	Calculator nmcacheshared.CacheCalculator
}

func RunCLI(args []string, opts CLIOptions) (int, error) {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	runOpts, input, code, err := ParsePipelineFlags(args, help, usageLine, stdout, stderr)
	if code >= 0 {
		return code, err
	}

	runner := opts.Runner
	if runner == nil {
		runner = nmmigrate.ExecRunner{}
	}
	calc := opts.Calculator
	if calc == nil {
		calc = nmcacheshared.NewCalculator()
	}

	return RunReport(input, runOpts, runner, calc, stdout, stderr)
}

func RunReport(input Input, opts nminventory.RunOptions, runner nmmigrate.CommandRunner, calc nmcacheshared.CacheCalculator, stdout, stderr io.Writer) (int, error) {
	logf := func(format string, args ...any) {
		if opts.Verbose {
			fmt.Fprintf(stderr, "[node-modules-migration-report] "+format+"\n", args...)
		}
	}

	selected, err := resolveEntries(input, opts, logf)
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		return 1, err
	}

	if len(selected) == 0 {
		fmt.Fprintln(stdout, "No eligible paths to process.")
		return 0, nil
	}

	if !opts.DryRun {
		logf("warming cache indexes...")
		_ = calc.PnpmCacheShared(selected[0].Path)
		_ = calc.BunCacheShared(selected[0].Path)
	}

	rows := runPipeline(selected, opts, runner, calc, logf)
	sortRowsBySharedAdded(rows)

	if err := FormatTable(stdout, rows); err != nil {
		return 1, err
	}
	return 0, nil
}

func runPipeline(entries []nminventory.Entry, opts nminventory.RunOptions, runner nmmigrate.CommandRunner, calc nmcacheshared.CacheCalculator, logf func(string, ...any)) []Row {
	type job struct {
		entry nminventory.Entry
	}
	type result struct {
		row Row
	}

	jobs := make(chan job)
	results := make(chan result)
	workers := opts.Workers
	if workers < 1 {
		workers = 1
	}
	if workers > len(entries) {
		workers = len(entries)
	}

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				logf("start path=%s", j.entry.Path)
				row := processEntry(j.entry, opts.DryRun, runner, calc)
				logf("done path=%s shared_added=%s migrate_ok=%v",
					j.entry.Path, formatSignedSize(row.SharedAdded()), row.Migrate.Success)
				results <- result{row: row}
			}
		}()
	}

	go func() {
		for _, entry := range entries {
			jobs <- job{entry: entry}
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	rows := make([]Row, 0, len(entries))
	for r := range results {
		rows = append(rows, r.row)
	}
	return rows
}