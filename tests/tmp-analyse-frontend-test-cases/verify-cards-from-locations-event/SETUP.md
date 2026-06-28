# Scenario

**Feature**: cards-from-locations-event UI verification

```
# React page renders tmp-analyse cards and scan UI
User -> /tmp-analyse -> playwright-debug checks data-testid elements
```

## Preconditions
- The React app is compiled and served at `SERVER_URL`
- The tmp-analyse page receives a "locations" SSE event before any scanning
- Cards are rendered from this event data, not from hardcoded defaults

## Steps
1. Set req.ScriptFile to "cards-from-locations.js"
2. The Run function executes playwright-debug with the script
3. The script navigates to /tmp-analyse and verifies cards appear before scan

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "cards-from-locations.js"
	return nil
}
```
