## Expected

- Non-zero exit code.
- Error indicates an **unknown/unsupported kind value** (registry miss), not merely that
  `--kind` is an unrecognized CLI flag.
- Message should reference the bad kind id (`not-a-kind`) and/or “unknown kind” /
  “invalid kind” / “unsupported kind”.
- Error (or usage) **lists supported kinds**, including **`xcode`**, **`grok`**,
  **`android-sdk`**, **`iterm2`**, and **`codex`**.
- Must not succeed with `KIND: xcode` / `KIND: grok-home` / `KIND: android-sdk` /
  `KIND: iterm2` / `KIND: codex-home` (or any successful pack).

## Exit Code

- Non-zero

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit for unknown --kind, got 0 (err=%v stdout=%q)", resp.Err, resp.Stdout)
	}
	msg := ""
	if resp.Err != nil {
		msg = resp.Err.Error()
	}
	msg += " " + resp.Stderr + " " + resp.Stdout
	lower := strings.ToLower(msg)

	// --kind must be a recognized flag; only the kind *value* is invalid.
	// Fail if the CLI still treats --kind as an unknown/undefined flag (pre-implementation).
	if strings.Contains(lower, "unknown flag") ||
		strings.Contains(lower, "flag provided but not defined") ||
		strings.Contains(lower, "undefined flag") ||
		(strings.Contains(lower, "unknown") && strings.Contains(lower, "flag") && strings.Contains(lower, "kind") &&
			!strings.Contains(lower, "unknown kind") && !strings.Contains(lower, "invalid kind") &&
			!strings.Contains(lower, "unsupported kind")) {
		t.Fatalf("expected --kind accepted as a flag with unknown kind value, got flag-parse style error: %q", msg)
	}

	// Kind-value error: mention the bad id and/or unknown/invalid/unsupported kind.
	hasBadID := strings.Contains(msg, "not-a-kind")
	hasKindValueErr := strings.Contains(lower, "unknown kind") ||
		strings.Contains(lower, "invalid kind") ||
		strings.Contains(lower, "unsupported kind") ||
		strings.Contains(lower, "unrecognized kind") ||
		(strings.Contains(lower, "kind") && (strings.Contains(lower, "unknown") ||
			strings.Contains(lower, "invalid") || strings.Contains(lower, "unsupported") ||
			strings.Contains(lower, "not supported") || strings.Contains(lower, "not found")))
	if !hasBadID && !hasKindValueErr {
		t.Fatalf("error should mention unknown/invalid kind value (e.g. not-a-kind), got exit=%d err=%v stderr=%q stdout=%q",
			resp.ExitCode, resp.Err, resp.Stderr, resp.Stdout)
	}
	// Supported kinds listed for the user (xcode + grok + android-sdk + iterm2 + codex).
	if !strings.Contains(lower, "xcode") {
		t.Fatalf("unknown-kind error should list supported kind xcode: %q", msg)
	}
	if !strings.Contains(lower, "grok") {
		t.Fatalf("unknown-kind error should list supported kind grok: %q", msg)
	}
	if !strings.Contains(lower, "android-sdk") {
		t.Fatalf("unknown-kind error should list supported kind android-sdk: %q", msg)
	}
	if !strings.Contains(lower, "iterm2") {
		t.Fatalf("unknown-kind error should list supported kind iterm2: %q", msg)
	}
	if !strings.Contains(lower, "codex") {
		t.Fatalf("unknown-kind error should list supported kind codex: %q", msg)
	}
	if strings.Contains(resp.Stdout, "KIND: xcode") || strings.Contains(resp.Stdout, "KIND: grok-home") ||
		strings.Contains(resp.Stdout, "KIND: android-sdk") || strings.Contains(resp.Stdout, "KIND: iterm2") ||
		strings.Contains(resp.Stdout, "KIND: codex-home") {
		t.Fatalf("unknown kind must not emit successful KIND line:\n%s", resp.Stdout)
	}
}
```
