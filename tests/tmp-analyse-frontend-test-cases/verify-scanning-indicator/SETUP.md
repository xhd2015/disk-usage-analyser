## Preconditions
- The Tmp Files Analyse page shows scanning indicators per card
- `scanning-badge` appears while a location is being scanned
- `done-badge` appears after a location scan completes

## Steps
1. Set req.ScriptFile to "scanning-indicator.js"
2. The script starts a scan and verifies badges appear and transition

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "scanning-indicator.js"
	return nil
}
```
