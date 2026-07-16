# Xcode Card — Frontend UI Tests

Playwright UI tests for expanded Xcode breakdown rows, simulator runtime section,
and extended cleanup tooltip on `/tmp-analyse`.

## Version

0.0.2

# DSN (Domain Specific Notion)

The **React tmp-analyse page** renders the Xcode card from SSE **locations** and
**location** events. Five **breakdown-items** rows show DerivedData plus four
ExtraPaths. When **runtimeItems** are present on the xcode location, the existing
**runtime-section** renders per-runtime rows. The cleanup popover for **xcode**
lists DerivedData, iOS DeviceSupport, Archives, DocumentationCache, simulator
device commands, and simulator runtime delete-by-UUID guidance.

The **doctest harness** builds (or will build) a compiled `disk-usage-analyser`
binary, starts it with `--dev`, discovers the listening port from stdout, runs
playwright scripts (or teardown-only checks), then tears down by killing the
**listener PID** — not only a transient `go run` parent.

## Test Tree

```
frontend/
├── breakdown-five-rows/
├── runtime-section/
├── cleanup-popover/
└── server-teardown/
    └── no-listener-orphan/
```

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
doctest vet ./tests/xcode-card/frontend
doctest test ./tests/xcode-card/frontend
doctest test --label ui-automation ./tests/xcode-card/frontend
doctest vet ./tests/xcode-card/frontend/server-teardown
doctest test ./tests/xcode-card/frontend/server-teardown/...
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