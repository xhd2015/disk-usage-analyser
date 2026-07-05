# Scenario

**Feature**: inventory filtering by size threshold and limit

```go
func Setup(t *testing.T, req *Request) error {
	if len(req.Args) == 0 {
		req.Args = []string{"--dry-run", "--size-threshold", "10M", inventoryPath(t)}
	}
	return nil
}
```