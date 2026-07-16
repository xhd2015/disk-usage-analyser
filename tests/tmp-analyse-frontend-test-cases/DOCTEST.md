# Tmp Files Analyse — Frontend Test Cases

UI tests for the tmp-analyse page: structure, navigation, scan interaction, breakdown
display, runtime stats sections, live breakdown progress, and Repository Scans
(worktrees, binaries, node_modules, and vendor sections with independent scan
controls, filter/sort/select/delete UX).

## Version

0.0.2

# DSN (Domain Specific Notion)

The **React tmp-analyse page** loads location cards from the initial SSE **locations**
event, then starts a scan via **Start Scan**. During scanning, **progress** SSE events
update card header sizes and (for multi-path cards) individual **breakdown rows** live.
When each location completes, a **location** event finalizes sizes and may attach
**runtimeItems** for Docker/Podman, and **vmInternal** for Podman VM storage on macOS.
The page renders cards with **data-testid** hooks for headings, summary bar, breakdown
table rows, vm-internal section rows, runtime section rows, and scan badges.
Repository Scans sections (worktrees, binaries, node_modules, vendor) each have
independent scan buttons, SSE event streams, and delete endpoints. Tests drive the
page through **playwright-debug** scripts against a locally started Go dev server.
Repository Scans leaves use `ui-automation` and `slow` labels when they depend on
completing a subsection scan.

The **doctest harness** builds (or will build) a compiled `disk-usage-analyser`
binary, starts it with `--dev`, discovers the listening port from stdout, runs
playwright scripts (or teardown-only checks), then tears down by killing the
**listener PID** — not only a transient `go run` parent.

## Test Tree

```
tmp-analyse-frontend-test-cases/
├── [existing] verify-page-renders … verify-npm-breakdown
├── breakdown-live-progress/
│   └── go-multi-path/
├── runtime-section/
│   ├── docker-after-scan/
│   └── podman-after-scan/
├── podman-vm-internal/
│   └── after-scan/
├── repository-scans/
│   └── renders/
├── worktrees-section/
│   ├── after-scan/
│   ├── left-aligned/
│   ├── filter-under-10m-default/
│   ├── filter-show-under-10m/
│   ├── sort-by-size-desc/
│   └── running-total/
├── worktrees-live-stream/
├── binaries-section/
│   ├── after-scan/
│   ├── left-aligned/
│   ├── filter-under-1m-default/
│   ├── filter-show-under-1m/
│   ├── sort-by-size-desc/
│   ├── select-and-total/
│   ├── delete-selected/
│   ├── repo-select-all/
│   └── running-total/
├── named-section/
│   ├── after-scan/
│   ├── left-aligned/
│   ├── filter-under-1m-default/
│   ├── filter-show-under-1m/
│   ├── sort-by-size-desc/
│   ├── select-and-total/
│   ├── delete-selected/
│   ├── repo-select-all/
│   ├── running-total/
│   ├── enrichment/
│   │   ├── rows-before-shared/
│   │   ├── per-row-shared-loading/
│   │   ├── shared-accumulates-during-enrichment/
│   │   └── first-row-timing/
│   ├── table-columns/
│   │   ├── after-scan/
│   │   ├── full-width/
│   │   ├── pm-no-wrap/
│   │   ├── package-json-column/
│   │   └── path-visible-limit/
│   ├── path-truncation/
│   │   ├── tooltip-full-path/
│   │   ├── prefix-suffix-visible/
│   │   ├── tooltip-copy-button/
│   │   └── truncate-path-keep-suffix/
│   │       ├── fits-within-limit/
│   │       ├── truncates-with-prefix-ellipsis/
│   │       └── prefers-slash-boundary/
│   └── column-filters/
│       ├── filter-named-repos-by-column/
│       │   ├── defaults-show-all/
│       │   ├── git-yes/
│       │   ├── git-no/
│       │   ├── package-json-yes/
│       │   ├── package-json-no/
│       │   ├── pm-pnpm/
│       │   ├── pm-unknown/
│       │   └── combined-filters/
│       ├── controls-present/
│       ├── git-no-hides-tracked/
│       └── pm-filter-npm/
├── vendor-section/
│   ├── after-scan/
│   └── running-total/
├── named-independent-scans/
├── independent-scan-controls/
└── server-teardown/
    └── no-listener-orphan/
```

