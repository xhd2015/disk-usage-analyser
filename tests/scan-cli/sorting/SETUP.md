# Scenario

**Feature**: tree child ordering matches web UI sort rules

```
Scan -> build tree -> sort each node's children by size desc -> dirs before files on equal size
```

## Preconditions

- Fixture entries have deterministic byte sizes.
- Sorting applies to each node's `children` slice (root children asserted here).

## Context

- Tie-break: when two children share the same `size`, directories precede files.

```go
import "disk-usage-analyser/usagescan"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "scan"
	req.ScanOpts = &usagescan.ScanOptions{
		Threshold: 0,
		MaxDepth:  3,
	}
	return nil
}
```