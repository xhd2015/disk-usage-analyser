# Scenario

**Leaf**: Corepack `packageManager` field `pnpm@9.0.0` yields `packageManager=pnpm` without lockfile

## Steps

1. Create git repo with `package.json` containing `"packageManager": "pnpm@9.0.0"` and `node_modules/` (no lockfiles).
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := nodeModulesRepo(t, req.HomeDir, "Projects/corepack-pnpm-app")
	writeFile(t, app, "package.json", []byte("{\n  \"name\": \"corepack-pnpm-app\",\n  \"packageManager\": \"pnpm@9.0.0\"\n}\n"), 0644)
	return nil
}
```