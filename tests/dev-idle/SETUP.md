# Scenario

**Feature**: DevIdleWatch idle watchdog with injectable fake clock

```
# first HTTP touch starts idle clock; each touch resets lastActivity
HTTP request -> DevIdleWatch.Wrap -> Touch -> next handler

# background ticker compares Now() - lastActivity against Timeout
DevIdleWatch.Start -> ticker -> idle? -> OnIdle (once)

# tests drive time via injectable Now (no real 10m sleep)
fake clock -> DevIdleWatch.Now -> ForceCheckForTest pumps ticker logic
```

## Preconditions

- Target package: `disk-usage-analyser/server` (`server/dev_idle.go` not yet implemented).
- Harness uses injectable **fake clock** via `DevIdleWatch.Now`; no long real sleeps.
- Test hooks expected from implementer (`dev_idle_export_test.go` or equivalent):
  - `DevIdleWatchConfigureForTest(w *DevIdleWatch, tickInterval time.Duration)`
  - `DevIdleWatchForceCheckForTest(w *DevIdleWatch)`
  - `DevIdleWatchStopForTest(w *DevIdleWatch)`
- Default tick interval in harness: 50ms (overridden only when `req.TickInterval` is set).

## Steps

1. Initialize `req.Scenario` and `req.Timeout` per leaf under `watchdog/`.
2. Run shared `Run` harness with fake clock and configured watch.
3. Leaf `Assert` checks `OnIdle` call count and HTTP side effects.

## Context

- P1 (`watchdog/`): `DevIdleWatch` core behavior only — no CLI flag, no `Serve()` wiring.
- P2 (`flag/`): `--dev-idle-life` CLI plumbing into `run.ServerOptions.DevIdleLife` via fake
  `StartServer`; no `Serve()` watchdog wiring, no SSE touch.
- P3 (`integration/`): `DevIdleWatch` wired into `server.Serve` / `server.ServeComponent` via
  in-process `server.ServeForTest` with injectable clock (no long real sleeps).
- P4 (`sse/`): nested DOCTEST — `sendSSEEvent` must call `DevIdleWatch.Touch()`; test route
  `/api/test-sse` via `ServeForTestOptions.EnableTestSSE`.
- P5 (`subprocess/`): nested DOCTEST — real `go build -o` binary, `exec.Command` subprocess,
  real wall-clock sleeps (1–3s); distinct from P3 in-process `ServeForTest`.
- Clock starts on **first touch**; no countdown before first request.

```go
import "time"

func Setup(t *testing.T, req *Request) error {
	if req.TickInterval == 0 {
		req.TickInterval = 50 * time.Millisecond
	}
	return nil
}
```