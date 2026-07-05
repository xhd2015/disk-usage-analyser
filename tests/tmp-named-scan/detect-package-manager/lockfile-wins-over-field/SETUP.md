# Scenario

**Leaf**: lockfile marker beats Corepack `packageManager` field when both are present

## Steps

1. Create git repo with `package-lock.json`, `package.json` containing `"packageManager": "pnpm@9.0.0"`, and `node_modules/`.
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := nodeModulesRepo(t, req.HomeDir, "Projects/lockfile-wins-app")
	writeFile(t, app, "package-lock.json", []byte("{\n  \"name\": \"lockfile-wins-app\"\n}\n"), 0644)
	writeFile(t, app, "package.json", []byte("{\n  \"name\": \"lockfile-wins-app\",\n  \"packageManager\": \"pnpm@9.0.0\"\n}\n"), 0644)
	return nil
}
```