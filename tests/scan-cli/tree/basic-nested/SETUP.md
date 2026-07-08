# Scenario

**Leaf**: nested tree children carry correct recursive sizes

## Steps

1. Write `alpha/beta/deep.bin` (200 bytes).
2. Write `gamma.txt` (300 bytes) at fixture root.
3. Run `usagescan.Scan` with default `ScanOptions`.

```go
import "disk-usage-analyser/usagescan"

func Setup(t *testing.T, req *Request) error {
	mkdir(t, req.FixtureDir, "alpha/beta")
	writeSizedFile(t, req.FixtureDir, "alpha/beta/deep.bin", 200)
	writeSizedFile(t, req.FixtureDir, "gamma.txt", 300)
	req.ScanOpts = &usagescan.ScanOptions{
		Threshold: 0,
		MaxDepth:  3,
	}
	return nil
}
```