# Scenario

**Leaf**: inspect `--max-depth 2` shows deeper tree nodes

## Steps

1. Write sample inspect JSON (depth ≥ 2).
2. Run `RunCLI --inspect <file> --max-depth 2`.

```go
func Setup(t *testing.T, req *Request) error {
	p := prepareSampleInspect(t, req)
	req.Args = []string{"--inspect", p.JSONPath, "--max-depth", "2"}
	return nil
}
```
