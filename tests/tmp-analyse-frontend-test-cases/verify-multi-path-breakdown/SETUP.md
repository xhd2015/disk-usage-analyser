# Scenario

**Feature**: multi-path-breakdown UI verification

```
# React page renders tmp-analyse cards and scan UI
User -> /tmp-analyse -> playwright-debug checks data-testid elements
```

## Preconditions
- The React app is compiled and served at `SERVER_URL`
- Go and Xcode cards use unified breakdown-items to show sizes across their multiple scan paths

## Steps
1. Set req.ScriptFile to "multi-path-breakdown.js"
2. The Run function executes playwright-debug with the script
3. The script navigates to /tmp-analyse and looks for breakdown-items elements

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "multi-path-breakdown.js"
	return nil
}
```
