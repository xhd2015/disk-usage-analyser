## Preconditions
- Go and Xcode cards show breakdown items in table-like rows
- Each row has data-testid="breakdown-row-{idx}" and uses flexbox layout

## Steps
1. Set req.ScriptFile to "breakdown-table-layout.js"
2. The script navigates to /tmp-analyse and checks the row layout structure

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "breakdown-table-layout.js"
	return nil
}
```
