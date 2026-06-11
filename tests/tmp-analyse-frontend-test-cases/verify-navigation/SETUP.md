## Preconditions
- The React app is served at `SERVER_URL`
- The home page has a navigation bar with links

## Steps
1. Set req.ScriptFile to "navigation.js"
2. The script navigates to `/`, finds the "Tmp Files" nav link, clicks it, verifies URL changed to `/tmp-analyse`

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "navigation.js"
	return nil
}
```
