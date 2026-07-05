# Scenario

**Leaf**: multi-hit pipelined scan completes without SSE write races or server errors

## Steps

1. Create three git repos each with `node_modules/` (multi-hit fixture).
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	nodeModulesRepo(t, req.HomeDir, "Projects/concurrent-app-a")
	nodeModulesRepo(t, req.HomeDir, "Projects/concurrent-app-b")
	nodeModulesRepo(t, req.HomeDir, "Projects/concurrent-app-c")
	return nil
}
```