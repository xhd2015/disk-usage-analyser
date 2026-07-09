# Scenario

**Leaf**: invalid `--min` value exits non-zero

## Steps

1. Run `RunCLI --min foo <fixture>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--min", "foo", req.FixtureDir}
	return nil
}
```
