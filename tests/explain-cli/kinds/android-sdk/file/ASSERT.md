## Expected

- Exit code 0.
- Kind is `android-sdk` (SDK ContextRoot preferred over `generic-file`).
- Output includes the explained file path and/or the SDK context root path.
- Human sections present; no `rm -rf`; trailing blank line.
- HOW TO PURGE still android-sdk-oriented (`sdkmanager` and/or scan); RAW COMMANDS use `$ `
  prefix; no ANSI under default auto.

## Exit Code

- 0

```go
import (
	"path/filepath"
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
	assertKindLine(t, resp.Stdout, "android-sdk")
	assertHumanSectionsPresent(t, resp.Stdout)
	if !strings.Contains(resp.Stdout, req.TargetPath) {
		t.Fatalf("output must include target file path %q:\n%s", req.TargetPath, resp.Stdout)
	}
	// ContextRoot should be the SDK directory (…/Library/Android/sdk), parent of platform-tools.
	// platform-tools/f → Dir → platform-tools → Dir → sdk.
	sdkDir := filepath.Dir(filepath.Dir(req.TargetPath))
	if !strings.Contains(resp.Stdout, sdkDir) {
		t.Fatalf("output must include Android SDK context root %q:\n%s", sdkDir, resp.Stdout)
	}
	// Must not mis-classify as generic-file when inside SDK with signatures.
	for _, line := range strings.Split(resp.Stdout, "\n") {
		if strings.TrimSpace(line) == "KIND: generic-file" {
			t.Fatalf("file under Android SDK must prefer android-sdk, got generic-file:\n%s", resp.Stdout)
		}
	}
	how := sectionBody(resp.Stdout, "HOW TO PURGE")
	lowerHow := strings.ToLower(how)
	if !strings.Contains(lowerHow, "sdkmanager") && !strings.Contains(lowerHow, "disk-usage-analyser scan") {
		t.Fatalf("file-under-SDK explain should still include sdkmanager and/or scan HOW TO PURGE:\n%s", how)
	}
	assertHowToPurgeOfficialDollarPrefix(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, "")
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
