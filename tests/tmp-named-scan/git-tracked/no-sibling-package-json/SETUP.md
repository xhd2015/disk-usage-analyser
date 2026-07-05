# Scenario

**Leaf**: `node_modules/` without sibling `package.json` yields `gitTracked=false`

## Steps

1. Create git repo with `node_modules/` only (no `package.json` beside it).
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	gitRepoWithNodeModules(t, req.HomeDir, "Projects/no-pkgjson-app")
	return nil
}
```