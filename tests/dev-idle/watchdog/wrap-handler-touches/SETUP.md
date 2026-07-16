# Scenario

**Leaf**: `Wrap` touches on HTTP request; idle only after request activity ages out

```
httptest GET -> DevIdleWatch.Wrap -> Touch -> handler 200
advance past Timeout without further traffic -> OnIdle once
```

## Steps

1. Set `Timeout` to 2s.
2. Issue one HTTP GET through wrapped handler.
3. Confirm `OnIdle` not fired immediately after request.
4. Advance fake clock 3s and pump idle checks.

```go
import "time"

func Setup(t *testing.T, req *Request) error {
	req.Scenario = "wrap-handler-touches"
	req.Timeout = 2 * time.Second
	return nil
}
```