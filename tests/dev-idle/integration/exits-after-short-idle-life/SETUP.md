# Scenario

**Leaf**: dev server closes listening port after short idle life expires

```
GET /ping -> DevIdleWatch.Touch -> idle clock starts
advance Now past DevIdleLife -> OnIdle -> shutdownDev() -> server.Close()
port no longer accepts connections
```

## Steps

1. Set `DevIdleLife` to 2s.
2. Start `ServeForTest` with `Dev: true`.
3. `GET /ping` once (response body `pong`).
4. Advance fake clock 3s and pump idle checks.
5. Wait until the port is closed.

```go
import "time"

func Setup(t *testing.T, req *Request) error {
	req.Scenario = "exits-after-short-idle-life"
	req.DevIdleLife = 2 * time.Second
	return nil
}
```