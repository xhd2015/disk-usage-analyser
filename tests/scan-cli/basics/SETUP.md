# Scenario

**Feature**: fundamental recursive tree scan without links

```
Scan(path, ScanOptions) -> walk tree -> TreeResult{tree, totalSize}
```

## Preconditions

- Fixture trees contain only regular files and directories.
- No symlinks or special filesystem features in this branch.

## Context

- Default `ScanOptions`: threshold `1M`, maxDepth `3`.
- Nested content rolls up into parent directory sizes.
- Root-level regular files appear as immediate `tree.children`.

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