# Scenario

**Feature**: scan-progress UI verification

```
# React page renders tmp-analyse cards and scan UI
User -> /tmp-analyse -> playwright-debug checks data-testid elements
```

## Preconditions
- The Tmp Files Analyse page serves real-time progress during scans
- Each card shows a `scanning-badge` while its location is being scanned
- Each card shows a `done-badge` after its location scan completes

## Steps
1. Set req.ScriptFile to "scan-progress.js"
2. The script starts a scan and checks that card sizes are updated from progress events

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "scan-progress.js"
	return nil
}
```
