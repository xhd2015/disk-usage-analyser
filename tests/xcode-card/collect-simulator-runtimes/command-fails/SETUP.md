# Scenario

**Feature**: simctl command failure returns empty slice without error

```
# CollectSimulatorRuntimeStats handles CLI errors gracefully
simctl error -> nil items, nil error
```

## Preconditions

- Mock runner returns error (simctl unavailable).

## Steps

1. Set `req.Op` to `collect-simulator-fails`.

## Context

- Same graceful pattern as Docker daemon unavailable.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "collect-simulator-fails"
	req.MockFail = true
	return nil
}
```