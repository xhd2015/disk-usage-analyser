# Scenario

**Leaf**: maxDepth 2 truncates branches; ancestor sizes include deeper bytes

## Steps

1. Write `d1/d2/leaf.bin` (1000 bytes).
2. Run `usagescan.Scan` with `MaxDepth = 2`.

```go
import "disk-usage-analyser/usagescan"

func Setup(t *testing.T, req *Request) error {
	mkdir(t, req.FixtureDir, "d1/d2")
	writeSizedFile(t, req.FixtureDir, "d1/d2/leaf.bin", 1000)
	req.ScanOpts = &usagescan.ScanOptions{
		Min: 0,
		MaxDepth:  2,
	}
	return nil
}
```