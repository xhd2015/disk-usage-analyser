package nmscan

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	lessflags "github.com/xhd2015/less-flags"

	"disk-usage-analyser/nmcacheshared"
	"disk-usage-analyser/nminventory"
	"disk-usage-analyser/tmpfiles"

	"github.com/xhd2015/dot-pkgs/go-pkgs/git/scan_repo"
)

const help = `
Usage: node-modules-scan [options]

Scan node_modules directories under ~ and emit JSONL to stdout as each entry
completes. Optionally write the full inventory JSON with -o.

Options:
  -o FILE                write {"version":"1.0","node_modules":[...]} to FILE
  --workers N            concurrent enrichment workers (default: NumCPU)
  --quick                skip shared_size measurement (emit 0B)
  --verbose              progress logs to stderr
  -h, --help             show help
`

type CLIOptions struct {
	Stdout io.Writer
	Stderr io.Writer
	HomeDir string
	Calculator nmcacheshared.CacheCalculator
}

type RunOptions struct {
	OutputPath string
	Workers    int
	Quick      bool
	Verbose    bool
}

type scanHit struct {
	absPath    string
	totalBytes int64
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

	runOpts, code, err := parseFlags(args, stdout, stderr)
	if code >= 0 {
		return code, err
	}

	homeDir := opts.HomeDir
	if homeDir == "" {
		homeDir, err = os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(stderr, "home: %v\n", err)
			return 1, err
		}
	}

	calc := opts.Calculator
	if calc == nil && !runOpts.Quick {
		calc = nmcacheshared.NewCalculator()
	}

	return Run(context.Background(), homeDir, runOpts, calc, stdout, stderr)
}

func parseFlags(args []string, stdout, stderr io.Writer) (RunOptions, int, error) {
	var (
		outputPath string
		workers    = runtime.NumCPU()
		quick      bool
		verbose    bool
	)

	_, err := lessflags.Int("--workers", &workers).
		Bool("--quick", &quick).
		Bool("--verbose", &verbose).
		String("-o", &outputPath).
		HelpFunc("-h,--help", func() {
			fmt.Fprint(stdout, help)
		}).
		HelpNoExit().
		Parse(args)
	if err == lessflags.ErrHelp {
		return RunOptions{}, 0, nil
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return RunOptions{}, 1, nil
	}

	return RunOptions{
		OutputPath: outputPath,
		Workers:    workers,
		Quick:      quick,
		Verbose:    verbose,
	}, -1, nil
}

func Run(ctx context.Context, homeDir string, opts RunOptions, calc nmcacheshared.CacheCalculator, stdout, stderr io.Writer) (int, error) {
	logf := func(format string, args ...any) {
		if opts.Verbose {
			fmt.Fprintf(stderr, "[node-modules-scan] "+format+"\n", args...)
		}
	}

	hits, err := discoverHits(ctx, homeDir, logf)
	if err != nil {
		if ctx.Err() != nil {
			return 1, ctx.Err()
		}
		fmt.Fprintf(stderr, "scan: %v\n", err)
		return 1, err
	}
	if len(hits) == 0 {
		logf("scan complete: 0 hits")
		if opts.OutputPath != "" {
			if err := writeOutputFile(opts.OutputPath, nil); err != nil {
				fmt.Fprintf(stderr, "write output: %v\n", err)
				return 1, err
			}
		}
		return 0, nil
	}

	logf("scan complete: %d hits; enriching with %d workers", len(hits), opts.Workers)

	if !opts.Quick && calc != nil {
		warm := hits[0].absPath
		_ = calc.PnpmCacheShared(warm)
		_ = calc.BunCacheShared(warm)
	}

	records, err := enrichHits(ctx, hits, calc, opts, stdout, logf)
	if err != nil {
		return 1, err
	}

	if opts.OutputPath != "" {
		if err := writeOutputFile(opts.OutputPath, records); err != nil {
			fmt.Fprintf(stderr, "write output: %v\n", err)
			return 1, err
		}
		logf("wrote %d entries to %s", len(records), opts.OutputPath)
	}

	return 0, nil
}

