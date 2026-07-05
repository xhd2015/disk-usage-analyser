# Scenario

**Leaf**: tooltip copy button copies full node_modules path to clipboard

```
hover path cell -> tooltip shows Copy button -> click -> clipboard has full path
```

## Steps

1. Set req.ScriptFile to path-tooltip-copy-button.js.
2. Run node_modules scan, hover first path, click copy, read clipboard.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "path-tooltip-copy-button.js"
	return nil
}
```