## Expected

- Exit code 0.
- Kind is `android-avd` (parent AVD context preferred over bare `generic-qcow2`).
- Output includes the explained file path.
- Human sections present; no `rm -rf`; trailing blank line.
- HOW TO PURGE CLI-first; RAW COMMANDS use `$ ` prefix; no ANSI under default auto.

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
	assertKindLine(t, resp.Stdout, "android-avd")
	assertHumanSectionsPresent(t, resp.Stdout)
	if !strings.Contains(resp.Stdout, req.TargetPath) {
		t.Fatalf("output must include target file path %q:\n%s", req.TargetPath, resp.Stdout)
	}
	// Must not mis-classify as only generic-qcow2 when inside *.avd.
	for _, line := range strings.Split(resp.Stdout, "\n") {
		if strings.TrimSpace(line) == "KIND: generic-qcow2" {
			t.Fatalf("file inside AVD must prefer android-avd, got generic-qcow2:\n%s", resp.Stdout)
		}
	}
	assertAndroidAVDCLIFirstPurge(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertNoRmRf(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, "")
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
