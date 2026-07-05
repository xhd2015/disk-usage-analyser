## Expected

- HTTP status 200.
- Response `opened` is the absolute parent project directory (not `node_modules`).
- Captured AppleScript contains `create window`.
- Captured AppleScript references the parent project directory for `cd`.

## Errors

- No harness error is returned.

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d, want 200; body=%s", resp.StatusCode, resp.Body)
	}
	wantParent := projectParentAbs(req.HomeDir, "Projects/iterm-app")
	if resp.Opened != wantParent {
		t.Fatalf("opened = %q, want %q", resp.Opened, wantParent)
	}
	if resp.CapturedScript == "" {
		t.Fatal("expected captured AppleScript from KOOL_ITERM2_SCRIPT_OUT")
	}
	if !strings.Contains(resp.CapturedScript, "create window") {
		t.Fatalf("script missing create window:\n%s", resp.CapturedScript)
	}
	if !scriptContainsTargetDir(resp.CapturedScript, wantParent) {
		t.Fatalf("script missing parent dir %q:\n%s", wantParent, resp.CapturedScript)
	}
}
```