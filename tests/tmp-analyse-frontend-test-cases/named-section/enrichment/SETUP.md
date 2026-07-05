# Scenario

**Feature**: two-pass node_modules scan — rows before shared enrichment completes

```
node-modules-scan-btn -> pass1 named rows (shared 0 B) -> scan_complete done badge -> named_enriched merge -> final shared column
```

## Preconditions

- Rows must appear when `scan_complete` fires (done badge) before enrichment finishes.
- Shared column may show `0 B` at scan_complete and update after enrichment, or stay `0 B`.
- During enrichment, each pending row should show a per-row loading indicator in Shared.
- Shared totals should accumulate incrementally as `named_enriched` events arrive.
- First row should appear within 10s of scan click on a typical dev machine.

```go
func Setup(t *testing.T, req *Request) error {
	_ = req.ScriptFile
	return nil
}
```