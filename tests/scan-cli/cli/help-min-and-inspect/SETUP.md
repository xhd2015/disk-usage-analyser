# Scenario

**Leaf**: `-h` documents inspect and query flags on scan

## Steps

1. Run `RunCLI -h`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"-h"}
	return nil
}
```
