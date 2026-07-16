## Expected Output

```
PATH: __PATH__
TOTAL: 512 B
MIN: 1B
MAX-DEPTH: 3

.
└── note.txt  512B

```

## Expected

- No harness error from `run.RunWithOptions`.
- Web server hook is not invoked.
- Stdout contains summary `PATH:` with the fixture directory and `note.txt` on an aligned-size tree line.

## Exit Code

- 0 (via nil error from run)

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
	if resp.ServerWasStarted {
		t.Fatal("web server must not start for scan dispatch")
	}
	assert.Output(t, resp.Stdout, `---
version: 3
__PATH__: type=string
---
PATH: __PATH__
TOTAL: 512 B
MIN: 1B
MAX-DEPTH: 3

\.
└── note\.txt  512B

`)
	if !strings.Contains(resp.Stdout, req.FixtureDir) {
		t.Fatalf("stdout missing fixture path %q:\n%s", req.FixtureDir, resp.Stdout)
	}
}
```
