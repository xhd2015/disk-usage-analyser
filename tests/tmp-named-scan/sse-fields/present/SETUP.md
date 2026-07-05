# Scenario

**Leaf**: SSE `named_enriched` JSON includes all package manager and shared size fields

## Steps

1. Create git repo with `package-lock.json` and `node_modules/`.
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	app := nodeModulesRepo(t, req.HomeDir, "Projects/sse-app")
	writeFile(t, app, "package-lock.json", []byte("{\n  \"name\": \"sse-app\"\n}\n"), 0644)
	return nil
}
```