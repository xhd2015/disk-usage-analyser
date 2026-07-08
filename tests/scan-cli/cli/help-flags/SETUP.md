# Scenario

**Leaf**: `-h` documents threshold and max-depth flags

## Steps

1. Run `RunCLI -h` without a path argument.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"-h"}
	return nil
}
```