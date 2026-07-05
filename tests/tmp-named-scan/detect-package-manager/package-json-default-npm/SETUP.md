# Scenario

**Leaf**: `package.json` only (no lockfile, no `packageManager` field) defaults to `packageManager=npm`

## Steps

1. Create git repo with `package.json` and `node_modules/` (no lockfiles).
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := nodeModulesRepo(t, req.HomeDir, "Projects/default-npm-app")
	writeFile(t, app, "package.json", []byte("{\n  \"name\": \"default-npm-app\"\n}\n"), 0644)
	return nil
}
```