# Scenario

**Feature**: REST endpoint returns locations JSON without scanning

```
# handler discovers locations, scans paths, streams SSE
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events
```

## Preconditions
- A REST endpoint GET /api/tmp-analyse-locations is registered
- The endpoint returns all discovered locations as a JSON array
- No scanning occurs — the response is immediate JSON, not an SSE stream

## Steps
1. Set req.HomeDir to "/Users/testuser"
2. Create an httptest server with HandleTmpAnalyseLocations
3. Make a GET request and parse the JSON response
4. Verify the response is a valid JSON array with at least 22 locations
5. Verify each location has Detected field, categories, labels, and paths

## Context
- This endpoint is for initial page load — frontend fetches it on mount to show all cards immediately
- Separate from the SSE /api/tmp-analyse endpoint which handles scanning
- Response Content-Type must be application/json, not text/event-stream

```go
import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"io"

	"disk-usage-analyser/server"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "locations-rest-endpoint"
	req.HomeDir = "/Users/testuser"
	return nil
}

```
