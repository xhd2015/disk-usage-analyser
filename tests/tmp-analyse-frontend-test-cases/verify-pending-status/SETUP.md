# Scenario

**Feature**: pending-status UI verification

```
# React page renders tmp-analyse cards and scan UI
User -> /tmp-analyse -> playwright-debug checks data-testid elements
```

## Preconditions
- The Tmp Files Analyse page is running
- When scan starts, all cards show `pending-badge` before receiving progress

## Steps
1. Set req.ScriptFile to "pending-status.js"
2. The script starts a scan and checks for pending badges during scan

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "pending-status.js"
	return nil
}
```
