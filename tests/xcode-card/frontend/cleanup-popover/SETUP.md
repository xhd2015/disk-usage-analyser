# Scenario

**Feature**: Xcode cleanup popover covers expanded paths and runtime delete guidance

```
# cleanup indicator -> cleanup-popover-xcode with DeviceSupport, runtime UUID delete, DerivedData
User clicks cleanup-indicator -> popover lists cleanup commands
```

## Preconditions

- Xcode card detected with cleanup indicator.

## Steps

1. Set `req.ScriptFile` to `cleanup-popover-xcode.js`.
2. Script opens popover and checks expanded cleanup content.

## Context

- Runtime delete must reference UUID pattern, not SimRuntime bundle IDs.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "cleanup-popover-xcode.js"
	return nil
}
```