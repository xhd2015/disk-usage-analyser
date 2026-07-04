## Expected

- Exit code 0.
- Help text documents `analyse [DIR]`, `--header`, `--json`, and `-h,--help`.
- Help states table output always includes a header row.

## Exit Code

- 0

```go
import (
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
	assert.Output(t, resp.Stdout, ``+
`<contains>
disk-usage-analyser analyse
header is always printed
--header
--json
-h, --help
</contains>`)
}
```