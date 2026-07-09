## Expected

- Exit code 0 (PATH optional when `--kind` is set).
- `KIND: android-sdk`.
- `PATH:` / output reflects the measured `{HomeDir}/Library/Android/sdk` tree (not a missing-path error).
- BREAKDOWN covers Android SDK roles/names + reclaim checkboxes under that home.
- HOW TO PURGE CLI-first (`sdkmanager` / scan); no `rm -rf`; `$` on runnable lines; trailing blank; no ANSI.

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
		t.Fatalf("expected exit 0 with --kind android-sdk and no PATH (HomeDir inject), got %d (err=%v stderr=%q stdout=%q)",
			resp.ExitCode, resp.Err, resp.Stderr, resp.Stdout)
	}
	assertHumanSectionsPresent(t, resp.Stdout)
	assertHumanSectionOrder(t, resp.Stdout)
	assertKindLine(t, resp.Stdout, "android-sdk")
	// Measured context is {HomeDir}/Library/Android/sdk
	sdkDir := req.TargetPath
	if sdkDir == "" || !strings.Contains(resp.Stdout, sdkDir) {
		t.Fatalf("PATH/output must include measured SDK %q:\n%s", sdkDir, resp.Stdout)
	}
	if req.HomeDir != "" && !strings.Contains(sdkDir, req.HomeDir) {
		t.Fatalf("measured SDK %q should be under HomeDir %q", sdkDir, req.HomeDir)
	}
	assertAndroidSDKBreakdownMentions(t, resp.Stdout)
	assertAndroidSDKReclaimCheckboxes(t, resp.Stdout)
	assertAndroidSDKCLIFirstPurge(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, sdkDir)
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
