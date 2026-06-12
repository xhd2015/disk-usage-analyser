## Preconditions
- The HandleTmpAnalyse handler checks runtime.GOOS before scanning
- On non-darwin platforms the handler must not scan any directories
- Instead it sends an unsupported_platform event and exits

## Steps
1. Note: this test cannot actually change runtime.GOOS at runtime
2. Instead, call DiscoverLocations which should always work regardless of OS
3. Verify the type system ensures TmpLocation has the fields needed
4. The actual platform gating is tested via integration with a simulated OS env

## Context
- This test verifies the response shape is correct and the function compiles
- Full platform gating verification requires running on a non-darwin machine or using build tags
- For now, we verify DiscoverLocations still works and the SSE flow includes initial locations event

```go
import (
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

	sseOutput := string(body)
	return &Response{SSEOutput: sseOutput}, nil
}
```
