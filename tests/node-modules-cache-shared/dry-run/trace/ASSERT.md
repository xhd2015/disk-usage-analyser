## Expected

- Exit code 0.
- Stdout is empty.
- Stderr contains dry-run summary and `would scan` lines.

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
	if strings.TrimSpace(resp.Stdout) != "" {
		t.Fatalf("expected empty stdout, got %q", resp.Stdout)
	}
	assert.Output(t, resp.Stderr, ``+
`<contains>
dry-run summary
would_scan=3
dry-run would scan
</contains>`)
}
```