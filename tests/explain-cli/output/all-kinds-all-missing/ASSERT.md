## Expected

- Exit code **0** even when every pack root is missing.
- Human report still prints multi-kind header and INDEX (or equivalent status listing).
- All five cli kinds appear with status **missing** (INDEX and/or body).
- TOTAL (present) is zero / no present packs.
- No `rm -rf`; trailing blank line.

## Exit Code

- 0

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0 when all packs are missing, got %d (err=%v stderr=%q stdout=%q)",
			resp.ExitCode, resp.Err, resp.Stderr, resp.Stdout)
	}
	scope := req.HomeDir
	if scope == "" {
		scope = req.TargetPath
	}
	assertAllKindsHumanHeader(t, resp.Stdout, scope)
	// All v1 packs missing in INDEX (status missing; no present packs).
	assertAllKindsIndex(t, resp.Stdout, nil, allKindsCLIKinds)
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```