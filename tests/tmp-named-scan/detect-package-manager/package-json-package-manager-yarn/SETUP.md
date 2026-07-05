# Scenario

**Leaf**: Corepack `packageManager` field `yarn@4.0.0` yields `packageManager=yarn` without lockfile

## Steps

1. Create git repo with `package.json` containing `"packageManager": "yarn@4.0.0"` and `node_modules/` (no lockfiles).
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := nodeModulesRepo(t, req.HomeDir, "Projects/corepack-yarn-app")
	writeFile(t, app, "package.json", []byte("{\n  \"name\": \"corepack-yarn-app\",\n  \"packageManager\": \"yarn@4.0.0\"\n}\n"), 0644)
	return nil
}
```