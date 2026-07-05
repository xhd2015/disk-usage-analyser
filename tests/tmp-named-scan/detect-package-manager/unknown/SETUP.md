# Scenario

**Leaf**: no recognized lockfiles yields `packageManager=unknown`

## Steps

1. Create git repo with `node_modules/` only (no lockfiles, no `.pnpm` store dir).
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	nodeModulesRepo(t, req.HomeDir, "Projects/plain-app")
	return nil
}
```