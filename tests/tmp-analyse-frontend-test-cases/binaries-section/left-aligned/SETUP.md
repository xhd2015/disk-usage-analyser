# Scenario

**Leaf**: binaries tree content is flush left within the card

## Steps

1. Set req.ScriptFile to binaries-left-aligned.js.
2. Run binaries scan, wait for done, check computed text-align on binaries-tree.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "binaries-left-aligned.js"
	return nil
}
```