func discoverHits(ctx context.Context, homeDir string, logf func(string, ...any)) ([]scanHit, error) {
	seen := make(map[string]struct{})
	var hits []scanHit
	var hitsMu sync.Mutex

	scanRoots := tmpfiles.PrioritizedScanRoots(homeDir)
	logf("scan roots: %s", strings.Join(scanRoots, ", "))

	_, err := tmpfiles.Scan(ctx, tmpfiles.ScanOptions{
		Roots:   scanRoots,
		Names:   []string{"node_modules"},
		Verbose: false,
		OnNamedHit: func(hit tmpfiles.NamedHit, _ scan_repo.Repo) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			absPath, err := resolveAbsPath(hit.Path, homeDir)
			if err != nil {
				logf("skip path=%q resolve err=%v", hit.Path, err)
				return nil
			}
			key := strings.ToLower(absPath)
			hitsMu.Lock()
			if _, ok := seen[key]; ok {
				hitsMu.Unlock()
				return nil
			}
			seen[key] = struct{}{}
			hits = append(hits, scanHit{absPath: absPath, totalBytes: hit.Size})
			hitsMu.Unlock()
			logf("discovered path=%s size=%d", absPath, hit.Size)
			return nil
		},
	}, homeDir, nil)
	return hits, err
}

func enrichHits(ctx context.Context, hits []scanHit, calc nmcacheshared.CacheCalculator, opts RunOptions, stdout io.Writer, logf func(string, ...any)) ([]nminventory.Record, error) {
	type job struct {
		index int
		hit   scanHit
	}

	workers := opts.Workers
	if workers < 1 {
		workers = 1
	}
	if workers > len(hits) {
		workers = len(hits)
	}

	jobs := make(chan job)
	var wg sync.WaitGroup
	var outMu sync.Mutex
	out := bufio.NewWriter(stdout)
	records := make([]nminventory.Record, len(hits))

	worker := func() {
		defer wg.Done()
		for j := range jobs {
			select {
			case <-ctx.Done():
				return
			default:
			}

			record := BuildRecord(j.hit.absPath, j.hit.totalBytes, calc, opts.Quick)
			line, err := json.Marshal(record)
			if err != nil {
				logf("skip path=%s encode err=%v", j.hit.absPath, err)
				continue
			}

			outMu.Lock()
			_, _ = out.Write(line)
			_ = out.WriteByte('\n')
			_ = out.Flush()
			outMu.Unlock()

			records[j.index] = record

			logf("done path=%s total=%s shared=%s belongs_to_git=%v",
				record.Path, record.TotalSize, record.SharedSize, record.BelongsToGit)
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}
	for i, hit := range hits {
		select {
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			return compactRecords(records), ctx.Err()
		case jobs <- job{index: i, hit: hit}:
		}
	}
	close(jobs)
	wg.Wait()
	return compactRecords(records), out.Flush()
}

func compactRecords(records []nminventory.Record) []nminventory.Record {
	out := make([]nminventory.Record, 0, len(records))
	for _, record := range records {
		if record.Path == "" {
			continue
		}
		out = append(out, record)
	}
	return out
}

func resolveAbsPath(displayPath, homeDir string) (string, error) {
	p := strings.TrimSpace(displayPath)
	if strings.HasPrefix(p, "~/") {
		p = filepath.Join(homeDir, p[2:])
	} else if p == "~" {
		p = homeDir
	}
	return filepath.Abs(filepath.Clean(p))
}

func writeOutputFile(path string, records []nminventory.Record) error {
	if records == nil {
		records = []nminventory.Record{}
	}
	payload := nminventory.OutputFile{
		Version:     "1.0",
		NodeModules: records,
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0644)
}