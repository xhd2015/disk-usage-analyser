## Expected

- Exit code 0.
- `HOW TO PURGE` section present with Official command lines.
- Official commands are **CLI-first**: include `emulator` and/or `avdmanager`.
- Runnable official lines are not UI navigation (`Android Studio` / `Device Manager`).
- If Device Manager / Android Studio appears, it is only under **Notes** (e.g. `UI (optional): …`).
- No `rm -rf`. Default non-TTY output does not require ANSI.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d (err=%v stderr=%q)", resp.ExitCode, resp.Err, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "HOW TO PURGE") {
		t.Fatalf("missing HOW TO PURGE:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Official command:") {
		t.Fatalf("HOW TO PURGE should include Official command:\n%s", resp.Stdout)
	}
	assertAndroidAVDCLIFirstPurge(t, resp.Stdout)
	assertNoRmRf(t, resp.Stdout)

	// Any Device Manager / Android Studio mention must sit on a Notes line.
	how := sectionBody(resp.Stdout, "HOW TO PURGE")
	for _, line := range strings.Split(how, "\n") {
		trim := strings.TrimSpace(line)
		if trim == "" {
			continue
		}
		if strings.Contains(trim, "Device Manager") || strings.Contains(trim, "Android Studio") {
			if !strings.HasPrefix(trim, "Notes:") {
				t.Fatalf("UI navigation must only appear on Notes lines, got: %q", line)
			}
		}
	}
}
```
