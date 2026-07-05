# Scenario

**Leaf**: `--dry-run` traces planned scans to stderr and emits no stdout JSONL

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--dry-run", inventoryPath(t)}
	return nil
}
```