# Scenario

**Feature**: cleanup-popover-npm UI verification

```
# React page renders tmp-analyse cards and scan UI
User -> /tmp-analyse -> playwright-debug checks data-testid elements
```

## Preconditions
- The React app is compiled and served at `SERVER_URL`
- The npm card has a cleanup indicator that opens a popover with cleanup suggestions on click

## Steps
1. Set req.ScriptFile to "cleanup-popover-npm.js"
2. The Run function executes playwright-debug with the script
3. The script clicks the npm cleanup indicator and verifies popover content

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "cleanup-popover-npm.js"
	return nil
}
```
