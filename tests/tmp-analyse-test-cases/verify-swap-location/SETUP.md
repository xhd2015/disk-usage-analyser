# Scenario

**Feature**: Swap at /private/var/vm/ with category swap

```
# handler discovers locations, scans paths, streams SSE
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events
```

## Preconditions
- DiscoverLocations includes swap at /private/var/vm/ as a core system location
- Swap has category="swap", rebootSafe=true, and a field indicating it is not reclaimable

## Steps
1. Call DiscoverLocations with a test home directory
2. Find the swap location in the results
3. Verify its properties

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "swap-location"
	return nil
}
```
