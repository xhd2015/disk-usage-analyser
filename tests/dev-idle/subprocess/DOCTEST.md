# Dev Idle Subprocess CLI

P5 subprocess tests: end-to-end CLI path `disk-usage-analyser --dev --dev-idle-life <dur>`
shuts down on idle using a **real compiled binary**, **real subprocess**, and **real wall-clock**
waits. Distinct from P3 `ServeForTest` in-process integration tests.

## Version

0.0.2

# DSN (Domain Specific Notion)

The **CLI root** (`run.Run`) parses **`--dev`** and **`--dev-idle-life`**, then starts the web
UI via **`server.Serve`**. In subprocess tests the harness **`go build -o`** the module binary
once per session, **`exec.Command`** launches it with **`NO_BROWSER=1`**, and reads stdout for
**`Serving directory preview at http://localhost:<port>`** to discover the listening port.

An HTTP **`GET /ping`** (body **`pong`**) touches **`DevIdleWatch`** through the wrapped mux.
After real wall-clock idle exceeds **`--dev-idle-life`**, the background ticker invokes
**`shutdownDev()`**: **`server.Close()`**, cancel bun context (if started), and log
**`[dev] no requests for <duration>; shutting down`** on stderr. The subprocess exits; TCP dials
to the former port fail.

Tests use **real `time.Sleep`** (1–3s) — no injectable fake clock. Teardown kills the server
PID if still running (cleanup only; idle shutdown is what we assert).

## Decision Tree

```
subprocess/
├── exits-after-short-idle-life/
└── idle-shutdown-log/
```

## Test Index

| Leaf | Scenario |
|------|----------|
| exits-after-short-idle-life | build binary; `--dev --dev-idle-life 2s`; `GET /ping`; sleep past tick → port closed, process exited |
| idle-shutdown-log | `--dev --dev-idle-life 1s`; capture stderr; sleep past idle → shutdown log line |

## How to Run

```sh
doctest vet ./tests/dev-idle/subprocess
doctest test ./tests/dev-idle/subprocess/...
```

