## Preconditions
- The React app is compiled and served at `SERVER_URL`
- npm card dynamically shows breakdown items when ~/.npm has subdirectories, or single path when empty

## Steps
1. Set req.ScriptFile to "npm-breakdown.js"
2. The Run function executes playwright-debug with the script
3. The script navigates to /tmp-analyse, starts a scan, and checks npm card structure

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "npm-breakdown.js"
	return nil
}
```
