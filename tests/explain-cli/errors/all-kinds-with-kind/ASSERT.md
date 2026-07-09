## Expected

- Non-zero exit code when **`--all-kinds`** and **`--kind`** are both set.
- Error (or stderr/stdout) indicates the flags are **mutually exclusive** / cannot be combined
  (wording may vary: "mutually exclusive", "cannot be used with", "not both", "conflict", …).
- Message should reference **`--all-kinds`** and/or **`--kind`** (or the kind value).
- Must not emit a successful multi-kind report (`MODE: all-kinds` with INDEX + present details)
  nor a successful single-kind `KIND: xcode` explain.

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
		t.Fatalf("expected non-zero exit when --all-kinds and --kind are combined, got 0 (err=%v stdout=%q)",
			resp.Err, resp.Stdout)
	}
	msg := ""
	if resp.Err != nil {
		msg = resp.Err.Error()
	}
	msg += " " + resp.Stderr + " " + resp.Stdout
	lower := strings.ToLower(msg)

	// Prefer mutual-exclusion language; also accept "cannot" + both flag names.
	hasExclusive := strings.Contains(lower, "mutually exclusive") ||
		strings.Contains(lower, "exclusive") ||
		strings.Contains(lower, "conflict") ||
		strings.Contains(lower, "incompatible") ||
		strings.Contains(lower, "cannot be used together") ||
		strings.Contains(lower, "cannot combine") ||
		strings.Contains(lower, "not both") ||
		(strings.Contains(lower, "cannot") && (strings.Contains(lower, "all-kinds") || strings.Contains(lower, "kind"))) ||
		(strings.Contains(lower, "not allowed") && (strings.Contains(lower, "all-kinds") || strings.Contains(lower, "kind")))
	hasFlagNames := strings.Contains(lower, "all-kinds") || strings.Contains(lower, "all kinds") ||
		strings.Contains(lower, "--kind") || strings.Contains(lower, "kind")
	if !hasExclusive && !hasFlagNames {
		t.Fatalf("error should indicate --all-kinds/--kind mutual exclusion, got exit=%d err=%v stderr=%q stdout=%q",
			resp.ExitCode, resp.Err, resp.Stderr, resp.Stdout)
	}
	if !hasExclusive {
		// Soft fail only if neither exclusive cue nor both flags mentioned.
		if !(strings.Contains(lower, "all-kinds") && strings.Contains(lower, "kind")) {
			t.Fatalf("error should mention mutual exclusion or both --all-kinds and --kind, got: %q", msg)
		}
	}

	// Must not succeed as all-kinds report or single xcode explain.
	outLower := strings.ToLower(resp.Stdout)
	if strings.Contains(outLower, "mode:") && strings.Contains(outLower, "all-kinds") &&
		strings.Contains(resp.Stdout, "INDEX") {
		t.Fatalf("mutual exclusion must not emit successful all-kinds INDEX report:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "KIND: xcode") {
		t.Fatalf("mutual exclusion must not emit successful KIND: xcode:\n%s", resp.Stdout)
	}
}
```
