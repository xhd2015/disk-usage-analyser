# Scenario

**Leaf**: single touch then idle expiry invokes `OnIdle` once

```
Touch at T=0 -> advance Now past Timeout -> idle check -> OnIdle once
```

## Steps

1. Set `Timeout` to 2s.
2. Touch once at T=0.
3. Advance fake clock 3s (exceeds timeout).
4. Pump idle checks.

```go
import "time"

func Setup(t *testing.T, req *Request) error {
	req.Scenario = "exits-after-idle"
	req.Timeout = 2 * time.Second
	return nil
}
```