## Test Index

| Leaf | Script |
|------|--------|
| verify-page-renders | page-renders.js |
| verify-breakdown-live-progress/go-multi-path | breakdown-live-progress.js |
| verify-runtime-section/docker-after-scan | docker-runtime-section.js |
| verify-runtime-section/podman-after-scan | podman-runtime-section.js |
| podman-vm-internal/after-scan | podman-vm-internal-section.js |
| repository-scans/renders | repository-scans-renders.js |
| worktrees-section/after-scan | worktrees-after-scan.js |
| worktrees-section/left-aligned | worktrees-left-aligned.js |
| worktrees-section/filter-under-10m-default | worktrees-filter-under-10m-default.js |
| worktrees-section/filter-show-under-10m | worktrees-filter-show-under-10m.js |
| worktrees-section/sort-by-size-desc | worktrees-sort-by-size-desc.js |
| worktrees-section/running-total | worktrees-running-total.js |
| worktrees-live-stream | worktrees-live-stream.js |
| binaries-section/after-scan | binaries-after-scan.js |
| binaries-section/left-aligned | binaries-left-aligned.js |
| binaries-section/filter-under-1m-default | binaries-filter-under-1m-default.js |
| binaries-section/filter-show-under-1m | binaries-filter-show-under-1m.js |
| binaries-section/sort-by-size-desc | binaries-sort-by-size-desc.js |
| binaries-section/select-and-total | binaries-select-total.js |
| binaries-section/delete-selected | binaries-delete-selected.js |
| binaries-section/repo-select-all | binaries-repo-select-all.js |
| binaries-section/running-total | binaries-running-total.js |
| named-section/after-scan | named-after-scan.js |
| named-section/left-aligned | named-left-aligned.js |
| named-section/filter-under-1m-default | named-filter-under-1m-default.js |
| named-section/filter-show-under-1m | named-filter-show-under-1m.js |
| named-section/sort-by-size-desc | named-sort-by-size-desc.js |
| named-section/select-and-total | named-select-total.js |
| named-section/delete-selected | named-delete-selected.js |
| named-section/repo-select-all | named-repo-select-all.js |
| named-section/running-total | node-modules-running-total.js |
| named-section/enrichment/rows-before-shared | rows-before-shared.js |
| named-section/enrichment/per-row-shared-loading | per-row-shared-loading.js |
| named-section/enrichment/shared-accumulates-during-enrichment | shared-accumulates.js |
| named-section/enrichment/first-row-timing | first-row-timing.js |
| named-section/table-columns/after-scan | table-columns-after-scan.js |
| named-section/table-columns/full-width | table-columns-full-width.js |
| named-section/table-columns/pm-no-wrap | table-columns-pm-no-wrap.js |
| named-section/table-columns/package-json-column | table-columns-package-json-column.js |
| named-section/table-columns/path-visible-limit | path-visible-limit-harness.ts |
| named-section/path-truncation/tooltip-full-path | path-truncation-tooltip.js |
| named-section/path-truncation/prefix-suffix-visible | path-prefix-suffix-visible.js |
| named-section/path-truncation/tooltip-copy-button | path-tooltip-copy-button.js |
| named-section/path-truncation/truncate-path-keep-suffix/fits-within-limit | path-display-harness.ts |
| named-section/path-truncation/truncate-path-keep-suffix/truncates-with-prefix-ellipsis | path-display-harness.ts |
| named-section/path-truncation/truncate-path-keep-suffix/prefers-slash-boundary | path-display-harness.ts |
| named-section/column-filters/filter-named-repos-by-column/defaults-show-all | column-filters-harness.ts |
| named-section/column-filters/filter-named-repos-by-column/git-yes | column-filters-harness.ts |
| named-section/column-filters/filter-named-repos-by-column/git-no | column-filters-harness.ts |
| named-section/column-filters/filter-named-repos-by-column/package-json-yes | column-filters-harness.ts |
| named-section/column-filters/filter-named-repos-by-column/package-json-no | column-filters-harness.ts |
| named-section/column-filters/filter-named-repos-by-column/pm-pnpm | column-filters-harness.ts |
| named-section/column-filters/filter-named-repos-by-column/pm-unknown | column-filters-harness.ts |
| named-section/column-filters/filter-named-repos-by-column/combined-filters | column-filters-harness.ts |
| named-section/column-filters/controls-present | column-filters-controls-present.js |
| named-section/column-filters/git-no-hides-tracked | column-filters-git-no-hides-tracked.js |
| named-section/column-filters/pm-filter-npm | column-filters-pm-filter-npm.js |
| vendor-section/after-scan | vendor-after-scan.js |
| vendor-section/running-total | vendor-running-total.js |
| named-independent-scans | named-independent-scans.js |
| independent-scan-controls | independent-scan-controls.js |
| server-teardown/no-listener-orphan | (harness-only; no playwright) |

