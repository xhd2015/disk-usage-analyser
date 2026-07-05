# Scenario

**Leaf**: pass 1 `named` events expose `sharedSize=0` before enrichment

## Steps

1. Create git repo with `node_modules/`.
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	nodeModulesRepo(t, req.HomeDir, "Projects/fast-zero-app")
	return nil
}
```