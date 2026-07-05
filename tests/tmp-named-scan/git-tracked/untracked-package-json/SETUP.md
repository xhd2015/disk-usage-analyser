# Scenario

**Leaf**: on-disk but unindexed sibling `package.json` yields `gitTracked=false`

## Steps

1. Create git repo with `package.json` on disk (not `git add`) and `node_modules/`.
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := gitRepoWithNodeModules(t, req.HomeDir, "Projects/untracked-app")
	writeFile(t, app, "package.json", []byte("{\n  \"name\": \"untracked-app\"\n}\n"), 0644)
	return nil
}
```