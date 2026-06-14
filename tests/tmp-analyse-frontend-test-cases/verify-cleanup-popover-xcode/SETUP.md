## Preconditions
- The React app is compiled and served at `SERVER_URL`
- The Xcode card has a cleanup indicator that shows simulator-specific cleanup on click

## Steps
1. Set req.ScriptFile to "cleanup-popover-xcode.js"
2. The Run function executes playwright-debug with the script
3. The script clicks the Xcode cleanup indicator and verifies popover content

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "cleanup-popover-xcode.js"
	return nil
}
```
