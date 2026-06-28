## Expected

- After scan completes, if `PODMAN_VM_INTERNAL_SECTION: present` then `vm-internal-row-0`, `vm-internal-label-0`, and `vm-internal-size-0` must be present.
- If `PODMAN_VM_INTERNAL_SECTION: absent`, `PODMAN_VM_INTERNAL_GRACEFUL: true` is logged.
- Card header size remains host filesystem scan total.

## Side Effects

- None.

## Errors

- SKIP when Podman card not detected.

## Exit Code

- 0 when vm-internal section present with rows, or gracefully absent.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "SKIP podman-vm-internal") {
		t.Skip("Podman card not detected on this machine")
	}
	if strings.Contains(resp.Output, "PODMAN_VM_INTERNAL_SECTION: present") {
		for _, elem := range []string{"vm-internal-row-0", "vm-internal-label-0", "vm-internal-size-0"} {
			if strings.Contains(resp.Output, "ELEM "+elem+": MISSING") {
				t.Fatalf("expected %s when vm-internal section present", elem)
			}
		}
		return
	}
	if !strings.Contains(resp.Output, "PODMAN_VM_INTERNAL_GRACEFUL: true") {
		t.Fatalf("expected vm-internal section or graceful absence\n%s", resp.Output)
	}
}
```