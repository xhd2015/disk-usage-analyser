# Scenario

**Leaf**: maxDepth 0 expands all branch levels (subject to threshold)

## Steps

1. Write `a/b/c/deep.bin` (42 bytes).
2. Run `usagescan.Scan` with `MaxDepth = 0` (unlimited).

```go
import "disk-usage-analyser/usagescan"

func Setup(t *testing.T, req *Request) error {
	mkdir(t, req.FixtureDir, "a/b/c")
	writeSizedFile(t, req.FixtureDir, "a/b/c/deep.bin", 42)
	req.ScanOpts = &usagescan.ScanOptions{
		Threshold: 0,
		MaxDepth:  0,
	}
	return nil
}
```