## Expected

- Exit code 0.
- Human sections in order: `PATH:`, `KIND:`, `TOTAL:`, `CONFIDENCE:`, `SUMMARY`, `BREAKDOWN`,
  `SAFE TO RECLAIM`, `HOW TO PURGE`, `RAW COMMANDS`.
- Exact kind line: `KIND: android-sdk`.
- `PATH:` includes the absolute SDK directory.
- BREAKDOWN/summary mentions Android SDK roles or basenames (system-images, emulator,
  sources, build-tools, platform-tools, platforms, tmp / `.temp`).
- BREAKDOWN size DESC roughly: system-images → emulator → sources → build-tools →
  platform-tools → platforms → `.temp`.
- RECLAIMABLE: system-images / sources / tmp (`.temp`) → `☑`; emulator / build-tools /
  platform-tools / platforms → `☐`.
- `SAFE TO RECLAIM`: temp/system-images/sources reclaimable language; platform-tools /
  build-tools / platforms not usually-safe-only purge.
- `HOW TO PURGE`:
  - CLI-first: `$ sdkmanager --list_installed` / `--uninstall` and/or
    `$ disk-usage-analyser scan` on SDK / system-images.
  - Never `rm -rf`.
  - UI (Android Studio SDK settings) only in Notes.
  - Runnable official lines use **`$ `** prefix.
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
		t.Fatalf("PATH/output must include target %q:\n%s", req.TargetPath, resp.Stdout)
	}
	// Must not fall through to generic-dir when signatures are present.
	for _, line := range strings.Split(resp.Stdout, "\n") {
		if strings.TrimSpace(line) == "KIND: generic-dir" {
			t.Fatalf("Android SDK tree must prefer android-sdk, got generic-dir:\n%s", resp.Stdout)
		}
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
