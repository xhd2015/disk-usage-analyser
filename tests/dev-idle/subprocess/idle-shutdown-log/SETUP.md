# Scenario

**Leaf**: subprocess idle shutdown writes `[dev] no requests for <duration>; shutting down` to stderr

```
exec.Command(bin, --dev, --dev-idle-life 1s) -> stderr captured
GET /ping -> idle expiry -> shutdownDev() -> log idle message on stderr
```

## Steps

1. Start subprocess with `--dev --dev-idle-life 1s` and `NO_BROWSER=1`.
2. Parse port; `GET /ping` once.
3. Sleep 12s wall clock (past 1s idle life plus default 10s watchdog tick interval).
4. Wait for process exit; assert stderr contains shutdown log line.

```go
import "time"

func Setup(t *testing.T, req *Request) error {
	req.Scenario = "idle-shutdown-log"
	req.DevIdleLife = 1 * time.Second
	req.IdleSleep = 12 * time.Second
	return nil
}
```