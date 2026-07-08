# Scenario

**Leaf**: default 1M threshold hides sub-1M nodes from tree display

## Steps

1. Write `large.bin` (2 MiB) and `small.bin` (512 KiB) at fixture root.
2. Run `usagescan.Scan` with default `ScanOptions` (threshold `1M`).

```go
func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "large.bin", 2<<20)
	writeSizedFile(t, req.FixtureDir, "small.bin", 512<<10)
	return nil
}
```