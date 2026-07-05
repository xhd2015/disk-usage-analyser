# Scenario

**Leaf**: project root without `package.json` yields `hasPackageJson=false`

## Steps

1. Create git repo with `node_modules/` only (no `package.json`, no lockfiles).
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	nodeModulesRepo(t, req.HomeDir, "Projects/no-pkgjson-app")
	return nil
}
```