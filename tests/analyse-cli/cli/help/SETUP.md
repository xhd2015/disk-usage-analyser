# Scenario

**Leaf**: `-h` prints analyse usage and documented flags

## Steps

1. Run `RunCLI -h` without a path argument.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	req.Args = []string{"-h"}
	return nil
}
```