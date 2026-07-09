## Expected

- Exit code 0.
- Human sections in order with exact kind line: `KIND: android-sdk`.
- Output includes the measured absolute SDK path under the scope (and typically the scope).
- BREAKDOWN mentions Android SDK roles/names; RECLAIMABLE ☑ for system-images/sources/tmp,
  ☐ for emulator/build-tools/platform-tools/platforms.
- BREAKDOWN size DESC: system-images → emulator → sources → build-tools → platform-tools → platforms.
- `SAFE TO RECLAIM` present with reclaim language for temp/system-images; not platform-tools bulk alone.
- `HOW TO PURGE` is CLI-first (`sdkmanager` and/or `disk-usage-analyser scan`); runnable lines use
  **`$ `**; never `rm -rf`.
- `RAW COMMANDS` includes `$ disk-usage-analyser scan` and the SDK path.
- Full stdout: no `rm -rf`; trailing blank line; no ANSI under default non-TTY auto color.

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
	assertHumanSectionsPresent(t, resp.Stdout)
	assertHumanSectionOrder(t, resp.Stdout)
	assertKindLine(t, resp.Stdout, "android-sdk")
	if !strings.Contains(resp.Stdout, req.TargetPath) {
		t.Fatalf("PATH/output must include measured SDK %q:\n%s", req.TargetPath, resp.Stdout)
	}
	assertAndroidSDKBreakdownMentions(t, resp.Stdout)
	assertBreakdownTableHeader(t, resp.Stdout)
	assertBreakdownNoRoleEqualsPrefix(t, resp.Stdout)
	assertBreakdownNamesInOrder(t, resp.Stdout, []string{
		"system-images",
		"emulator",
		"sources",
		"build-tools",
		"platform-tools",
		"platforms",
	})
	assertAndroidSDKReclaimCheckboxes(t, resp.Stdout)
	assertAndroidSDKSafeToReclaim(t, resp.Stdout)
	assertAndroidSDKCLIFirstPurge(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, req.TargetPath)
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
