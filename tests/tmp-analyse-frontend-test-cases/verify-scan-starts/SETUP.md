# Scenario

**Feature**: scan-starts UI verification

```
# React page renders tmp-analyse cards and scan UI
User -> /tmp-analyse -> playwright-debug checks data-testid elements
```

## Preconditions
- The Tmp Files Analyse page is implemented and running
- The backend SSE endpoint `/api/tmp-analyse` is functional (or will be in implementation phase)

## Steps
1. Set req.ScriptFile to "scan-starts.js"
2. The script navigates to /tmp-analyse, clicks Start Scan, monitors SSE events, checks card sizes update and button toggles

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "scan-starts.js"
	return nil
}
```
