## Expected Output

```
PATH: __PATH__
TOTAL: 350 B
THRESHOLD: 1B
MAX-DEPTH: 3

.
├── long-name-file.txt  200B
├── a.txt               100B
└── shortdir/           50B
    └── b.txt           50B

```

## Expected

- Exit code 0.
- Root children sorted by size descending: `long-name-file.txt`, `a.txt`, `shortdir/`.
- Every visible tree row places the size at the same column (padding after the name).
- Directory `shortdir/` shows trailing `/`; nested `b.txt` does not.
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
TOTAL: 350 B
THRESHOLD: 1B
MAX-DEPTH: 3

.
├── long-name-file.txt  200B
├── a.txt               100B
└── shortdir/           50B
    └── b.txt           50B

`)
	if !strings.Contains(resp.Stdout, req.FixtureDir) {
		t.Fatalf("stdout PATH should include fixture dir %q:\n%s", req.FixtureDir, resp.Stdout)
	}
	assertTreeSizeColumnAligned(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```