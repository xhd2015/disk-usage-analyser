# Tmp Files Scan

CLI tests for `disk-usage-analyser tmp-files scan`, covering dispatch, root
resolution, git repository discovery, Go/Mach-O/ELF classification, ignored
paths, output formats, streaming, and graceful error handling.

## Version

0.0.2

# DSN (Domain Specific Notion)

The **tmp-files scan command** is a CLI-only subcommand below `tmp-files`; it
must bypass the web server dispatcher. **Root resolution** turns repeatable
`--root` flags, or the default home directory, into scan roots. **Repo
discovery** uses `scan_repo.Scan` with `MaxDepth` and default ignored basenames.
**Remote-backed filesystem detection** identifies macOS cloud-sync roots such as
`~/Library/CloudStorage/...`, prints a stderr warning, and skips them before
walking so `fdopendir` does not hang on remote-backed paths. For every
discovered git repository, **binary discovery** walks regular files,
skips `.git` and ignored directories, uses `detect.DetectFileType`, and gives
`debug/buildinfo.ReadFile` precedence over file magic so Go binaries classify as
`go` even when their container is Mach-O or ELF. **BinaryHit** records path,
size, human size, kind, type description, repo path, and repo name. **Named-hit
detection** uses a repeatable `--name=NAME` flag; when a directory basename
matches, it computes the recursive size (skipping nested same-name subdirs),
reports a **NamedHit**, and skips the subtree. When a regular file basename
matches, it reports the file's direct size. Named hits and binary hits are
additive in a single scan. **ScanResult** aggregates roots, repo count, binary
hits, named hits, and totals. The CLI streams one human or NDJSON line per hit
immediately, flushes after each hit, then prints the summary line
`Found N binaries, M named entries, total <human-size>`.

## Decision Tree

```
tmp-files-scan/
├── command/
│   ├── help/                  # scan -h documents flags and exits 0
│   ├── parent-help/
│   │   ├── long/              # tmp-files --help (parent level) → usage, exit 0
│   │   └── short/             # tmp-files -h (parent level) → usage, exit 0
│   └── cli-dispatch/          # run.Run dispatches tmp-files without web server
├── roots/
│   ├── default-root-tilde/    # no --root scans fixture HOME
│   └── custom-root/           # repeatable explicit --root limits scope
├── discovery/
│   ├── skip-outside-repo/     # binary outside .git repo is ignored
│   ├── multi-repos/           # multiple repos preserve RepoPath/RepoName
│   ├── respect-ignore-dirs/   # vendor/node_modules/.venv paths are skipped
│   ├── max-depth/             # deep repos outside MaxDepth are undiscovered
│   ├── permission-denied/     # unreadable subtree yields partial results
│   ├── skip-cloud-storage/           # CloudStorage paths skip silently by default
│   ├── skip-cloud-storage-verbose/   # -v warns when CloudStorage is skipped
│   ├── warn-remote-root/             # explicit --root on CloudStorage skips
│   └── no-repos/              # empty discovery emits zero summary
├── classification/
│   ├── single-repo-macho/     # Mach-O binary hit
│   ├── go-binary/             # buildinfo hit takes Kind=go precedence
│   ├── elf-binary/            # ELF binary hit
│   └── skip-text-files/       # .go and .txt files are not hits
├── output/
│   ├── human-output/          # default streaming human format
│   ├── json-output/           # --json NDJSON hit lines plus text summary
│   ├── streaming/             # first Write is a hit, not buffered summary
│   └── tilde-paths/           # home-relative paths use ~/...
└── named/
    ├── match-dir/
    │   ├── single/
    │   │   ├── basic/              # single dir match, human output, correct size
    │   │   ├── nested/             # nested dirs produce separate hits
    │   │   ├── override-ignore/    # --name overrides ignoredDirBasenames
    │   │   ├── additive/           # named dir + binary in same scan
    │   │   ├── json/               # --json NDJSON with type:named
    │   │   ├── size-accuracy/      # recursive size matches expected
    │   │   └── custom-root/        # --root limits named search scope
    │   └── multiple/
    │       ├── both-found/         # two names both found in same repo
    │       └── mixed-ignore/       # named ok, non-named ignored dir still skipped
    ├── match-file/
    │   └── basic/                  # regular file matches --name
    ├── no-matches/                 # nothing matches --name
    └── no-repos/                   # no git repos at all
```

