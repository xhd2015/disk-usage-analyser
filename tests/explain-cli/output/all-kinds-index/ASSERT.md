## Expected

- Exit code 0 (PATH optional when `--all-kinds` is set; missing packs do not fail the run).
- Human multi-kind header: `SCOPE:`, `MODE: all-kinds`, `TOTAL` (present).
- `SCOPE` / output reflects the injected `HomeDir`.
- `INDEX` table lists all five v1 **cli** kinds: `xcode`, `grok`, `android-sdk`, `iterm2`, `codex`.
- INDEX columns include SIZE, KIND, STATUS, PATH (optional NOTE).
- Present packs (codex, android-sdk, grok, iterm2) marked **present**; **xcode** marked **missing**.
- Present rows ordered **size DESC** (fixture: codex sqlite multi-KB > android-sdk > grok > iterm2).
- No `rm -rf`; no ANSI (non-TTY auto); trailing blank line after last content.

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
		t.Fatalf("expected exit 0 for --all-kinds, got %d (err=%v stderr=%q stdout=%q)",
			resp.ExitCode, resp.Err, resp.Stderr, resp.Stdout)
	}
	scope := req.HomeDir
	if scope == "" {
		scope = req.TargetPath
	}
	assertAllKindsHumanHeader(t, resp.Stdout, scope)
	// Present packs in writeAllKindsHomeFixture ordered by fixture payload size DESC:
	// codex (≥1198 non-DB) > android-sdk (890) > grok (798) > iterm2 (674); xcode missing.
	assertAllKindsIndex(t, resp.Stdout,
		[]string{"codex", "android-sdk", "grok", "iterm2"},
		[]string{"xcode"},
	)
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
