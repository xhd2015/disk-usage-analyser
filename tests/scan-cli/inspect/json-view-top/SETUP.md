# Scenario

**Leaf**: inspect `--json --top` emits ViewResult with tree and matches

## Steps

1. Write sample inspect JSON.
2. Run `RunCLI --inspect <file> --json --top 2`.

```go
func Setup(t *testing.T, req *Request) error {
	p := prepareSampleInspect(t, req)
	req.Args = []string{"--inspect", p.JSONPath, "--json", "--top", "2"}
	return nil
}
```
