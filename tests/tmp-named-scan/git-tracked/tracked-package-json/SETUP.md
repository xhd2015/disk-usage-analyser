# Scenario

**Leaf**: committed sibling `package.json` yields `gitTracked=true` on enriched hit

## Steps

1. Create git repo with tracked `package.json` and `node_modules/`.
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := gitRepoWithNodeModules(t, req.HomeDir, "Projects/tracked-app")
	writeFile(t, app, "package.json", []byte("{\n  \"name\": \"tracked-app\"\n}\n"), 0644)
	runGit(t, app, "add", "package.json")
	runGit(t, app, "commit", "-m", "track package.json")
	return nil
}
```