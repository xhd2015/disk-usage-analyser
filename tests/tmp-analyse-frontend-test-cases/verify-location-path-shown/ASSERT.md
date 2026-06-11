## Expected
- PATH_ELEMENTS > 0 (card-path elements exist)
- HAS_NONEMPTY_PATH: true (at least one path is non-empty)
- HAS_RECOGNIZABLE_PATH: true (at least one path contains known dir name)

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	if !strings.Contains(resp.Output, "PATH_ELEMENTS:") {
		t.Fatal("expected PATH_ELEMENTS line")
	}
	if strings.Contains(resp.Output, "PATH_ELEMENTS: 0") {
		t.Fatal("expected PATH_ELEMENTS > 0 (cards should show file paths)")
	}
	if !strings.Contains(resp.Output, "HAS_NONEMPTY_PATH: true") {
		t.Fatal("expected HAS_NONEMPTY_PATH: true (paths should not be empty)")
	}
	if !strings.Contains(resp.Output, "HAS_RECOGNIZABLE_PATH: true") {
		t.Fatal("expected HAS_RECOGNIZABLE_PATH: true (paths should contain recognizable dir names)")
	}
}
```
