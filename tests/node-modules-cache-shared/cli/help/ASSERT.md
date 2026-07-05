## Expected

- Exit code 0.
- Help documents `--workers`, `--limit`, `--size-threshold`, `--dry-run`, `--verbose`.

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
node-modules-cache-shared
--workers
--limit
--size-threshold
--dry-run
--verbose
</contains>`)
}
```