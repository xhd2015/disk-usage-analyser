# Scenario

**Leaf**: nested pnpm `node_modules/.pnpm/pkg@ver/node_modules` without sibling `package.json` yields `gitTracked=false`

## Steps

1. Create git repo with pnpm-style nested `node_modules` (no `package.json` beside nested hit).
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := gitRepoWithNodeModules(t, req.HomeDir, "Projects/nested-pnpm-git")
	writeFile(t, app, "pnpm-lock.yaml", []byte("lockfileVersion: 5.4\n"), 0644)
	pkg := "pkg@1.0.0"
	nested := filepath.Join("node_modules", ".pnpm", pkg, "node_modules", "pkg")
	writeFile(t, app, filepath.Join(nested, "index.js"), []byte("module.exports = {}\n"), 0644)
	return nil
}
```