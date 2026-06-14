## Preconditions
- The React app is compiled and served at `SERVER_URL`
- The Go card has a cleanup indicator that opens a popover with Go-specific cleanup suggestions

## Steps
1. Set req.ScriptFile to "cleanup-popover-go.js"
2. The Run function executes playwright-debug with the script
3. The script clicks the Go cleanup indicator and verifies popover content

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "cleanup-popover-go.js"
	return nil
}
```
