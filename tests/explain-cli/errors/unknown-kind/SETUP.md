# Scenario

**Leaf**: unknown `--kind` value exits non-zero

## Steps

1. Run `explain.RunCLI --kind not-a-kind` (no PATH needed to surface kind error).

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--kind", "not-a-kind"}
	return nil
}
```
