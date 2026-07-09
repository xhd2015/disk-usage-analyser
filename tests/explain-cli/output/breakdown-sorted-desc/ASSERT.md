## Expected

- Exit code 0.
- `KIND: android-avd`.
- Human BREAKDOWN lists known AVD basenames in **size DESC** order:
  `userdata-qemu.img.qcow2` → `sdcard.img` → `snapshots` → `config.ini`.
- Table header present; trailing blank line; no ANSI; no `rm -rf`.

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
		t.Fatalf("expected exit 0, got %d (err=%v stderr=%q)", resp.ExitCode, resp.Err, resp.Stderr)
	}
	assertKindLine(t, resp.Stdout, "android-avd")
	assertBreakdownTableHeader(t, resp.Stdout)
	// Exact payload sizes from writeAVDFixture: 400, 200, 100, 32.
	assertBreakdownNamesInOrder(t, resp.Stdout, []string{
		"userdata-qemu.img.qcow2",
		"sdcard.img",
		"snapshots",
		"config.ini",
	})
	assertNoANSI(t, resp.Stdout)
	assertNoRmRf(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
