## Expected

- After scan completes, if `DOCKER_RUNTIME_SECTION: present` then runtime-row-0, runtime-label-0, runtime-count-0, and runtime-size-0 must be present.
- If `DOCKER_RUNTIME_SECTION: absent`, `DOCKER_RUNTIME_GRACEFUL: true` is logged (CLI unavailable — no failure).
- Card header size remains filesystem scan total (not replaced by runtime stats).

## Side Effects

- None.

## Errors

- SKIP when Docker card not detected.

## Exit Code

- 0 when runtime section present with rows, or gracefully absent.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "SKIP docker-runtime") {
		t.Skip("Docker card not detected on this machine")
	}
	if strings.Contains(resp.Output, "DOCKER_RUNTIME_SECTION: present") {
		for _, elem := range []string{"runtime-row-0", "runtime-label-0", "runtime-count-0", "runtime-size-0"} {
			if strings.Contains(resp.Output, "ELEM "+elem+": MISSING") {
				t.Fatalf("expected %s when runtime section present", elem)
			}
		}
		return
	}
	if !strings.Contains(resp.Output, "DOCKER_RUNTIME_GRACEFUL: true") {
		t.Fatalf("expected runtime section or graceful absence\n%s", resp.Output)
	}
}
```