# Scenario

**Decision**: SSE activity resets dev idle timer via `sendSSEEvent` → `DevIdleWatch.Touch`

```
# initial HTTP request touches via Wrap; each SSE frame re-touches
GET /api/test-sse -> DevIdleWatch.Wrap -> handler -> sendSSEEvent -> Touch

# stream keeps lastActivity fresh across interval > DevIdleLife
tick every 500ms for 3s with idle=2s -> no shutdown while streaming

# after stream ends, idle clock resumes from last touch
stream done -> advance 3s without traffic -> OnIdle -> server.Close()
```

## Preconditions

- Depends on P1–P3 (`DevIdleWatch`, `--dev-idle-life` flag, `ServeForTest` integration wiring).
- Implementer adds:
  - **`DevIdleWatch.Touch()`** call inside **`sendSSEEvent`** (or shared SSE write helper) using a
    watch ref reachable from SSE handlers (e.g. `serveRuntime.idleWatch` or context value).
  - Test route **`GET /api/test-sse`** registered when **`ServeForTestOptions.EnableTestSSE`**
    is true (or equivalent test hook). Query params: **`intervalMs`**, **`ticks`**.
  - Test SSE handler uses injectable **`Now`** from **`ServeForTest`** (poll `now()` between
    events; no multi-second `time.Sleep` when `Now` is set).
  - Events: **`tick`** (one per interval) then **`done`**.
- Harness uses **`ServeForTest`** with injectable clock; advances fake time between ticks and
  pumps **`DevIdleWatchForceCheckForTest`**.
- Nested root: `sse/DOCTEST.md` owns `Run`; leaves set `req.Scenario` to `<leaf-slug>`.

## Steps

1. Start `ServeForTest` with `Dev: true`, `EnableTestSSE: true`, and configured `DevIdleLife`.
2. Open long-lived `GET /api/test-sse` in a background goroutine.
3. Advance fake clock one `EventInterval` at a time; wait for each `tick` event.
4. Assert port state per leaf (still listening during stream vs closed after post-stream idle).

```go
import "time"

func Setup(t *testing.T, req *Request) error {
	if req.TickInterval <= 0 {
		req.TickInterval = 50 * time.Millisecond
	}
	if req.EventInterval <= 0 {
		req.EventInterval = 500 * time.Millisecond
	}
	if req.StreamTicks <= 0 {
		req.StreamTicks = 6
	}
	return nil
}
```