## Harness contract (P6 — implementer)

Current `Run()` still uses `go run` and `SIGTERM` on the wrapper PID, which can
orphan the compiled `--dev` listener (PPID=1). Implementer must replace startup
and teardown with:

| Helper | Contract |
|--------|----------|
| `buildDevServerBinary(t, projectRoot) string` | `go build -o $tmpdir/disk-usage-analyser .`; session-cache via `DOCTEST_SESSION_ID` |
| `startDevServer(t, bin) (port string, pid int, cleanup func())` | `exec.Command(bin, "--dev")` with `Setpgid: true`; read port from stdout; `pid` is **listener** on that port |
| `cleanup` | Kill listener PID or process group (`Kill(-pgid, SIGKILL)`); `waitForPortClosed(port)` |

Teardown-only leaves set `req.TeardownOnly` and skip playwright. After `cleanup`,
`Run` records whether the port still accepts TCP and whether the listener PID is
alive — both must be false.

## How to Run

```sh
doctest vet ./tests/tmp-analyse-frontend-test-cases
doctest test ./tests/tmp-analyse-frontend-test-cases
doctest vet ./tests/tmp-analyse-frontend-test-cases/server-teardown
doctest test ./tests/tmp-analyse-frontend-test-cases/server-teardown/...
```

```go
import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

type Request struct {
	ScriptFile    string
	TeardownOnly  bool
}

type Response struct {
	Output                       string
	Passed                       bool
	ServerPort                   string
	ListenerPID                  int
	PortListeningAfterTeardown   bool
	ListenerAliveAfterTeardown   bool
}

func sessionCacheDir() string {
	return filepath.Join(os.TempDir(), "frontend-doctest-"+DOCTEST_SESSION_ID)
}

func withFileLock(t *testing.T, lockPath string, fn func() error) error {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(lockPath), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return fn()
}

func buildDevServerBinary(t *testing.T, projectRoot string) string {
	t.Helper()
	cacheDir := sessionCacheDir()
	bin := filepath.Join(cacheDir, "disk-usage-analyser")
	lock := filepath.Join(cacheDir, "build.lock")
	ready := filepath.Join(cacheDir, "binaries.ready")

	err := withFileLock(t, lock, func() error {
		if _, err := os.Stat(ready); err == nil {
			if _, err := os.Stat(bin); err == nil {
				return nil
			}
		}
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			return err
		}
		cmd := exec.Command("go", "build", "-o", bin, ".")
		cmd.Dir = projectRoot
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("go build -o %s: %w\n%s", bin, err, string(output))
		}
		return os.WriteFile(ready, []byte("ok"), 0644)
	})
	if err != nil {
		t.Fatal(err)
	}
	return bin
}

func startDevServer(t *testing.T, projectRoot, bin string) (port string, listenerPID int, cleanup func()) {
	t.Helper()

	serverCmd := exec.Command(bin, "--dev")
	serverCmd.Dir = projectRoot
	serverCmd.Env = append(os.Environ(), "NO_BROWSER=1")
	serverCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, err := serverCmd.StdoutPipe()
	if err != nil {
		t.Fatalf("failed to get stdout pipe: %v", err)
	}

	if err := serverCmd.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	portCh := make(chan string, 1)
	errCh := make(chan error, 1)

	reader := bufio.NewReader(stdout)
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					errCh <- err
				}
				return
			}
			if p, ok := parsePort(line); ok {
				portCh <- p
				return
			}
		}
	}()

	select {
	case port = <-portCh:
	case err := <-errCh:
		killServerProcessGroup(serverCmd)
		t.Fatalf("failed to read server output: %v", err)
	case <-time.After(60 * time.Second):
		killServerProcessGroup(serverCmd)
		t.Fatal("server did not start within 60 seconds")
	}

	go func() {
		io.Copy(io.Discard, stdout)
	}()

	waitForPortOpen(t, port, 10*time.Second)
	listenerPID = listenerPIDForPort(t, port)

	cleanup = func() {
		if listenerPID > 0 {
			if pgid, err := syscall.Getpgid(listenerPID); err == nil && pgid > 0 {
				_ = syscall.Kill(-pgid, syscall.SIGKILL)
			} else {
				_ = syscall.Kill(listenerPID, syscall.SIGKILL)
			}
		}
		killServerProcessGroup(serverCmd)
		waitForPortClosed(t, port, 3*time.Second)
	}

	return port, listenerPID, cleanup
}

func killServerProcessGroup(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	_, _ = cmd.Process.Wait()
}

func Run(t *testing.T, req *Request) (*Response, error) {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to find project root: %v", err)
	}

	bin := buildDevServerBinary(t, projectRoot)
	port, listenerPID, teardown := startDevServer(t, projectRoot, bin)

	if req.TeardownOnly {
		teardown()
		return &Response{
			ServerPort:                 port,
			ListenerPID:                listenerPID,
			PortListeningAfterTeardown: portOpen(port),
			ListenerAliveAfterTeardown: processAlive(listenerPID),
		}, nil
	}

	defer teardown()

	script, err := os.ReadFile(req.ScriptFile)
	if err != nil {
		return nil, err
	}

	serverURL := fmt.Sprintf("http://localhost:%s", port)

	cmd := exec.Command("playwright-debug", string(script))
	cmd.Env = append(os.Environ(), "SERVER_URL="+serverURL, "NO_BROWSER=1")
	out, runErr := cmd.CombinedOutput()

	passed := cmd.ProcessState != nil && cmd.ProcessState.ExitCode() == 0

	return &Response{
		Output: string(out),
		Passed: passed,
	}, runErr
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found in any parent directory")
		}
		dir = parent
	}
}

func parsePort(line string) (string, bool) {
	prefix := "Serving directory preview at http://localhost:"
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	rest := strings.TrimPrefix(line, prefix)
	return strings.TrimSpace(rest), true
}

func portOpen(port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), 200*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func waitForPortOpen(t *testing.T, port string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if portOpen(port) {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("port %s did not open within %v", port, timeout)
}

func waitForPortClosed(t *testing.T, port string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !portOpen(port) {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func listenerPIDForPort(t *testing.T, port string) int {
	t.Helper()
	cmd := exec.Command("lsof", "-nP", "-iTCP:"+port, "-sTCP:LISTEN", "-t")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	line := strings.TrimSpace(string(out))
	if line == "" {
		return 0
	}
	pidStr := strings.Split(line, "\n")[0]
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		t.Fatalf("parse listener pid %q: %v", pidStr, err)
	}
	return pid
}

func processAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	return syscall.Kill(pid, 0) == nil
}
```