# Scenario

**Leaf**: `--help` prints usage and documented flags

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--help"}
	return nil
}
```