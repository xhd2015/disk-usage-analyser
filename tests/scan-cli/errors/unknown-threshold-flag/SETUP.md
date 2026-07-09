# Scenario

**Leaf**: removed `--threshold` flag is not accepted

## Steps

1. Run `RunCLI --threshold 1B <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--threshold", "1B", req.FixtureDir}
	return nil
}
```
