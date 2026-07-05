# Scenario

**Feature**: CLI help surface

```go
func Setup(t *testing.T, req *Request) error {
	if len(req.Args) == 0 {
		req.Args = []string{"--help"}
	}
	return nil
}
```