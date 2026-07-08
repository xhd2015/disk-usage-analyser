---
label: ui-automation
explanation: Playwright opens xcode cleanup popover and checks expanded tips
---

## Expected

- Popover mentions DerivedData cleanup.
- Popover mentions iOS DeviceSupport cleanup.
- Popover mentions `simctl shutdown all` and `simctl delete all` for simulator devices.
- Popover mentions `simctl runtime delete` with UUID guidance (not SimRuntime bundle IDs).
- At least one suggestion indicates recoverability.
- SKIP when Xcode card not detected.

## Side Effects

- None.

## Errors

- Missing required cleanup topics fails the test.

## Exit Code

- 0 on PASS or SKIP.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil && !strings.Contains(resp.Output, "SKIP xcode-cleanup-popover") {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}
	if strings.Contains(resp.Output, "SKIP xcode-cleanup-popover") {
		t.Skip("Xcode card not detected on this machine")
	}

	checks := map[string]string{
		"CLEANUP_XCODE_DERIVED":         "true",
		"CLEANUP_XCODE_DEVICE_SUPPORT":  "true",
		"CLEANUP_XCODE_SIMULATORS":      "true",
		"CLEANUP_XCODE_RUNTIME_DELETE":  "true",
		"CLEANUP_XCODE_RECOVERABLE":     "true",
	}

	for key, expected := range checks {
		line := findLine(resp.Output, key)
		if line == "" {
			t.Fatalf("missing expected output line: %s", key)
		}
		if !strings.Contains(line, expected) {
			t.Fatalf("line %q: expected to contain %q", line, expected)
		}
	}
}

func findLine(output, prefix string) string {
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			return line
		}
	}
	return ""
}
```