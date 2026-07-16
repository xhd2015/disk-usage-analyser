# Dev Idle Watchdog

P1 tests for `server.DevIdleWatch`: injectable clock, `Touch`, `Wrap`, background idle
ticker, and `OnIdle` callback. P2 tests for `--dev-idle-life` CLI flag plumbing into
`run.ServerOptions.DevIdleLife`. P3 integration tests wire `DevIdleWatch` into
`server.Serve` / `server.ServeComponent` when `dev=true`. P4 SSE tests assert
`sendSSEEvent` touches the idle watch so long streams prevent shutdown. P5 subprocess tests
exercise the full CLI path via a real compiled binary and real wall-clock idle shutdown.

## Version

0.0.5

# DSN (Domain Specific Notion)

The **DevIdleWatch** guards a dev server process against lingering after traffic stops.
It holds a configurable **Timeout**, an **OnIdle** shutdown callback, and an injectable
**Now** clock (defaults to `time.Now`). **Start** launches a background goroutine that
ticks periodically and compares `Now()` against **lastActivity**. **Touch** records
activity; the idle clock does not run until the first touch. **Wrap** wraps an
`http.Handler`, calling **Touch** before delegating to the next handler.

When **Timeout** is zero or negative, the watch is disabled: **Start** is a no-op and
**OnIdle** is never invoked. When enabled, each **Touch** resets **lastActivity**. If
`Now() - lastActivity >= Timeout`, **OnIdle** fires once per idle episode.

The **CLI root** (`run.RunWithOptions`) parses global server flags before starting the web
UI. **`--dev`** enables development mode. **`--dev-idle-life <duration>`** sets how long
the dev server may remain idle before shutdown; values are parsed with `time.ParseDuration`.
The literals **`0`** and **`off`** disable idle life (zero duration). When **`--dev`** is
set and **`--dev-idle-life`** is omitted, the default is **10 minutes**. When **`--dev`**
is not set, **`--dev-idle-life`** is ignored and **DevIdleLife** stays zero. Invalid
duration strings are rejected with an error. Root **`-h`** help documents **`--dev-idle-life`**.

Parsed duration is passed to **`run.ServerOptions.DevIdleLife`** via a fake **StartServer**
hook in P2 flag tests. In P3, **`run.Run`** passes **`DevIdleLife`** into **`server.Serve`**
and **`server.ServeComponent`**. When **`dev=true`** and **`DevIdleLife > 0`**, the server
wraps the mux with **`DevIdleWatch.Wrap`**, starts the idle ticker, and on expiry runs a
shared **`shutdownDev()`** path: **`server.Close()`**, cancel bun context (if a frontend
dev child was started), and log **`[dev] no requests for <duration>; shutting down`** to
stderr. Signal handling (SIGTERM) and idle shutdown share one shutdown path.

Long-lived **SSE** handlers stream progress via **`sendSSEEvent`**. P4 requires each SSE write
to call **`DevIdleWatch.Touch()`** so periodic events reset **`lastActivity`** during scans.
After the stream ends, the idle clock resumes from the last touch and shutdown can fire.

## Decision Tree

```
dev-idle/
├── integration/          # nested DOCTEST.md (P3 — separate Run contract)
│   ├── exits-after-short-idle-life/
│   ├── stays-up-when-disabled/
│   └── idle-shutdown-log/
├── sse/                  # nested DOCTEST.md (P4 — SSE touch contract)
│   ├── stream-prevents-idle-exit/
│   └── stream-end-then-idle-exits/
├── subprocess/           # nested DOCTEST.md (P5 — real CLI subprocess)
│   ├── exits-after-short-idle-life/
│   └── idle-shutdown-log/
├── watchdog/
│   ├── disabled-zero-timeout/
│   ├── exits-after-idle/
│   ├── traffic-resets-idle/
│   └── wrap-handler-touches/
└── flag/
    ├── default-10m-when-dev/
    ├── override-1h/
    ├── disabled-off/
    ├── ignored-without-dev/
    ├── invalid-duration-errors/
    └── help-lists-flag/
```

## Test Index

| Leaf | Scenario |
|------|----------|
| integration/exits-after-short-idle-life | `ServeForTest` dev=true, idle=2s, `/ping`, advance 3s → port closed |
| integration/stays-up-when-disabled | `DevIdleLife=0`, advance 5s → port still listening |
| integration/idle-shutdown-log | idle=1s, capture stderr → shutdown log line |
| sse/stream-prevents-idle-exit | SSE tick every 500ms for 3s, idle=2s → port still up |
| sse/stream-end-then-idle-exits | SSE ends, advance 3s idle → port closed |
| subprocess/exits-after-short-idle-life | build binary; `--dev --dev-idle-life 2s`; sleep 3s → port closed, process exited |
| subprocess/idle-shutdown-log | `--dev --dev-idle-life 1s`; capture stderr → shutdown log line |
| watchdog/disabled-zero-timeout | Timeout=0, advance fake clock, pump checks |
| watchdog/exits-after-idle | Timeout=2s, one touch, advance 3s |
| watchdog/traffic-resets-idle | Timeout=3s, touches at T=0 and T=2s, advance to T=4s |
| watchdog/wrap-handler-touches | Wrap + httptest request, idle only after request |
| flag/default-10m-when-dev | `--dev` only → `ServerOptions.DevIdleLife == 10m` |
| flag/override-1h | `--dev --dev-idle-life 1h` → `DevIdleLife == 1h` |
| flag/disabled-off | `--dev --dev-idle-life off` → `DevIdleLife == 0` (`0` equivalent) |
| flag/ignored-without-dev | `--dev-idle-life 1h` without `--dev` → `DevIdleLife == 0` |
| flag/invalid-duration-errors | `--dev --dev-idle-life not-a-duration` → parse error, no server |
| flag/help-lists-flag | root `-h` stdout contains `--dev-idle-life`, no server |

