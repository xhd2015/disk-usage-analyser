# Scenario

**Leaf**: active SSE stream prevents idle shutdown even when elapsed time exceeds `DevIdleLife`

```
GET /api/test-sse -> tick every 500ms
each sendSSEEvent -> Touch -> lastActivity reset
advance 3s total with idle=2s -> port still listening
```

## Steps

1. Set `DevIdleLife` to 2s.
2. Start `ServeForTest` with test SSE route enabled.
3. Open SSE stream: 6 ticks at 500ms (3s total).
4. Advance fake clock per tick and pump idle checks after each advance.
5. Assert port still accepts TCP connections after the stream completes.

```go
import "time"

func Setup(t *testing.T, req *Request) error {
	req.Scenario = "stream-prevents-idle-exit"
	req.DevIdleLife = 2 * time.Second
	req.EventInterval = 500 * time.Millisecond
	req.StreamTicks = 6
	return nil
}
```