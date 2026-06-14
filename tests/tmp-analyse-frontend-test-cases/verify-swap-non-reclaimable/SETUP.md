## Preconditions
- The React app is compiled and served at `SERVER_URL`
- The swap card displays a non-reclaimable indicator distinct from the reclaimable safe-to-delete badge

## Steps
1. Set req.ScriptFile to "swap-non-reclaimable.js"
2. The Run function executes playwright-debug with the script
3. The script navigates to /tmp-analyse and checks the swap card for non-reclaimable indicator

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "swap-non-reclaimable.js"
	return nil
}
```