## How to Run

```sh
doctest vet ./tests/dev-idle
doctest test ./tests/dev-idle/...
doctest vet ./tests/dev-idle/integration
doctest test ./tests/dev-idle/integration/...
doctest vet ./tests/dev-idle/sse
doctest test ./tests/dev-idle/sse/...
doctest vet ./tests/dev-idle/subprocess
doctest test ./tests/dev-idle/subprocess/...
doctest test ./tests/dev-idle/flag/...
doctest test ./tests/dev-idle/watchdog/...
```

```go
import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"disk-usage-analyser/run"
	"disk-usage-analyser/server"
)

type Request struct {
	Scenario     string
	Timeout      time.Duration
	TickInterval time.Duration

	// P2 flag plumbing via run.RunWithOptions
	Args   []string
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

type Response struct {
	OnIdleCount          int
	OnIdleCountAfterHTTP int
	HTTPStatus           int

	Stdout           string
	Stderr           string
	Err              error
	ServerWasStarted bool
	DevIdleLife      time.Duration
	ServerDev        bool
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

func newWatch(t *testing.T, req *Request, clock *fakeClock, onIdleCount *atomic.Int32) *server.DevIdleWatch {
	t.Helper()
	tick := req.TickInterval
	if tick <= 0 {
		tick = 50 * time.Millisecond
	}
	watch := &server.DevIdleWatch{
		Timeout: req.Timeout,
		Now:     clock.Now,
		OnIdle: func() {
			onIdleCount.Add(1)
		},
	}
	server.DevIdleWatchConfigureForTest(watch, tick)
	return watch
}

func pumpIdleChecks(watch *server.DevIdleWatch, n int) {
	for i := 0; i < n; i++ {
		server.DevIdleWatchForceCheckForTest(watch)
	}
}

func runFlagCLI(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req.Stdout == nil {
		req.Stdout = &bytes.Buffer{}
	}
	if req.Stderr == nil {
		req.Stderr = &bytes.Buffer{}
	}

	var captured run.ServerOptions
	serverStarted := false

	err := run.RunWithOptions(context.Background(), req.Args, run.Options{
		Stdout: req.Stdout,
		Stderr: req.Stderr,
		StartServer: func(_ context.Context, opts run.ServerOptions) error {
			serverStarted = true
			captured = opts
			return nil
		},
	})

	return &Response{
		Stdout:           strings.ReplaceAll(req.Stdout.String(), "\r\n", "\n"),
		Stderr:           strings.ReplaceAll(req.Stderr.String(), "\r\n", "\n"),
		Err:              err,
		ServerWasStarted: serverStarted,
		DevIdleLife:      captured.DevIdleLife,
		ServerDev:        captured.Dev,
	}, nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	switch {
	case strings.HasPrefix(req.Scenario, "flag/"):
		return runFlagCLI(t, req)

	case req.Scenario == "disabled-zero-timeout":
		clock := newFakeClock(time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC))
		var onIdleCount atomic.Int32

		watch := newWatch(t, req, clock, &onIdleCount)
		watch.Start()
		defer server.DevIdleWatchStopForTest(watch)

		watch.Touch()
		clock.Advance(10 * time.Minute)
		pumpIdleChecks(watch, 5)

		return &Response{OnIdleCount: int(onIdleCount.Load())}, nil

	case req.Scenario == "exits-after-idle":
		clock := newFakeClock(time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC))
		var onIdleCount atomic.Int32

		watch := newWatch(t, req, clock, &onIdleCount)
		watch.Start()
		defer server.DevIdleWatchStopForTest(watch)

		watch.Touch()
		clock.Advance(3 * time.Second)
		pumpIdleChecks(watch, 3)

		return &Response{OnIdleCount: int(onIdleCount.Load())}, nil

	case req.Scenario == "traffic-resets-idle":
		clock := newFakeClock(time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC))
		var onIdleCount atomic.Int32

		watch := newWatch(t, req, clock, &onIdleCount)
		watch.Start()
		defer server.DevIdleWatchStopForTest(watch)

		watch.Touch()
		clock.Advance(2 * time.Second)
		watch.Touch()
		clock.Advance(2 * time.Second)
		pumpIdleChecks(watch, 3)

		return &Response{OnIdleCount: int(onIdleCount.Load())}, nil

	case req.Scenario == "wrap-handler-touches":
		clock := newFakeClock(time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC))
		var onIdleCount atomic.Int32

		watch := newWatch(t, req, clock, &onIdleCount)
		watch.Start()
		defer server.DevIdleWatchStopForTest(watch)

		handler := watch.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		srv := httptest.NewServer(handler)
		defer srv.Close()

		httpResp, err := http.Get(srv.URL)
		if err != nil {
			return nil, err
		}
		httpResp.Body.Close()

		afterHTTP := int(onIdleCount.Load())
		clock.Advance(3 * time.Second)
		pumpIdleChecks(watch, 3)

		return &Response{
			OnIdleCount:          int(onIdleCount.Load()),
			OnIdleCountAfterHTTP: afterHTTP,
			HTTPStatus:           httpResp.StatusCode,
		}, nil

	default:
		t.Fatalf("unknown dev-idle scenario: %q", req.Scenario)
		return nil, nil
	}
}
```