# Scenario

**Leaf**: `package.json` in project root yields `hasPackageJson=true` and `packageManager=unknown`

## Steps

1. Create git repo with `package.json` and `node_modules/` (no lockfiles).
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := nodeModulesRepo(t, req.HomeDir, "Projects/pkgjson-app")
	writeFile(t, app, "package.json", []byte("{\n  \"name\": \"pkgjson-app\"\n}\n"), 0644)
	return nil
}
```