```go
import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

type Request struct {
	Scenario    string
	DevIdleLife time.Duration
	IdleSleep   time.Duration
	BinPath     string
}

type Response struct {
	PingStatus    int
	PingBody      string
	Port          int
	PortClosed    bool
	ProcessExited bool
	ExitCode      int
	Stdout        string
	Stderr        string
	ServerPID     int
}

func moduleRoot() string {
	return filepath.Clean(filepath.Join(DOCTEST_ROOT, "..", "..", ".."))
}

func sessionCacheDir() string {
	return filepath.Join(os.TempDir(), "dev-idle-subprocess-doctest-"+DOCTEST_SESSION_ID)
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func buildBinaryOnce(t *testing.T) string {
	t.Helper()
	cacheDir := sessionCacheDir()
	bin := filepath.Join(cacheDir, "disk-usage-analyser")
	lock := filepath.Join(cacheDir, "build.lock")
	ready := filepath.Join(cacheDir, "binaries.ready")

	err := withFileLock(t, lock, func() error {
		if fileExists(ready) && fileExists(bin) {
			return nil
		}
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			return err
		}
		cmd := exec.Command("go", "build", "-o", bin, ".")
		cmd.Dir = moduleRoot()
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

func parsePort(line string) (int, bool) {
	prefix := "Serving directory preview at http://localhost:"
	if !strings.HasPrefix(line, prefix) {
		return 0, false
	}
	rest := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	port, err := strconv.Atoi(rest)
	if err != nil {
		return 0, false
	}
	return port, true
}

func portOpen(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 200*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func waitForPortOpen(t *testing.T, port int, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if portOpen(port) {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("port %d did not open within %v", port, timeout)
}

func waitForProcessExit(t *testing.T, cmd *exec.Cmd, timeout time.Duration) int {
	t.Helper()
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case err := <-done:
		if err == nil {
			return 0
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		t.Fatalf("wait subprocess: %v", err)
		return -1
	case <-time.After(timeout):
		t.Fatalf("subprocess pid %d did not exit within %v", cmd.Process.Pid, timeout)
		return -1
	}
}

func killServerIfRunning(t *testing.T, cmd *exec.Cmd) {
	t.Helper()
	if cmd == nil || cmd.Process == nil {
		return
	}
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	_, _ = cmd.Process.Wait()
}

func startDevServer(t *testing.T, req *Request) (*exec.Cmd, int, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()

	idleLife := req.DevIdleLife
	if idleLife <= 0 {
		t.Fatalf("DevIdleLife must be positive for subprocess scenarios, got %v", idleLife)
	}

	stdoutBuf := &bytes.Buffer{}
	stderrBuf := &bytes.Buffer{}

	cmd := exec.Command(
		req.BinPath,
		"--dev",
		"--dev-idle-life", idleLife.String(),
	)
	cmd.Dir = moduleRoot()
	cmd.Env = append(os.Environ(), "NO_BROWSER=1")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("start subprocess: %v", err)
	}

	go func() {
		_, _ = io.Copy(stderrBuf, stderrPipe)
	}()

	portCh := make(chan int, 1)
	errCh := make(chan error, 1)
	go func() {
		reader := bufio.NewReader(stdoutPipe)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					errCh <- err
				}
				return
			}
			stdoutBuf.WriteString(line)
			if port, ok := parsePort(line); ok {
				portCh <- port
			}
		}
	}()

	var port int
	select {
	case port = <-portCh:
	case err := <-errCh:
		killServerIfRunning(t, cmd)
		t.Fatalf("read server stdout: %v", err)
	case <-time.After(60 * time.Second):
		killServerIfRunning(t, cmd)
		t.Fatal("server did not print listening port within 60s")
	}

	waitForPortOpen(t, port, 10*time.Second)
	return cmd, port, stdoutBuf, stderrBuf
}

// subprocessRunMu serializes parallel leaves: each subprocess calls FindAvailablePort
// inside the binary and concurrent starts can race on the same port.
var subprocessRunMu sync.Mutex

func Run(t *testing.T, req *Request) (*Response, error) {
	subprocessRunMu.Lock()
	defer subprocessRunMu.Unlock()

	if req.BinPath == "" {
		req.BinPath = buildBinaryOnce(t)
	}

	cmd, port, stdoutBuf, stderrBuf := startDevServer(t, req)
	defer killServerIfRunning(t, cmd)

	pingURL := fmt.Sprintf("http://127.0.0.1:%d/ping", port)
	pingResp, err := http.Get(pingURL)
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	pingBody, err := io.ReadAll(pingResp.Body)
	if err != nil {
		return nil, err
	}
	pingResp.Body.Close()

	idleSleep := req.IdleSleep
	if idleSleep <= 0 {
		idleSleep = req.DevIdleLife + time.Second
	}
	time.Sleep(idleSleep)

	resp := &Response{
		PingStatus: pingResp.StatusCode,
		PingBody:   string(pingBody),
		Port:       port,
		PortClosed: !portOpen(port),
		ServerPID:  cmd.Process.Pid,
		Stdout:     strings.ReplaceAll(stdoutBuf.String(), "\r\n", "\n"),
		Stderr:     strings.ReplaceAll(stderrBuf.String(), "\r\n", "\n"),
	}

	switch req.Scenario {
	case "exits-after-short-idle-life":
		resp.ExitCode = waitForProcessExit(t, cmd, 15*time.Second)
		resp.ProcessExited = true
		resp.PortClosed = !portOpen(port)

	case "idle-shutdown-log":
		resp.ExitCode = waitForProcessExit(t, cmd, 15*time.Second)
		resp.ProcessExited = true
		resp.PortClosed = !portOpen(port)

	default:
		t.Fatalf("unknown subprocess scenario: %q", req.Scenario)
	}

	return resp, nil
}
```