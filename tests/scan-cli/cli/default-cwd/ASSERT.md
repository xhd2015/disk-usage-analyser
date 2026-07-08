## Expected Output

```
PATH: __PATH__
TOTAL: 2M
THRESHOLD: 1M
MAX-DEPTH: 3

.
└── only.txt  2M

```

## Expected

- Exit code 0.
- Stdout contains summary lines: `PATH:`, `TOTAL:`, `THRESHOLD:`, `MAX-DEPTH:`.
- Stdout contains root tree line `.` and box-drawing child with name then aligned size (no brackets).
- Stdout ends with a trailing blank line after the last content line.

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
	assert.Output(t, resp.Stdout, `---
version: 2
__PATH__: type=string
---
PATH: __PATH__
TOTAL: 2M
THRESHOLD: 1M
MAX-DEPTH: 3

.
└── only.txt  2M

`)
	if !strings.Contains(resp.Stdout, req.FixtureDir) {
		t.Fatalf("stdout PATH should include fixture dir %q:\n%s", req.FixtureDir, resp.Stdout)
	}
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```