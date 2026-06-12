## Preconditions
- The React app is compiled and served at `SERVER_URL`
- Go and Xcode cards show a breakdown of sizes across their multiple scan paths

## Steps
1. Set req.ScriptFile to "multi-path-breakdown.js"
2. The Run function executes playwright-debug with the script
3. The script navigates to /tmp-analyse and looks for extra-path breakdown elements

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "multi-path-breakdown.js"
	return nil
}
```
