# Scenario

**Leaf**: maxDepth 1 shows only immediate children of the scan root

## Steps

1. Write `dir/nested.bin` (80 bytes) and `top.bin` (20 bytes) at fixture root.
2. Run `usagescan.Scan` with `MaxDepth = 1`.

```go
import "disk-usage-analyser/usagescan"

func Setup(t *testing.T, req *Request) error {
	mkdir(t, req.FixtureDir, "dir")
	writeSizedFile(t, req.FixtureDir, "dir/nested.bin", 80)
	writeSizedFile(t, req.FixtureDir, "top.bin", 20)
	req.ScanOpts = &usagescan.ScanOptions{
		Min: 0,
		MaxDepth:  1,
	}
	return nil
}
```