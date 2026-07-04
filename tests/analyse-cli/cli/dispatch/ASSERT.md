## Expected

- No harness error from `run.RunWithOptions`.
- Web server hook is not invoked.
- Stdout contains a TSV summary line with `pathfmt.Short` of the fixture path.

## Exit Code

- 0 (via nil error from run)

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/dot-pkgs/go-pkgs/pathfmt"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ServerWasStarted {
		t.Fatal("web server must not start for analyse dispatch")
	}
	shortPath := pathfmt.Short(req.FixtureDir)
	if !strings.Contains(resp.Stdout, shortPath) {
		t.Fatalf("stdout missing short fixture path %q:\n%s", shortPath, resp.Stdout)
	}
}
```