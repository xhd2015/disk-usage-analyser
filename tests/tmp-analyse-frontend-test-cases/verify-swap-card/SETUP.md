# Scenario

**Feature**: swap-card UI verification

```
# React page renders tmp-analyse cards and scan UI
User -> /tmp-analyse -> playwright-debug checks data-testid elements
```

## Preconditions
- The React app is compiled and served at `SERVER_URL`
- Swap appears as a system (core) location card in the System Locations section

## Steps
1. Set req.ScriptFile to "swap-card.js"
2. The Run function executes playwright-debug with the script
3. The script navigates to /tmp-analyse and checks for swap card

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "swap-card.js"
	return nil
}
```
