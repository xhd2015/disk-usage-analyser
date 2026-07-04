# Scenario

**Leaf**: empty analysed root has no subdirectory rows and zero-byte summary

## Steps

1. Create an empty fixture directory.
2. Run `analyse.Analyse` on the fixture root.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "analyse"
	return nil
}
```