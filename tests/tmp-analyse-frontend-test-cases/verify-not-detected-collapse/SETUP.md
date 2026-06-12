## Preconditions
- The React app is compiled and served at `SERVER_URL`
- The tmp-analyse page groups not-detected software into a collapsed section

## Steps
1. Set req.ScriptFile to "not-detected-collapse.js"
2. The Run function executes playwright-debug with the script
3. The script navigates to /tmp-analyse and checks the collapse panel

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "not-detected-collapse.js"
	return nil
}
```
