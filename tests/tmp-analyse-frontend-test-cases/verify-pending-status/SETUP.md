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
