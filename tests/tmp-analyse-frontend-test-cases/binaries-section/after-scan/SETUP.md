# Scenario

**Leaf**: binaries scan populates repo-grouped tree with kind badges

## Steps

1. Set req.ScriptFile to binaries-after-scan.js.
2. Run binaries scan and verify tree structure.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "binaries-after-scan.js"
	return nil
}
```