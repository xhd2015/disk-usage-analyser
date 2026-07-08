# Scenario

**Feature**: cleanup-popover-xcode UI verification

```
# React page renders tmp-analyse cards and scan UI
User -> /tmp-analyse -> playwright-debug checks data-testid elements
```

## Preconditions
- The React app is compiled and served at `SERVER_URL`
- The Xcode card has a cleanup indicator that shows expanded cleanup tips (DeviceSupport, runtime UUID delete, simulators) on click

## Steps
1. Set req.ScriptFile to "cleanup-popover-xcode.js"
2. The Run function executes playwright-debug with the script
3. The script clicks the Xcode cleanup indicator and verifies popover content

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "cleanup-popover-xcode.js"
	return nil
}
```
