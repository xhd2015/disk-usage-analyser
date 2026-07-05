# Scenario

**Leaf**: `--limit 1` caps dry-run scheduled scans after filtering

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--dry-run", "--size-threshold", "10M", "--limit", "1", inventoryPath(t)}
	return nil
}
```