## Test Index

| Leaf | Description |
|------|-------------|
| command/help | `tmp-files scan -h` documents `--go-binaries`, `--root`, `--max-depth`, and `--json`. |
| command/parent-help/long | Parent-level `tmp-files --help` (no `scan` token) prints the same usage family as empty/`scan -h`, exit 0, trailing `\n`, no scan. |
| command/parent-help/short | Parent-level `tmp-files -h` (no `scan` token) prints the same usage family as empty/`scan -h`, exit 0, trailing `\n`, no scan. |
| command/cli-dispatch | `run.Run` dispatches `tmp-files scan` through the CLI path and never starts the web server. |
| roots/default-root-tilde | With no `--root`, the scan uses the fixture home directory and renders `~` paths. |
| roots/custom-root | Explicit `--root` scans only the selected subtree. |
| discovery/skip-outside-repo | Binary files outside git repositories are not reported. |
| discovery/multi-repos | Two repositories produce two hits with correct repository metadata. |
| discovery/respect-ignore-dirs | Binaries below ignored basenames such as `vendor` are skipped. |
| discovery/max-depth | Repositories deeper than `--max-depth` are not discovered. |
| discovery/permission-denied | Permission-denied paths are skipped gracefully while other repos are reported. |
| discovery/skip-cloud-storage | `~/Library/CloudStorage/...` is skipped silently; local repos still scan. |
| discovery/skip-cloud-storage-verbose | `-v` warns on stderr when CloudStorage paths are skipped. |
| discovery/warn-remote-root | `--root` on a CloudStorage provider skips without reporting hits. |
| discovery/no-repos | A root with no `.git` directories emits no hits and a zero summary. |
| classification/single-repo-macho | A Mach-O fixture is classified as `macho` and streamed. |
| classification/go-binary | A real Go-built fixture is classified as `go`. |
| classification/elf-binary | An ELF fixture is classified as `elf`. |
| classification/skip-text-files | Source and text files inside a repo are ignored. |
| output/human-output | Default output has size, kind, path, repo path, and summary. |
| output/json-output | `--json` emits valid NDJSON hit objects followed by the text summary. |
| output/streaming | The first stdout write contains a hit line before the final summary is written. |
| output/tilde-paths | All home-contained paths are rendered with a `~/` prefix. |
| named/match-dir/single/basic | Single `node_modules` dir reported with path, size, repo metadata. |
| named/match-dir/single/nested | Nested `node_modules` dirs produce two separate hits with correct size exclusion. |
| named/match-dir/single/override-ignore | `--name=node_modules` overrides `ignoredDirBasenames` to report the dir. |
| named/match-dir/single/additive | Repo with `node_modules` + Mach-O binary reports both hits. |
| named/match-dir/single/json | `--json` emits NDJSON with `"type":"named"` for named hits. |
| named/match-dir/single/size-accuracy | Recursive size of known file contents matches expected value. |
| named/match-dir/single/custom-root | `--root` limits named dir search to selected subtree. |
| named/match-dir/multiple/both-found | `--name=node_modules --name=vendor` finds both in same repo. |
| named/match-dir/multiple/mixed-ignore | `--name=node_modules` reports it; `vendor` (not named) is still skipped. |
| named/match-file/basic | Regular file named `node_modules` is reported as a named hit. |
| named/no-matches | `--name=nonexistent` produces zero named hits and a zero summary. |
| named/no-repos | Root with no `.git` dirs produces zero named hits with `--name`.

## How to Run

```sh
doctest vet ./tests/tmp-files-scan
doctest test ./tests/tmp-files-scan
```

