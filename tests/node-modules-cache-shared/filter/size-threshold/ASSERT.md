## Expected

- Exit code 0.
- Dry-run summary reports `would_scan=2` (large + huge only).
- Stderr mentions `below_threshold=1`.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d (err=%v)", resp.ExitCode, resp.Err)
	}
	assert.Output(t, resp.Stderr, ``+
`<contains>
below_threshold=1
would_scan=2
</contains>`)
	if strings.Contains(resp.Stderr, "small/node_modules") && strings.Contains(resp.Stderr, "would scan") {
		t.Fatalf("small entry should be filtered out: %q", resp.Stderr)
	}
}
```