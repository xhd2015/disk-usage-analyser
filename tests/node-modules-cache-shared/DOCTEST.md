# node-modules-cache-shared CLI

CLI tests for `node-modules-cache-shared`: less-flags parsing, size-threshold
filtering, limit, dry-run tracing, and JSONL emission.

## Version

0.0.1

# DSN (Domain Specific Notion)

**node-modules-cache-shared** reads a versioned inventory JSON (`version`,
`node_modules[]` with `path`, `total_size`, etc.), filters entries by optional
`--size-threshold` (binary sizes), caps with `--limit`, then either traces
planned work (`--dry-run` → stderr only) or scans each path for APFS clone overlap
with the pnpm store and bun install cache, emitting JSONL rows with
`pnpm_cache_shared` and `bun_cache_shared` appended.

## Test Index

| Leaf | Mode | Description |
|------|------|-------------|
| cli/help | cli | `--help` documents flags and exits 0 |
| dry-run/trace | cli | `--dry-run` traces to stderr, empty stdout |
| filter/size-threshold | cli | `--size-threshold 10M` skips sub-10M in dry-run summary |
| filter/limit | cli | `--limit 1 --dry-run` schedules one path |

## How to Run

```sh
doctest vet ./tests/node-modules-cache-shared
doctest test ./tests/node-modules-cache-shared
```

```go
import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"disk-usage-analyser/nmcacheshared"
)

type Request struct {
	Args   []string
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Err      error
}

func inventoryPath(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for dir := wd; ; dir = filepath.Dir(dir) {
		candidate := filepath.Join(dir, "testdata", "inventory.json")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("testdata/inventory.json not found")
		}
	}
	return ""
}

func Run(t *testing.T, req *Request) (*Response, error) {
	if req.Stdout == nil {
		req.Stdout = &bytes.Buffer{}
	}
	if req.Stderr == nil {
		req.Stderr = &bytes.Buffer{}
	}
	exitCode, err := nmcacheshared.RunCLI(req.Args, nmcacheshared.CLIOptions{
		Stdout: req.Stdout,
		Stderr: req.Stderr,
	})
	return &Response{
		Stdout:   strings.ReplaceAll(req.Stdout.String(), "\r\n", "\n"),
		Stderr:   strings.ReplaceAll(req.Stderr.String(), "\r\n", "\n"),
		ExitCode: exitCode,
		Err:      err,
	}, nil
}
```