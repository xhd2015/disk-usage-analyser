# Dev Idle SSE Touch

P4 tests: `sendSSEEvent` (and the shared SSE write path) calls `DevIdleWatch.Touch()` so
long-lived SSE streams reset the dev idle timer. In-process via `server.ServeForTest` with
injectable clock and test-only `/api/test-sse`.

## Version

0.0.2

# DSN (Domain Specific Notion)

When **`dev=true`** and **`DevIdleLife > 0`**, the server wraps the mux with
**`DevIdleWatch.Wrap`** (touch on each HTTP request) and starts the idle ticker. Long scan
handlers stream progress over **SSE** via **`sendSSEEvent`**, which writes `event:` / `data:`
frames to the response.

P4 wires idle activity into the SSE path: each **`sendSSEEvent`** (or equivalent shared
helper) calls **`DevIdleWatch.Touch()`** using a watch reference reachable from SSE handlers
(global on **`serveRuntime`**, request context, or test export). A single long-lived SSE
request therefore keeps **`lastActivity`** fresh across many events, not only the initial
**`Wrap`** touch at request start.

**`/api/test-sse`** is a minimal test-only route registered by **`ServeForTest`**. It emits
**`tick`** events on a fixed interval for a fixed duration, then **`done`**, using the same
injectable **`Now`** clock as the idle watch (no real multi-second sleeps in tests).

## Decision Tree

```
sse/
├── stream-prevents-idle-exit/
└── stream-end-then-idle-exits/
```

## Test Index

| Leaf | Scenario |
|------|----------|
| stream-prevents-idle-exit | SSE `tick` every 500ms for 3s with idle=2s → port still listening |
| stream-end-then-idle-exits | SSE completes, advance 3s idle → port closed |

## How to Run

```sh
doctest vet ./tests/dev-idle/sse
doctest test ./tests/dev-idle/sse/...
```

```go
import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"disk-usage-analyser/server"
)

type Request struct {
	Scenario      string
	DevIdleLife   time.Duration
	TickInterval  time.Duration
	EventInterval time.Duration
	StreamTicks   int
	IdleAdvance   time.Duration
	Port          int
}

type Response struct {
	SSEEventCount int
	PortListening bool
	PortClosed    bool
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

func sseURL(port int, interval time.Duration, ticks int) string {
	ms := interval.Milliseconds()
	return fmt.Sprintf(
		"http://127.0.0.1:%d/api/test-sse?intervalMs=%d&ticks=%d",
		port, ms, ticks,
	)
}

func consumeSSEStream(t *testing.T, url string, wantTicks int, tickCount *atomic.Int32, done chan struct{}) {
	t.Helper()
	defer close(done)

	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("sse GET: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("sse status = %d, want 200; body=%q", resp.StatusCode, body)
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "event: tick") {
			tickCount.Add(1)
		}
		if strings.HasPrefix(line, "event: done") {
			return
		}
	}
	if err := scanner.Err(); err != nil {
		t.Errorf("sse read: %v", err)
	}
	if int(tickCount.Load()) < wantTicks {
		t.Errorf("sse tick count = %d, want >= %d", tickCount.Load(), wantTicks)
	}
}

func waitForTickCount(t *testing.T, tickCount *atomic.Int32, want int, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if int(tickCount.Load()) >= want {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("sse tick count = %d, want >= %d within %v", tickCount.Load(), want, timeout)
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

	interval := req.EventInterval
	if interval <= 0 {
		interval = 500 * time.Millisecond
	}
	ticks := req.StreamTicks
	if ticks <= 0 {
		ticks = 6
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
			NoBrowser:    true,
			SkipFrontend: true,
			EnableTestSSE: true,
		})
		if result != nil {
			watch = result.Watch
			httpServer = result.Server
		}
		serveDone <- err
	}()

	waitForPortOpen(t, port, 5*time.Second)

	var tickCount atomic.Int32
	sseDone := make(chan struct{})
	url := sseURL(port, interval, ticks)

	go consumeSSEStream(t, url, ticks, &tickCount, sseDone)

	resp := &Response{}

	switch req.Scenario {
	case "stream-prevents-idle-exit":
		for i := 1; i <= ticks; i++ {
			clock.Advance(interval)
			pumpIdleChecks(watch, 3)
			waitForTickCount(t, &tickCount, i, 2*time.Second)
		}
		<-sseDone
		resp.SSEEventCount = int(tickCount.Load())
		resp.PortListening = portOpen(port)
		if httpServer != nil {
			httpServer.Close()
		}
		<-serveDone

	case "stream-end-then-idle-exits":
		for i := 1; i <= ticks; i++ {
			clock.Advance(interval)
			pumpIdleChecks(watch, 3)
			waitForTickCount(t, &tickCount, i, 2*time.Second)
		}
		<-sseDone
		resp.SSEEventCount = int(tickCount.Load())

		idleAdvance := req.IdleAdvance
		if idleAdvance <= 0 {
			idleAdvance = 3 * time.Second
		}
		clock.Advance(idleAdvance)
		pumpIdleChecks(watch, 10)
		waitForPortClosed(t, port, 5*time.Second)
		resp.PortClosed = !portOpen(port)
		<-serveDone

	default:
		if httpServer != nil {
			httpServer.Close()
		}
		<-serveDone
		t.Fatalf("unknown sse scenario: %q", req.Scenario)
	}

	return resp, nil
}
```