# Scenario

**Leaf**: empty scan root has no items and zero total size

## Steps

1. Keep the fixture directory empty.
2. Run `usagescan.Scan` on the fixture root.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "scan"
	return nil
}
```