# Scenario

**Leaf**: disabled watch (`Timeout=0`) never calls `OnIdle`

```
DevIdleWatch.Timeout = 0 -> disabled -> Start noop -> OnIdle never fires
```

## Steps

1. Set `Timeout` to 0.
2. Touch once and advance fake clock well past any would-be idle window.
3. Pump idle checks to prove background ticker does not invoke `OnIdle`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Scenario = "disabled-zero-timeout"
	req.Timeout = 0
	return nil
}
```