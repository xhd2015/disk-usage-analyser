# Scenario

**Feature**: stop-scan UI verification

```
# React page renders tmp-analyse cards and scan UI
User -> /tmp-analyse -> playwright-debug checks data-testid elements
```

## Preconditions
- The Tmp Files Analyse page is running with functional SSE backend
- A scan can be started and stopped

## Steps
1. Set req.ScriptFile to "stop-scan.js"
2. The script starts a scan, then clicks stop, verifies button reverts and SSE stream closes

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "stop-scan.js"
	return nil
}
```
