# Scenario

**Feature**: BuildSummary computes total + reclaimable

```
# handler discovers locations, scans paths, streams SSE
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events
```

## Preconditions
- Four pre-scanned TmpLocation entries are provided:
  - User Trash: size=1000, rebootSafe=true
  - User Caches: size=500, rebootSafe=true
  - System Temp: size=2000, rebootSafe=false
  - User Logs: size=300, rebootSafe=true

## Steps
1. Create four TmpLocation entries with known sizes and rebootSafe flags
2. Call BuildSummary with these locations
3. TotalSize should be 1000+500+2000+300 = 3800
4. ReclaimableSize should be 1000+500+300 = 1800 (only rebootSafe=true items)

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "summary-totals"
	return nil
}
```
