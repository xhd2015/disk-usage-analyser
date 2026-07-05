# Scenario

**Leaf**: `yarn.lock` in project root yields `packageManager=yarn`

## Steps

1. Create git repo with `yarn.lock` (no bun/pnpm/npm lockfiles).
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := nodeModulesRepo(t, req.HomeDir, "Projects/yarn-app")
	writeFile(t, app, "yarn.lock", []byte("# yarn lockfile v1\n"), 0644)
	return nil
}
```