# Scenario

**Leaf**: Option B — tree section plus TOP N match list

## Steps

1. Write sample inspect JSON.
2. Run `RunCLI --inspect <file> --top 2`.

```go
func Setup(t *testing.T, req *Request) error {
	p := prepareSampleInspect(t, req)
	req.Args = []string{"--inspect", p.JSONPath, "--top", "2"}
	return nil
}
```
