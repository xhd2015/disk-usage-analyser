# Scenario

**Feature**: Xcode card shows runtime section after scan when simulator runtimes available

```
# location event runtimeItems -> runtime-section on xcode card
CollectSimulatorRuntimeStats -> location(xcode) -> runtime-section / runtime-row-0
```

## Preconditions

- Xcode card detected.
- When simctl returns runtimes, runtime-section appears with at least runtime-row-0.

## Steps

1. Set `req.ScriptFile` to `xcode-runtime-section.js`.
2. Script completes scan and checks runtime-section or graceful absence.

## Context

- SKIP when no simctl / no runtimes (graceful absence).

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "xcode-runtime-section.js"
	return nil
}
```