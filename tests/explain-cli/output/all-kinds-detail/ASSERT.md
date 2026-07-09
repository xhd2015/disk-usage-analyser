## Expected

- Exit code 0.
- Human multi-kind header (`SCOPE:`, `MODE: all-kinds`) with scope = PATH argument.
- For each **present** pack, a mini-explain detail includes:
  - `KIND: android-sdk`, `KIND: grok-home`, `KIND: iterm2`, `KIND: codex-home` (output kind ids)
  - `BREAKDOWN` section(s) covering those present kinds
- Missing **xcode** must not appear as a successful full present detail (`KIND: xcode` with
  non-empty BREAKDOWN success path is not required; soft: no present status forced for xcode).
- Full stdout must not contain `rm -rf`; no ANSI; trailing blank line.

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
		t.Fatalf("expected exit 0 for --all-kinds <SCOPE>, got %d (err=%v stderr=%q stdout=%q)",
			resp.ExitCode, resp.Err, resp.Stderr, resp.Stdout)
	}
	assertAllKindsHumanHeader(t, resp.Stdout, req.TargetPath)
	// Present output kind ids (cli grok → kind grok-home; cli codex → kind codex-home).
	assertAllKindsDetailPresent(t, resp.Stdout, []string{
		"android-sdk",
		"grok-home",
		"iterm2",
		"codex-home",
	})
	// Soft role/name cues from fixture packs inside the combined report.
	lower := strings.ToLower(resp.Stdout)
	if !strings.Contains(lower, "system-images") && !strings.Contains(lower, "emulator") {
		t.Fatalf("detail should mention android-sdk breakdown cues (system-images/emulator):\n%s", resp.Stdout)
	}
	if !strings.Contains(lower, "downloads") && !strings.Contains(lower, "installer-cache") {
		t.Fatalf("detail should mention grok-home breakdown cues (downloads/installer-cache):\n%s", resp.Stdout)
	}
	if !strings.Contains(lower, "iterm2env") && !strings.Contains(lower, "python-env") {
		t.Fatalf("detail should mention iterm2 breakdown cues (iterm2env/python-env):\n%s", resp.Stdout)
	}
	if !strings.Contains(lower, "logs_") && !strings.Contains(lower, "app-logs-db") &&
		!strings.Contains(lower, "codex") {
		t.Fatalf("detail should mention codex-home breakdown cues (logs_*.sqlite / app-logs-db):\n%s", resp.Stdout)
	}
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
