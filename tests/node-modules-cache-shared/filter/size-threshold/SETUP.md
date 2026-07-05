# Scenario

**Leaf**: `--size-threshold 10M` excludes sub-10M entries in dry-run summary

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--dry-run", "--size-threshold", "10M", inventoryPath(t)}
	return nil
}
```