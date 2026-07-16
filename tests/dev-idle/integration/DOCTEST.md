# Dev Idle Serve Integration

P3 integration tests: `DevIdleWatch` wired into `server.Serve` / `server.ServeComponent`
when `dev=true`. In-process via `server.ServeForTest` with injectable clock.

## Version

0.0.2

# DSN (Domain Specific Notion)

**`run.Run`** passes **`ServerOptions.DevIdleLife`** into **`server.Serve`** when starting the
web UI in dev mode. **`server.ServeForTest`** (test hook) starts the same idle wiring on an
ephemeral port without bun frontend or browser open.

When **`dev=true`** and **`DevIdleLife > 0`**, the server wraps the mux with
**`DevIdleWatch.Wrap`**, starts the idle ticker, and on expiry runs **`shutdownDev()`**:
**`server.Close()`**, cancel bun context (if started), and log
**`[dev] no requests for <duration>; shutting down`** on stderr. **`DevIdleLife=0`** disables
the watch; the port stays open.

Tests drive time via injectable **`Now`** on the watch (no long real sleeps).

## Decision Tree

```
integration/
├── exits-after-short-idle-life/
├── stays-up-when-disabled/
└── idle-shutdown-log/
```

## Test Index

| Leaf | Scenario |
|------|----------|
| exits-after-short-idle-life | `ServeForTest` dev=true, idle=2s, `/ping`, advance 3s → port closed |
| stays-up-when-disabled | `DevIdleLife=0`, advance 5s → port still listening |
| idle-shutdown-log | idle=1s, capture stderr → shutdown log line |

## How to Run

```sh
doctest vet ./tests/dev-idle/integration
doctest test ./tests/dev-idle/integration/...
```

```go
import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"disk-usage-analyser/server"
)

type Request struct {
	Scenario     string
	DevIdleLife  time.Duration
	TickInterval time.Duration
	Port         int
	Stderr       *bytes.Buffer
}

type Response struct {
	PingStatus    int
	PingBody      string
	PortClosed    bool
	PortListening bool
	Stderr        string
}

type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func newFakeClock(start time.Time) *fakeClock {
	return &fakeClock{now: start}
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

func pumpIdleChecks(watch *server.DevIdleWatch, n int) {
	for i := 0; i < n; i++ {
		server.DevIdleWatchForceCheckForTest(watch)
	}
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

func waitForPortClosed(t *testing.T, port int, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !portOpen(port) {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("port %d still open after %v", port, timeout)
}

func Run(t *testing.T, req *Request) (*Response, error) {
	clock := newFakeClock(time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC))

	port := req.Port
	if port == 0 {
		var err error
		port, err = server.FindAvailablePort(18080, 50)
		if err != nil {
			return nil, err
		}
	}

	stderrBuf := req.Stderr
	if stderrBuf == nil {
		stderrBuf = &bytes.Buffer{}
	}

	serveDone := make(chan error, 1)
	var watch *server.DevIdleWatch
	var httpServer *http.Server

	go func() {
		result, err := server.ServeForTest(port, server.ServeForTestOptions{
			Dev:          true,
			DevIdleLife:  req.DevIdleLife,
			Now:          clock.Now,
			TickInterval: req.TickInterval,
			Stderr:       stderrBuf,
			NoBrowser:    true,
			SkipFrontend: true,
		})
		if result != nil {
			watch = result.Watch
			httpServer = result.Server
		}
		serveDone <- err
	}()

	waitForPortOpen(t, port, 5*time.Second)

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

	resp := &Response{
		PingStatus: pingResp.StatusCode,
		PingBody:   string(pingBody),
		Stderr:     strings.ReplaceAll(stderrBuf.String(), "\r\n", "\n"),
	}

	switch req.Scenario {
	case "exits-after-short-idle-life":
		clock.Advance(3 * time.Second)
		pumpIdleChecks(watch, 10)
		waitForPortClosed(t, port, 5*time.Second)
		resp.PortClosed = !portOpen(port)
		<-serveDone

	case "stays-up-when-disabled":
		clock.Advance(5 * time.Second)
		pumpIdleChecks(watch, 10)
		resp.PortListening = portOpen(port)
		if httpServer != nil {
			httpServer.Close()
		}
		<-serveDone

	case "idle-shutdown-log":
		clock.Advance(2 * time.Second)
		pumpIdleChecks(watch, 10)
		waitForPortClosed(t, port, 5*time.Second)
		resp.PortClosed = !portOpen(port)
		<-serveDone

	default:
		if httpServer != nil {
			httpServer.Close()
		}
		<-serveDone
		t.Fatalf("unknown integration scenario: %q", req.Scenario)
	}

	return resp, nil
}
```