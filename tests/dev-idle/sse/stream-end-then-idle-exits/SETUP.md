# Scenario

**Leaf**: after SSE stream ends, idle shutdown fires once traffic stops

```
GET /api/test-sse -> ticks complete -> stream closes
no further Touch -> advance 3s past DevIdleLife=2s -> port closed
```

## Steps

1. Set `DevIdleLife` to 2s.
2. Start `ServeForTest` with test SSE route enabled.
3. Consume full SSE stream (6 ticks, 500ms interval).
4. Advance fake clock 3s with no further HTTP traffic.
5. Pump idle checks and wait until port is closed.

```go
import "time"

func Setup(t *testing.T, req *Request) error {
	req.Scenario = "stream-end-then-idle-exits"
	req.DevIdleLife = 2 * time.Second
	req.EventInterval = 500 * time.Millisecond
	req.StreamTicks = 6
	req.IdleAdvance = 3 * time.Second
	return nil
}
```