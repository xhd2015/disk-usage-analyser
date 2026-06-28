# Scenario

**Feature**: Swap excluded from reclaimable total

```
# handler discovers locations, scans paths, streams SSE
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events
```

## Preconditions
- TmpLocation has a Reclaimable field (defaults to rebootSafe for existing locations)
- BuildSummary uses Reclaimable field instead of rebootSafe
- Swap location has Reclaimable=false and rebootSafe=true

## Steps
1. Create locations including a swap entry with Reclaimable=false, rebootSafe=true
2. Call BuildSummary
3. Verify swap's size is included in TotalSize but excluded from ReclaimableSize

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "swap-not-reclaimable"
	return nil
}
```
