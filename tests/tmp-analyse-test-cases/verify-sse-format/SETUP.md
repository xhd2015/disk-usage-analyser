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
	req.HomeDir = "/Users/testuser"
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	handler := http.HandlerFunc(server.HandleTmpAnalyse)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	httpReq, err := http.NewRequest("GET", srv.URL+"/api/tmp-analyse", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{SSEOutput: string(body)}, nil
}
```
