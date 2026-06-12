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

	scanner := bufio.NewScanner(resp.Body)
	var sseOutput strings.Builder
	eventCount := 0
	var firstEventType string
	var firstEventData string

	for scanner.Scan() && eventCount < 1 {
		line := scanner.Text()
		sseOutput.WriteString(line + "\n")
		if strings.HasPrefix(line, "event: ") {
			firstEventType = strings.TrimPrefix(line, "event: ")
		}
		if strings.HasPrefix(line, "data: ") {
			firstEventData = strings.TrimPrefix(line, "data: ")
			eventCount++
		}
	}

	resp.Body.Close()

	var locations []map[string]interface{}
	if err := json.Unmarshal([]byte(firstEventData), &locations); err != nil {
		return &Response{
			SSEOutput: firstEventType + " | " + firstEventData,
		}, nil
	}

	detectedCount := 0
	notDetectedCount := 0
	for _, loc := range locations {
		if detected, ok := loc["detected"].(bool); ok && detected {
			detectedCount++
		} else {
			notDetectedCount++
		}
	}

	return &Response{
		SSEOutput:        firstEventType,
		DetectedCount:    detectedCount,
		NotDetectedCount: notDetectedCount,
	}, nil
}
```
