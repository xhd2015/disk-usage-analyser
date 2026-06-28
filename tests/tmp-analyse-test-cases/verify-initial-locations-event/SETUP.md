# Scenario

**Feature**: First SSE event is locations with full array

```
# handler discovers locations, scans paths, streams SSE
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events
```

## Preconditions
- An SSE handler is registered at /api/tmp-analyse
- Before any scanning begins, the handler sends a "locations" event listing all discovered locations

## Steps
1. Create an httptest server with HandleTmpAnalyse
2. Send a GET request and start reading the SSE stream
3. Read only the first SSE event from the stream
4. Verify it is "event: locations" with a JSON array containing all discovered locations

## Context
- The first SSE event must be "locations" (plural), sent before any scanning happens
- The data field should be a JSON array of TmpLocation objects
- The array must include both Detected=true and Detected=false entries
- This allows the frontend to render all cards immediately before scanning starts

```go
import (
	"bufio"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"disk-usage-analyser/server"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "initial-locations-event"
	req.HomeDir = "/Users/testuser"
	return nil
}

```
