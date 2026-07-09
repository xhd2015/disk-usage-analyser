## Expected

- Exit code 0.
- Help text documents `explain`, PATH, **`--kind`**, **`--all-kinds`**, `--json`, `--color`,
  and `-h`/`--help`.
- PATH is described as required **unless** `--kind` and/or **`--all-kinds`** is set
  (optional scope / default home when either flag is set).
- **`--kind`** documents supported pack/alias ids including **`xcode`**, **`grok`**,
  **`android-sdk`**, **`iterm2`**, and **`codex`** (e.g. wording like
  `xcode, grok, android-sdk, iterm2, codex` or separate mentions).
- **`--all-kinds`** is documented (analyse all registered packs under optional scope).

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
		t.Fatalf("expected exit code 0, got %d (err=%v)", resp.ExitCode, resp.Err)
	}
	out := resp.Stdout + resp.Stderr
	for _, want := range []string{
		"explain",
		"PATH",
		"--kind",
		"--all-kinds",
		"--json",
		"--color",
		"-h",
		"--help",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("help output missing %q:\n%s", want, out)
		}
	}
	// PATH optional when --kind or --all-kinds is set (wording may vary).
	lower := strings.ToLower(out)
	hasOptionalScope := strings.Contains(lower, "unless") ||
		strings.Contains(lower, "optional") ||
		((strings.Contains(lower, "kind") || strings.Contains(lower, "all-kinds")) &&
			(strings.Contains(lower, "not required") ||
				strings.Contains(lower, "without path") || strings.Contains(lower, "default") ||
				strings.Contains(lower, "home")))
	if !hasOptionalScope {
		t.Fatalf("help should document that PATH is optional when --kind or --all-kinds is set (optional scope / default home):\n%s", out)
	}
	// --all-kinds should be described as multi-pack / all registered kinds (soft cues).
	if !strings.Contains(lower, "all-kinds") && !strings.Contains(lower, "all kinds") {
		t.Fatalf("help must document --all-kinds:\n%s", out)
	}
	// --kind supported values must mention xcode, grok, android-sdk, iterm2, and codex.
	if !strings.Contains(lower, "grok") {
		t.Fatalf("help --kind docs must mention supported kind grok:\n%s", out)
	}
	if !strings.Contains(lower, "xcode") {
		t.Fatalf("help --kind docs must mention supported kind xcode:\n%s", out)
	}
	if !strings.Contains(lower, "android-sdk") {
		t.Fatalf("help --kind docs must mention supported kind android-sdk:\n%s", out)
	}
	if !strings.Contains(lower, "iterm2") {
		t.Fatalf("help --kind docs must mention supported kind iterm2:\n%s", out)
	}
	if !strings.Contains(lower, "codex") {
		t.Fatalf("help --kind docs must mention supported kind codex:\n%s", out)
	}
}
```