# Scenario

**Leaf**: long path cell shows prefix ellipsis with full suffix visible; full path in data-full-path

```
scan -> path cell -> visible truncated suffix -> data-full-path retains complete path
```

## Steps

1. Set req.ScriptFile to path-prefix-suffix-visible.js.
2. Run node_modules scan and inspect first long path cell.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "path-prefix-suffix-visible.js"
	return nil
}
```