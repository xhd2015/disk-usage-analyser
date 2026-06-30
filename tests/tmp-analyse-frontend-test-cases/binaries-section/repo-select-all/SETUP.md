# Scenario

**Leaf**: repo parent checkbox selects all binary leaves in that repo

## Steps

1. Set req.ScriptFile to binaries-repo-select-all.js.
2. Scan, click repo checkbox, verify all child checkboxes selected and total updates.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "binaries-repo-select-all.js"
	return nil
}
```