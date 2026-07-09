# Scenario

**Leaf**: root totalSize counts sub-min bytes omitted from tree display

## Steps

1. Write `visible.bin` (2 MiB) and `hidden.bin` (500 bytes) at fixture root.
2. Run `usagescan.Scan` with default min `1M`.

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "visible.bin", 2<<20)
	writeSizedFile(t, req.FixtureDir, "hidden.bin", 500)
	return nil
}
```