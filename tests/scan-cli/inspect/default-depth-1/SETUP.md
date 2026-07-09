# Scenario

**Leaf**: inspect defaults to max-depth 1, min 0, with SOURCE line

## Steps

1. Write fixed TreeResult JSON with depth ≥ 2 children (`sampleInspectTree`).
2. Run `RunCLI --inspect <file>` with no extra flags.

```go
func Setup(t *testing.T, req *Request) error {
	p := prepareSampleInspect(t, req)
	req.Args = []string{"--inspect", p.JSONPath}
	return nil
}
```
