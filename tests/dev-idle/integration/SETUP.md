# Scenario

**Decision**: DevIdleWatch wired into `server.Serve` / `server.ServeComponent` when `dev=true`

```
# run.Run passes DevIdleLife from ServerOptions into server
run.RunWithOptions -> ServerOptions.DevIdleLife -> server.Serve(dev=true)

# mux wrapped before ListenAndServe; idle shutdown shares signal path
HTTP request -> DevIdleWatch.Wrap -> mux -> handler
DevIdleWatch.Start -> idle? -> shutdownDev() -> server.Close() + cancel bun ctx

# disabled when DevIdleLife <= 0
DevIdleLife=0 -> no wrap, no idle shutdown, port stays open
```

## Preconditions

- Depends on P1 (`DevIdleWatch`) and P2 (`run.ServerOptions.DevIdleLife` flag plumbing).
- Target: `server.Serve` and `server.ServeComponent` with `Dev: true`.
- Harness uses **in-process** `server.ServeForTest` (implementer adds export test hook) with:
  - injectable `Now` clock (no long real sleeps),
  - `SkipFrontend: true` (no bun dev server on port 5173),
  - `NoBrowser: true` (`NO_BROWSER=1` equivalent),
  - `DevIdleWatchConfigureForTest` tick interval from `req.TickInterval` (default 50ms).
- `ServeForTest` returns the started `*http.Server` and `*DevIdleWatch` for clock pumping.
- Port chosen via `server.FindAvailablePort(18080, 50)` unless `req.Port` is preset.
- Nested root: `integration/DOCTEST.md` owns `Run`; leaves set `req.Scenario` to `<leaf-slug>`.

## Steps

1. Start `ServeForTest` in a background goroutine on an ephemeral port.
2. Wait until the port accepts TCP connections.
3. `GET /ping` once to start the idle clock.
4. Advance the fake clock and pump idle checks (or wait for auto-shutdown).
5. Assert port state and captured stderr per leaf.

```go
import "time"

func Setup(t *testing.T, req *Request) error {
	if req.TickInterval <= 0 {
		req.TickInterval = 50 * time.Millisecond
	}
	return nil
}
```