# Scenario

**Leaf**: `package-lock.json` in project root yields `packageManager=npm`

## Steps

1. Create git repo with `package-lock.json` (no bun or pnpm lockfiles).
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := nodeModulesRepo(t, req.HomeDir, "Projects/npm-app")
	writeFile(t, app, "package-lock.json", []byte("{\n  \"name\": \"npm-app\"\n}\n"), 0644)
	return nil
}
```