## Preconditions
- The React app is compiled and served at `SERVER_URL`
- The tmp-analyse page renders software-specific cards for detected tools

## Steps
1. Set req.ScriptFile to "software-cards-render.js"
2. The Run function executes playwright-debug with the script
3. The script navigates to /tmp-analyse and checks software card elements

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "software-cards-render.js"
	return nil
}
```
