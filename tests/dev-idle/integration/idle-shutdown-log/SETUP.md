# Scenario

**Leaf**: idle shutdown writes `[dev] no requests for <duration>; shutting down` to stderr

```
GET /ping -> idle expiry -> shutdownDev() -> log idle message on stderr
```

## Steps

1. Set `DevIdleLife` to 1s.
2. Start `ServeForTest` with stderr captured in `req.Stderr`.
3. `GET /ping` once.
4. Advance fake clock 2s and pump idle checks.
5. Wait for port close; assert stderr contains shutdown message.

```go
import (
	"bytes"
	"time"
)

func Setup(t *testing.T, req *Request) error {
	req.Scenario = "idle-shutdown-log"
	req.DevIdleLife = 1 * time.Second
	if req.Stderr == nil {
		req.Stderr = &bytes.Buffer{}
	}
	return nil
}
```