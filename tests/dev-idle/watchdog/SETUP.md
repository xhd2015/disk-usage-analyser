# Scenario

**Decision**: DevIdleWatch core idle behavior (timeout, touch, wrap)

```
DevIdleWatch.Timeout -> enabled/disabled branch
DevIdleWatch.Touch -> reset lastActivity (clock starts on first touch)
DevIdleWatch.Wrap -> Touch on each HTTP request
DevIdleWatch.Start -> background idle checks -> OnIdle
```

## Preconditions

- All leaves exercise `server.DevIdleWatch` through the shared `Run` harness.
- Fake clock starts at a fixed UTC instant; advances are deterministic.

## Steps

1. Set `req.Scenario` to the leaf slug (`disabled-zero-timeout`, etc.).
2. Set `req.Timeout` per leaf table in requirement.

```go
import "time"

func Setup(t *testing.T, req *Request) error {
	if req.TickInterval <= 0 {
		req.TickInterval = 50 * time.Millisecond
	}
	return nil
}
```