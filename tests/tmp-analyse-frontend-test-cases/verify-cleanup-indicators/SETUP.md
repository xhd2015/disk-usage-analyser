# Scenario

**Feature**: cleanup-indicators UI verification

```
# React page renders tmp-analyse cards and scan UI
User -> /tmp-analyse -> playwright-debug checks data-testid elements
```

## Preconditions
- The React app is compiled and served at `SERVER_URL`
- Each section card has a cleanup indicator icon that can be clicked

## Steps
1. Set req.ScriptFile to "cleanup-indicators.js"
2. The Run function executes playwright-debug with the script
3. The script navigates to /tmp-analyse and checks for cleanup indicator elements on each card

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "cleanup-indicators.js"
	return nil
}
```
