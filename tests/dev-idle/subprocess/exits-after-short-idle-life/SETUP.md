# Scenario

**Leaf**: compiled CLI subprocess closes port after short idle life expires

```
exec.Command(bin, --dev, --dev-idle-life 2s) -> listen on ephemeral port
GET /ping -> DevIdleWatch.Touch -> idle clock starts
sleep 3s real time -> OnIdle -> shutdownDev() -> process exit, port closed
```

## Steps

1. Reuse session-built binary from parent `Setup`.
2. Start subprocess with `--dev --dev-idle-life 2s` and `NO_BROWSER=1`.
3. Parse port from stdout; `GET /ping` once (body `pong`).
4. Sleep 12s wall clock (past 2s idle life plus default 10s watchdog tick interval).
5. Assert TCP dial fails and subprocess has exited.

```go
import "time"

func Setup(t *testing.T, req *Request) error {
	req.Scenario = "exits-after-short-idle-life"
	req.DevIdleLife = 2 * time.Second
	req.IdleSleep = 12 * time.Second
	return nil
}
```