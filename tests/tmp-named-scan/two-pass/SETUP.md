# Scenario

**Feature**: two-pass named scan — fast plain discovery then async shared enrichment

```
pass1 named (zero shared) -> enrichment worker may emit named_enriched early -> summary -> scan_complete -> done
```

## Preconditions

- `node_modules` scans run pass 2 enrichment; `vendor` scans skip it.
- Pass 1 `named` events must carry zeroed shared columns before any `named_enriched`.
- Event order: last `named` before `scan_complete`; `named_enriched` may interleave before `scan_complete`; `scan_complete` before `done`; every `named` path gets one `named_enriched` before `done`.
- Async `emitNamedDirSized` goroutines must not send enrichment jobs after `jobCh` is closed.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "named-scan"
	if req.Name == "" {
		req.Name = "node_modules"
	}
	return nil
}
```