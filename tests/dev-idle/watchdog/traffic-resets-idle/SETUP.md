# Scenario

**Leaf**: traffic within timeout window prevents idle shutdown

```
Touch T=0 -> advance 2s -> Touch T=2 -> advance 2s (T=4) -> gap since last touch < Timeout -> no OnIdle
```

## Steps

1. Set `Timeout` to 3s.
2. Touch at T=0, advance 2s, touch again at T=2.
3. Advance 2s more (T=4); elapsed since last touch is 2s (< 3s timeout).
4. Pump idle checks.

```go
import "time"

func Setup(t *testing.T, req *Request) error {
	req.Scenario = "traffic-resets-idle"
	req.Timeout = 3 * time.Second
	return nil
}
```