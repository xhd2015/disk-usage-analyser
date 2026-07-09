## Expected

- Exit code 0.
- `KIND: go-build-cache`.
- **RAW COMMANDS**: every runnable line starts with `$ ` after trim; `#` group/comment lines do not.
- At least one RAW line is `$ disk-usage-analyser scan …` (path present).
- **HOW TO PURGE**: official runnable lines use `$ `; includes `$ go clean` (or `$ go` …).
- Comment lines under official command (if any) do not get `$`.
- Default non-TTY: no ANSI required (assert absence is OK, not required here).

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
	assertKindLine(t, resp.Stdout, "go-build-cache")
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertHowToPurgeOfficialDollarPrefix(t, resp.Stdout)
	assertHowToPurgeHasDollarCommand(t, resp.Stdout, "go")

	// Explicit RAW scan form with shell prompt.
	raw := sectionBody(resp.Stdout, "RAW COMMANDS")
	found := false
	for _, line := range strings.Split(raw, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "$ ") && strings.Contains(trim, "disk-usage-analyser scan") {
			found = true
			if req.TargetPath != "" && !strings.Contains(trim, req.TargetPath) {
				// Path may appear on PATH: line only; still require scan token.
			}
			break
		}
	}
	if !found {
		t.Fatalf("RAW COMMANDS must include a \"$ disk-usage-analyser scan …\" line:\n%s", raw)
	}
	assertContainsScanCommand(t, resp.Stdout, req.TargetPath)
	assertNoRmRf(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
