package nmmigrate

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"disk-usage-analyser/nminventory"
)

const help = `
Usage: node-modules-migrate-pnpm [options] INPUT.json

Read a node_modules inventory JSON. For each eligible entry (sibling package.json
exists and is git-tracked), remove node_modules and run corepack use pnpm@latest
in the project root. Emits JSONL with operation details per completed entry.

Options:
  --workers N            concurrent workers (default: NumCPU)
  --limit N              process only first N eligible entries after filters (0 = all)
  --size-threshold SIZE  skip entries with total_size below SIZE (e.g. 10M, 10MB)
  --dry-run              emit planned operations as JSONL without removing or running corepack
  --verbose              progress logs to stderr
  -h, --help             show help
`

const corepackCommand = "corepack use pnpm@latest"

type CLIOptions struct {
	Stdout io.Writer
	Stderr io.Writer
	Runner CommandRunner
}

type CommandRunner interface {
	Run(dir string, name string, args ...string) (exitCode int, output string, err error)
}

// ExecRunner runs real shell commands.
type ExecRunner struct{}

func (ExecRunner) Run(dir string, name string, args ...string) (int, string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	if err == nil {
		return 0, buf.String(), nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), buf.String(), nil
	}
	return 1, buf.String(), err
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

	runOpts, inputPath, code, err := nminventory.ParseFlags(args, nminventory.FlagConfig{
		HelpText:  help,
		UsageLine: "usage: node-modules-migrate-pnpm [options] INPUT.json",
		Stdout:    stdout,
		Stderr:    stderr,
	})
	if code >= 0 {
		return code, err
	}

	runner := opts.Runner
	if runner == nil {
		runner = ExecRunner{}
	}

	return RunInventory(inputPath, runOpts, runner, stdout, stderr)
}

func RunInventory(inputPath string, opts nminventory.RunOptions, runner CommandRunner, stdout, stderr io.Writer) (int, error) {
	input, err := nminventory.Load(inputPath)
	if err != nil {
		fmt.Fprintf(stderr, "read input: %v\n", err)
		return 1, err
	}

	logf := func(format string, args ...any) {
		if opts.Verbose || opts.DryRun {
			fmt.Fprintf(stderr, "[node-modules-migrate-pnpm] "+format+"\n", args...)
		}
	}

	parsed, _ := nminventory.ParseEntries(input.NodeModules, logf)
	filtered, belowThreshold := nminventory.FilterBySizeThreshold(parsed, opts.SizeThreshold, logf)
	eligible, ineligible := filterEligible(filtered, logf)
	selected := nminventory.ApplyLimit(eligible, opts.Limit)

	logf("summary: input=%d parsed=%d below_threshold=%d ineligible=%d eligible=%d selected=%d dry_run=%v",
		len(input.NodeModules), len(parsed), belowThreshold, ineligible, len(eligible), len(selected), opts.DryRun)

	if len(selected) == 0 {
		return 0, nil
	}

	if err := migrateEntries(selected, opts, runner, stdout, logf); err != nil {
		return 1, err
	}
	return 0, nil
}

func migrateEntries(entries []nminventory.Entry, opts nminventory.RunOptions, runner CommandRunner, stdout io.Writer, logf func(string, ...any)) error {
	type job struct {
		entry nminventory.Entry
	}

	jobs := make(chan job)
	var wg sync.WaitGroup
	var outMu sync.Mutex
	out := bufio.NewWriter(stdout)

	workers := opts.Workers
	if workers < 1 {
		workers = 1
	}
	if workers > len(entries) {
		workers = len(entries)
	}

	worker := func() {
		defer wg.Done()
		for j := range jobs {
			result := migrateOne(j.entry, opts.DryRun, runner)
			line, err := json.Marshal(result)
			if err != nil {
				logf("skip index=%d path=%q encode err=%v", j.entry.Index, j.entry.Path, err)
				continue
			}

			outMu.Lock()
			_, _ = out.Write(line)
			_ = out.WriteByte('\n')
			_ = out.Flush()
			outMu.Unlock()

			if opts.Verbose || opts.DryRun {
				logf("done index=%d path=%s success=%v dry_run=%v",
					j.entry.Index, j.entry.Path, result["success"], result["dry_run"])
			}
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}
	for _, entry := range entries {
		jobs <- job{entry: entry}
	}
	close(jobs)
	wg.Wait()
	return out.Flush()
}