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
	req.HomeDir = "/Users/testuser"
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	handler := http.HandlerFunc(server.HandleTmpAnalyseLocations)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/tmp-analyse-locations")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	contentType := resp.Header.Get("Content-Type")

	var locations []server.TmpLocation
	if err := json.Unmarshal(body, &locations); err != nil {
		return nil, err
	}

	categoryCount := make(map[string]int)
	detectedCount := 0
	notDetectedCount := 0
	for _, loc := range locations {
		categoryCount[loc.Category]++
		if loc.Detected {
			detectedCount++
		} else {
			notDetectedCount++
		}
	}

	return &Response{
		Locations:        locations,
		CategoryCount:    categoryCount,
		DetectedCount:    detectedCount,
		NotDetectedCount: notDetectedCount,
		SSEOutput:        contentType,
	}, nil
}
```
