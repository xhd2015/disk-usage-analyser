# Scenario

**Feature**: SSE handler emits all expected event types

```
# handler discovers locations, scans paths, streams SSE
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events
```

## Preconditions
- An SSE handler is registered at /api/tmp-analyse
- The handler streams location data as Server-Sent Events
- The handler first sends a "locations" event with all discovered locations before scanning

## Steps
1. Create an httptest server with the HandleTmpAnalyse handler
2. Make a GET request to /api/tmp-analyse
3. Read the response body as SSE events
4. Capture the SSE output in resp.SSEOutput

## Context
- SSE format: "event: <type>\ndata: <json>\n\n"
- Expected event types: "locations" (before scanning, full list), "location" (per completed location), "summary" (totals), "done" (completion)
- Event data must be valid JSON

```go
import (
	"io"
	"net/http"
	"net/http/httptest"

	"disk-usage-analyser/server"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "sse-format"
	req.HomeDir = "/Users/testuser"
	return nil
}

```
