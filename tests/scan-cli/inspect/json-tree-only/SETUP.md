# Scenario

**Leaf**: inspect `--json` without query flags emits ViewResult with tree

## Steps

1. Write sample inspect JSON.
2. Run `RunCLI --inspect <file> --json`.

```go
func Setup(t *testing.T, req *Request) error {
	p := prepareSampleInspect(t, req)
	req.Args = []string{"--inspect", p.JSONPath, "--json"}
	return nil
}
```
