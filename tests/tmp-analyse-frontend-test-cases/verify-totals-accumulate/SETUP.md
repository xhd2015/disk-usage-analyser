## Preconditions
- The Tmp Files Analyse page shows summary totals that accumulate during scan
- MID_TOTAL is NOT "0 Bytes" (totals build up as files are scanned)
- After scan completes, FINAL_TOTAL equals sum of card sizes

## Steps
1. Set req.ScriptFile to "totals-accumulate.js"
2. The script starts a scan and checks that total/reclaimable accumulate mid-scan

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "totals-accumulate.js"
	return nil
}
```
