# Scenario

**Leaf**: `pnpm-lock.yaml` in project root yields `packageManager=pnpm`

## Steps

1. Create git repo with `pnpm-lock.yaml` (no bun lockfiles).
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := nodeModulesRepo(t, req.HomeDir, "Projects/pnpm-app")
	writeFile(t, app, "pnpm-lock.yaml", []byte("lockfileVersion: 5.4\n"), 0644)
	return nil
}
```