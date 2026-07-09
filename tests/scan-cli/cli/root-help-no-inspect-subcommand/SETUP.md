# Scenario

**Leaf**: root help no longer lists `inspect` as a subcommand

## Steps

1. Call `run.RunWithOptions` with `-h` (root help, not `scan -h`).

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "dispatch"
	req.Args = []string{"-h"}
	return nil
}
```
