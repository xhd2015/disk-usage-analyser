## Expected

- Exit code 0.
- Human sections in order: `PATH:`, `KIND:`, `TOTAL:`, `CONFIDENCE:`, `SUMMARY`, `BREAKDOWN`,
  `SAFE TO RECLAIM`, `HOW TO PURGE`, `RAW COMMANDS`.
- Exact kind line: `KIND: android-avd`.
- `PATH:` includes the absolute AVD directory.
- Breakdown mentions AVD artifacts (e.g. userdata, sdcard, snapshots) and/or roles.
- `SAFE TO RECLAIM` present; no `rm -rf` anywhere in stdout.
- `HOW TO PURGE` is **CLI-first** (`emulator` / `avdmanager`); Official command + Removes present.
- Runnable HOW TO PURGE / RAW COMMANDS lines use **`$ `** prefix.
- `RAW COMMANDS` includes `$ disk-usage-analyser scan` and the AVD path.
- Stdout ends with a trailing blank line; no ANSI under default non-TTY auto color.

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
	assertKindLine(t, resp.Stdout, "android-avd")
	if !strings.Contains(resp.Stdout, req.TargetPath) {
		t.Fatalf("PATH/output must include target %q:\n%s", req.TargetPath, resp.Stdout)
	}
	// Breakdown should reference known AVD basenames or roles.
	lower := strings.ToLower(resp.Stdout)
	hasArtifact := strings.Contains(lower, "userdata") ||
		strings.Contains(lower, "sdcard") ||
		strings.Contains(lower, "snapshot") ||
		strings.Contains(lower, "qcow2") ||
		strings.Contains(lower, "user-data")
	if !hasArtifact {
		t.Fatalf("BREAKDOWN/summary should mention AVD artifacts (userdata/sdcard/snapshots):\n%s", resp.Stdout)
	}
	assertNoRmRf(t, resp.Stdout)
	if !strings.Contains(resp.Stdout, "HOW TO PURGE") {
		t.Fatalf("missing HOW TO PURGE:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Official command:") {
		t.Fatalf("HOW TO PURGE should include Official command:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Removes:") {
		t.Fatalf("HOW TO PURGE should include Removes:\n%s", resp.Stdout)
	}
	assertAndroidAVDCLIFirstPurge(t, resp.Stdout)
	assertHowToPurgeOfficialDollarPrefix(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, req.TargetPath)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
