# Scenario

**Leaf**: SSE event order for two-pass node_modules scan with multiple hits

## Steps

1. Create two git repos each with `node_modules/`.
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	nodeModulesRepo(t, req.HomeDir, "Projects/order-app-a")
	nodeModulesRepo(t, req.HomeDir, "Projects/order-app-b")
	return nil
}
```