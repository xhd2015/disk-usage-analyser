package nmcacheshared

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"disk-usage-analyser/analyse"
	"disk-usage-analyser/nminventory"
)

const help = `
Usage: node-modules-cache-shared [options] INPUT.json

Read a node_modules inventory JSON and emit JSONL with pnpm_cache_shared and
bun_cache_shared per entry as each completes.

Options:
  --workers N            concurrent path workers (default: NumCPU)
  --limit N              process only first N entries after filters (0 = all)
  --size-threshold SIZE  skip entries with total_size below SIZE (e.g. 10M, 10MB)
  --dry-run              trace planned work to stderr; skip counting
  --verbose              progress logs to stderr
  -h, --help             show help
`

type CLIOptions struct {
	Stdout     io.Writer
	Stderr     io.Writer
	Calculator CacheCalculator
}

type CacheCalculator interface {
	PnpmCacheShared(path string) int64
	BunCacheShared(path string) int64
}

type analyseCalculator struct {
	inner *analyse.CacheSharedCalculator
}

func (c *analyseCalculator) PnpmCacheShared(path string) int64 {
	return c.inner.PnpmCacheShared(path)
}

func (c *analyseCalculator) BunCacheShared(path string) int64 {
	return c.inner.BunCacheShared(path)
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
		UsageLine: "usage: node-modules-cache-shared [options] INPUT.json",
		Stdout:    stdout,
		Stderr:    stderr,
	})
	if code >= 0 {
		return code, err
	}

	calc := opts.Calculator
	if calc == nil && !runOpts.DryRun {
		calc = &analyseCalculator{inner: analyse.NewCacheSharedCalculator()}
	}

	return RunInventory(inputPath, runOpts, calc, stdout, stderr)
}

func RunInventory(inputPath string, opts nminventory.RunOptions, calc CacheCalculator, stdout, stderr io.Writer) (int, error) {
	input, err := nminventory.Load(inputPath)
	if err != nil {
		fmt.Fprintf(stderr, "read input: %v\n", err)
		return 1, err
	}

	logf := func(format string, args ...any) {
		if opts.Verbose || opts.DryRun {
			fmt.Fprintf(stderr, "[node-modules-cache-shared] "+format+"\n", args...)
		}
	}

	parsed, _ := nminventory.ParseEntries(input.NodeModules, logf)
	filtered, belowThreshold := nminventory.FilterBySizeThreshold(parsed, opts.SizeThreshold, logf)
	selected := nminventory.ApplyLimit(filtered, opts.Limit)

	if opts.DryRun {
		if opts.SizeThreshold > 0 {
			logf("dry-run summary: input=%d parsed=%d below_threshold=%d filtered=%d limit=%d would_scan=%d",
				len(input.NodeModules), len(parsed), belowThreshold, len(filtered), opts.Limit, len(selected))
		} else {
			logf("dry-run summary: input=%d parsed=%d limit=%d would_scan=%d",
				len(input.NodeModules), len(parsed), opts.Limit, len(selected))
		}
		for _, entry := range selected {
			logf("dry-run would scan index=%d path=%s total_size=%s", entry.Index, entry.Path, entry.TotalSize)
		}
		return 0, nil
	}

	if len(selected) == 0 {
		return 0, nil
	}
	if calc == nil {
		return 1, fmt.Errorf("calculator required when not dry-run")
	}

	logf("indexes warming: selected=%d workers=%d", len(selected), opts.Workers)

	warmPath := selected[0].Path
	if warmPath == "" {
		warmPath = filepath.Join(os.TempDir(), "node_modules")
	}
	_ = calc.PnpmCacheShared(warmPath)
	_ = calc.BunCacheShared(warmPath)

	if err := scanEntries(selected, calc, opts, stdout, logf); err != nil {
		return 1, err
	}
	return 0, nil
}

func scanEntries(entries []nminventory.Entry, calc CacheCalculator, opts nminventory.RunOptions, stdout io.Writer, logf func(string, ...any)) error {
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
			pnpm := calc.PnpmCacheShared(j.entry.Path)
			bun := calc.BunCacheShared(j.entry.Path)

			obj, err := nminventory.DecodeRaw(j.entry.Raw)
			if err != nil {
				logf("skip index=%d path=%q decode err=%v", j.entry.Index, j.entry.Path, err)
				continue
			}
			obj["pnpm_cache_shared"] = nminventory.FormatCompactHumanSize(pnpm)
			obj["bun_cache_shared"] = nminventory.FormatCompactHumanSize(bun)

			line, err := json.Marshal(obj)
			if err != nil {
				logf("skip index=%d path=%q encode err=%v", j.entry.Index, j.entry.Path, err)
				continue
			}

			outMu.Lock()
			_, _ = out.Write(line)
			_ = out.WriteByte('\n')
			_ = out.Flush()
			outMu.Unlock()

			if opts.Verbose {
				logf("done index=%d path=%s pnpm=%s bun=%s",
					j.entry.Index, j.entry.Path, obj["pnpm_cache_shared"], obj["bun_cache_shared"])
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