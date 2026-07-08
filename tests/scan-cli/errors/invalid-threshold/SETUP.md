# Scenario

**Leaf**: invalid `--threshold` value exits non-zero

## Steps

1. Run `RunCLI --threshold foo <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--threshold", "foo", req.FixtureDir}
	return nil
}
```