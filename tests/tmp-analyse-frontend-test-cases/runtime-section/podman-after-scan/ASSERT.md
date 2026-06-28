## Expected

- After scan completes, if `PODMAN_RUNTIME_SECTION: present` then runtime-row-0, runtime-label-0, runtime-count-0, and runtime-size-0 must be present.
- If `PODMAN_RUNTIME_SECTION: absent`, `PODMAN_RUNTIME_GRACEFUL: true` is logged.
- Card header size remains filesystem scan total.

## Side Effects

- None.

## Errors

- SKIP when Podman card not detected.

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
	if strings.Contains(resp.Output, "SKIP podman-runtime") {
		t.Skip("Podman card not detected on this machine")
	}
	if strings.Contains(resp.Output, "PODMAN_RUNTIME_SECTION: present") {
		for _, elem := range []string{"runtime-row-0", "runtime-label-0", "runtime-count-0", "runtime-size-0"} {
			if strings.Contains(resp.Output, "ELEM "+elem+": MISSING") {
				t.Fatalf("expected %s when runtime section present", elem)
			}
		}
		return
	}
	if !strings.Contains(resp.Output, "PODMAN_RUNTIME_GRACEFUL: true") {
		t.Fatalf("expected runtime section or graceful absence\n%s", resp.Output)
	}
}
```