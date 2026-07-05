# Scenario

**Leaf**: Git filter `no` hides git-tracked node_modules rows

## Steps

1. Set req.ScriptFile to column-filters-git-no-hides-tracked.js.
2. Scan, count tracked rows, select Git=No, verify tracked rows disappear.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "column-filters-git-no-hides-tracked.js"
	return nil
}
```