## Expected Output

```
PATH: __PATH__
TOTAL: 128 B
MIN: 1B
MAX-DEPTH: 3

.
└── data.bin  128B

```

## Expected

- Exit code 0.
- `PATH:` line includes the absolute fixture directory.
- Summary includes `TOTAL:`, `MIN:`, `MAX-DEPTH:`.
- Tree output shows `data.bin` on a box-drawing line with aligned size column (no brackets).

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
TOTAL: 128 B
MIN: 1B
MAX-DEPTH: 3

.
└── data.bin  128B

`)
	if !strings.Contains(resp.Stdout, req.FixtureDir) {
		t.Fatalf("stdout should show scanned path %q:\n%s", req.FixtureDir, resp.Stdout)
	}
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
