# Scenario

**Leaf**: `--all-kinds` combined with `--kind` exits non-zero (mutually exclusive)

## Steps

1. Run `explain.RunCLI --all-kinds --kind xcode` (no PATH required to surface the flag conflict).

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--all-kinds", "--kind", "xcode"}
	return nil
}
```
