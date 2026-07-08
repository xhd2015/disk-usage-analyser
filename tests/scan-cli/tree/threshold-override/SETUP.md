# Scenario

**Leaf**: custom 10M threshold filters display

## Steps

1. Write `medium.bin` (5 MiB) and `huge.bin` (15 MiB) at fixture root.
2. Run `usagescan.Scan` with `Threshold = 10 MiB`.

```go
import "disk-usage-analyser/usagescan"

func Setup(t *testing.T, req *Request) error {
	writeSizedFile(t, req.FixtureDir, "medium.bin", 5<<20)
	writeSizedFile(t, req.FixtureDir, "huge.bin", 15<<20)
	req.ScanOpts = &usagescan.ScanOptions{
		Threshold: 10 << 20,
		MaxDepth:  3,
	}
	return nil
}
```