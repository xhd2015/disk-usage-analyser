## Expected
- Home page URL is loaded
- `NAV_LINK Tmp Files: FOUND` appears (not MISSING)
- After click, `URL after click` contains `/tmp-analyse`
- `URL tmp-analyse: REACHED` appears
- `ELEM page-heading` contains "Tmp Files Analyse"

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	if !strings.Contains(resp.Output, "NAV_LINK Tmp Files: FOUND") {
		t.Fatal("expected 'Tmp Files' nav link to be found")
	}
	if strings.Contains(resp.Output, "NAV_LINK Tmp Files: MISSING") {
		t.Fatal("nav link 'Tmp Files' is missing")
	}
	if !strings.Contains(resp.Output, "URL tmp-analyse: REACHED") {
		t.Fatal("expected URL to reach /tmp-analyse after click")
	}
	if !strings.Contains(resp.Output, "ELEM page-heading:") {
		t.Fatal("expected page-heading element on /tmp-analyse")
	}
	if !strings.Contains(resp.Output, "Tmp Files Analyse") {
		t.Fatal("expected page-heading to contain 'Tmp Files Analyse'")
	}
}
```
