## Expected
- npm directory has at least 2 subdirectories (_cacache, _logs)
- Each subdirectory produces a breakdown item
- Breakdown items count matches the number of subdirectories

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.DetectedCount < 2 {
		t.Fatalf("expected at least 2 npm subdirectories, got %d", resp.DetectedCount)
	}

	if len(resp.ExtraPaths) < 2 {
		t.Fatalf("expected at least 2 ExtraPaths (npm subdirs), got %d", len(resp.ExtraPaths))
	}

	foundCacache := false
	foundLogs := false
	for _, p := range resp.ExtraPaths {
		if strings.Contains(p, "_cacache") {
			foundCacache = true
		}
		if strings.Contains(p, "_logs") {
			foundLogs = true
		}
	}
	if !foundCacache {
		t.Fatal("expected npm breakdown to include _cacache subdirectory")
	}
	if !foundLogs {
		t.Fatal("expected npm breakdown to include _logs subdirectory")
	}
}
```