```go
import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"disk-usage-analyser/run"
	"disk-usage-analyser/tmpfiles"
)

var _ = os.Setenv("GO111MODULE", "off")

type Request struct {
	Op           string
	Args         []string
	HomeDir      string
	Stdout       *bytes.Buffer
	Stderr       *bytes.Buffer
	StreamStdout *recordingWriter
}

type Response struct {
	Stdout          string
	Stderr          string
	Result          *tmpfiles.ScanResult
	ExitCode        int
	Err             error
	FirstWrite      string
	WriteCount       int
	ServerWasStarted bool
}

type recordingWriter struct {
	*io.PipeWriter
	mu         sync.Mutex
	buf        bytes.Buffer
	firstWrite string
	writeCount int
	done       chan struct{}
}

func recordingWrite(w *recordingWriter, p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.writeCount == 0 {
		w.firstWrite = string(p)
	}
	w.writeCount++
	return w.buf.Write(p)
}

func recordingStart(w *recordingWriter) {
	if w.PipeWriter != nil {
		return
	}
	reader, writer := io.Pipe()
	w.PipeWriter = writer
	w.done = make(chan struct{})
	go func() {
		defer close(w.done)
		buf := make([]byte, 64*1024)
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				_, _ = recordingWrite(w, buf[:n])
			}
			if err != nil {
				return
			}
		}
	}()
}

func newRecordingWriter() *recordingWriter {
	w := &recordingWriter{}
	recordingStart(w)
	return w
}

func recordingClose(w *recordingWriter) {
	if w.PipeWriter == nil {
		return
	}
	_ = w.PipeWriter.Close()
	if w.done != nil {
		<-w.done
	}
}

func recordingString(w *recordingWriter) string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.String()
}

func recordingFirstWrite(w *recordingWriter) string {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.firstWrite == "" {
		text := strings.ReplaceAll(w.buf.String(), "\r\n", "\n")
		if idx := strings.IndexByte(text, '\n'); idx >= 0 {
			return text[:idx+1]
		}
		return text
	}
	return w.firstWrite
}

func recordingWriteCount(w *recordingWriter) int {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.writeCount == 0 {
		text := strings.ReplaceAll(w.buf.String(), "\r\n", "\n")
		return strings.Count(text, "\n")
	}
	return w.writeCount
}

func Run(t *testing.T, req *Request) (*Response, error) {
	if req.Stdout == nil {
		req.Stdout = &bytes.Buffer{}
	}
	if req.Stderr == nil {
		req.Stderr = &bytes.Buffer{}
	}
	stdout := any(req.Stdout).(interface{ Write([]byte) (int, error) })
	if req.StreamStdout != nil {
		recordingStart(req.StreamStdout)
		stdout = req.StreamStdout
	}

	if req.Op == "cli-dispatch" {
		serverStarted := false
		err := run.RunWithOptions(context.Background(), req.Args, run.Options{
			Stdout: stdout,
			Stderr: req.Stderr,
			HomeDir: req.HomeDir,
			StartServer: func(context.Context, run.ServerOptions) error {
				serverStarted = true
				t.Fatalf("tmp-files dispatch must not start web server")
				return nil
			},
		})
		if req.StreamStdout != nil {
			recordingClose(req.StreamStdout)
		}
		return &Response{
			Stdout: req.Stdout.String(),
			Stderr: req.Stderr.String(),
			Err: err,
			ServerWasStarted: serverStarted,
		}, nil
	}

	result, exitCode, err := tmpfiles.RunCLI(context.Background(), req.Args, tmpfiles.CLIOptions{
		Stdout: stdout,
		Stderr: req.Stderr,
		HomeDir: req.HomeDir,
	})
	if req.StreamStdout != nil {
		recordingClose(req.StreamStdout)
	}

	resp := &Response{
		Result: result,
		ExitCode: exitCode,
		Err: err,
		Stderr: req.Stderr.String(),
	}
	if req.StreamStdout != nil {
		resp.Stdout = recordingString(req.StreamStdout)
		resp.FirstWrite = recordingFirstWrite(req.StreamStdout)
		resp.WriteCount = recordingWriteCount(req.StreamStdout)
	} else {
		resp.Stdout = req.Stdout.String()
	}
	resp.Stdout = strings.ReplaceAll(resp.Stdout, "\r\n", "\n")
	return resp, nil
}
```
