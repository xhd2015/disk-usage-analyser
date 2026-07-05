# Scenario

**Leaf**: nested pnpm `node_modules/.pnpm/pkg@ver/node_modules` inherits PM from app-root lockfile

Reproduces production path shape:

```
frontend/node_modules/.pnpm/@babel+helpers@7.29.7/node_modules
```

`DetectPackageManager` currently uses only `parent(node_modules)` (the `.pnpm/pkg@ver`
folder), so it returns `unknown` even when `pnpm-lock.yaml` exists at the app root.

## Steps

1. Create git repo with `pnpm-lock.yaml` at app root.
2. Add pnpm-style nested `node_modules` tree (no lockfiles beside the nested dir).
3. Run `named-scan`.

```go
import "fmt"

func nestedPnpmNodeModulesRepo(t *testing.T, home, rel string) string {
	t.Helper()
	app := repo(t, home, rel)
	writeFile(t, app, "pnpm-lock.yaml", []byte("lockfileVersion: 5.4\n"), 0644)
	pkg := "@babel+helpers@7.29.7"
	nested := filepath.Join("node_modules", ".pnpm", pkg, "node_modules", "@babel", "helpers")
	writeFile(t, app, filepath.Join(nested, "index.js"), []byte("module.exports = {}\n"), 0644)
	return app
}

func Setup(t *testing.T, req *Request) error {
	nestedPnpmNodeModulesRepo(t, req.HomeDir, "Projects/nested-pnpm-app")
	return nil
}
```