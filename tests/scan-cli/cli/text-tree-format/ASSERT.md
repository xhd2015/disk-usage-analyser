## Expected Output

```
PATH: __PATH__
TOTAL: 500 B
MIN: 1B
MAX-DEPTH: 3

.
├── big.txt    400B
└── small.txt  100B

```

## Expected

- Exit code 0.
- Summary block with `PATH`, `TOTAL`, `MIN`, `MAX-DEPTH`.
- Root line `.` followed by box-drawing branches: name first, then aligned size column (no brackets).
- Dirs show trailing `/`; files do not.
- Stdout ends with a trailing blank line.

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
TOTAL: 500 B
MIN: 1B
MAX-DEPTH: 3

.
├── big.txt    400B
└── small.txt  100B

`)
	if !strings.Contains(resp.Stdout, req.FixtureDir) {
		t.Fatalf("stdout PATH should include fixture dir %q:\n%s", req.FixtureDir, resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "big.txt/") {
		t.Fatal("file entries must not have trailing slash")
	}
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
