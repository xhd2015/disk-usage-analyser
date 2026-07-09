# Scenario

**Leaf**: `--suffix` produces tree + match section (Option B)

## Steps

1. Write sample inspect JSON.
2. Run `RunCLI --inspect <file> --suffix .bin`.

```go
func Setup(t *testing.T, req *Request) error {
	p := prepareSampleInspect(t, req)
	req.Args = []string{"--inspect", p.JSONPath, "--suffix", ".bin"}
	return nil
}
```
