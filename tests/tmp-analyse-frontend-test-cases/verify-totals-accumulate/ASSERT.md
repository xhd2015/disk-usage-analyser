## Expected
- INIT total: 0 Bytes, INIT reclaimable: 0 Bytes
- MID total: NOT "0 Bytes" (totals are accumulating during scan)
- MID reclaimable: NOT "0 Bytes"
- FINAL total: NOT "0 Bytes"
- FINAL reclaimable: NOT "0 Bytes"

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	if !strings.Contains(resp.Output, "INIT total: 0 Bytes") {
		t.Fatal("expected INIT total: 0 Bytes")
	}
	if !strings.Contains(resp.Output, "INIT reclaimable: 0 Bytes") {
		t.Fatal("expected INIT reclaimable: 0 Bytes")
	}
	if !strings.Contains(resp.Output, "MID total:") {
		t.Fatal("expected MID total line")
	}
	if strings.Contains(resp.Output, "MID total: 0 Bytes") {
		t.Fatal("expected MID total != 0 Bytes (totals should accumulate during scan)")
	}
	if !strings.Contains(resp.Output, "MID reclaimable:") {
		t.Fatal("expected MID reclaimable line")
	}
	if !strings.Contains(resp.Output, "FINAL total:") {
		t.Fatal("expected FINAL total line")
	}
	if strings.Contains(resp.Output, "FINAL total: 0 Bytes") {
		t.Fatal("expected FINAL total != 0 Bytes")
	}
	if !strings.Contains(resp.Output, "FINAL reclaimable:") {
		t.Fatal("expected FINAL reclaimable line")
	}
}
```
