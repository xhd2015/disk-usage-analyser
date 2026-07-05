# Scenario

**Feature**: node_modules table columns for package.json, package manager, and shared size

```
node-modules-scan -> SSE named hits -> table header + pkgjson/pkgmgr/shared cells per row
```

## Preconditions

- Table header row shows Path, package.json, PM, Shared, Size when a repo has hits.
- Each hit row exposes `node-modules-pkgjson-{rowKey}`, `node-modules-pkgmgr-{rowKey}`, and `node-modules-shared-{rowKey}`.
- PM cells use `white-space: nowrap` so labels like `unknown` do not wrap.

```go
func Setup(t *testing.T, req *Request) error {
	if req.ScriptFile == "" {
		req.ScriptFile = "table-columns-after-scan.js"
	}
	return nil
}
```