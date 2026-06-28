# Scenario

**Feature**: self-contained-server UI verification

```
# React page renders tmp-analyse cards and scan UI
User -> /tmp-analyse -> playwright-debug checks data-testid elements
```

## Preconditions
- The test harness (root SETUP.md) starts its own server, not requiring an externally running server
- The server is started in --dev mode, auto-starts Vite dev server if needed
- The server port is discovered from stdout and passed to playwright via SERVER_URL

## Steps
1. Set req.ScriptFile to "self-contained-server.js"
2. The script verifies the server /ping endpoint responds
3. Then navigates to /tmp-analyse and checks basic page rendering

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "self-contained-server.js"
	return nil
}
```
