# Scenario

**Feature**: dry-run tracing without cache scans

```go
func Setup(t *testing.T, req *Request) error {
	if len(req.Args) == 0 {
		req.Args = []string{"--dry-run", inventoryPath(t)}
	}
	return nil
}
```