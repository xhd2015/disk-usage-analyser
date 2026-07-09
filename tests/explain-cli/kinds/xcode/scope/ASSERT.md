## Expected

- Exit code 0.
- Human sections in order: `PATH:`, `KIND:`, `TOTAL:`, `CONFIDENCE:`, `SUMMARY`, `BREAKDOWN`,
  `SAFE TO RECLAIM`, `HOW TO PURGE`, `RAW COMMANDS`.
- Exact kind line: `KIND: xcode`.
- `PATH:` includes the absolute scope (fixture home).
- BREAKDOWN mentions all five pack roles/names: derived-data/DerivedData, simulator/Devices,
  device-support/DeviceSupport, archives/Archives, docs-cache/DocumentationCache.
- BREAKDOWN size DESC roughly: DerivedData → Devices → Archives → DeviceSupport → DocumentationCache.
- RECLAIMABLE: derived-data / docs-cache / device-support / simulator → `☑`; archives → `☐`.
- `SAFE TO RECLAIM` present with reclaim/safe language for rebuildable roots.
- `HOW TO PURGE` is CLI-first (`xcrun` / `simctl` and/or `disk-usage-analyser scan`); runnable
  lines use **`$ `**; never `rm -rf`; UI only in Notes.
- `RAW COMMANDS` includes `$ disk-usage-analyser scan` and the scope path.
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
	assertKindLine(t, resp.Stdout, "xcode")
	if !strings.Contains(resp.Stdout, req.TargetPath) {
		t.Fatalf("PATH/output must include scope %q:\n%s", req.TargetPath, resp.Stdout)
	}
	assertXcodeBreakdownMentions(t, resp.Stdout)
	assertBreakdownTableHeader(t, resp.Stdout)
	assertBreakdownNoRoleEqualsPrefix(t, resp.Stdout)
	// Size DESC name order for full fixture (basename tokens; DeviceSupport after Devices).
	assertBreakdownNamesInOrder(t, resp.Stdout, []string{
		"DerivedData",
		"Devices",
		"Archives",
		"DeviceSupport",
		"DocumentationCache",
	})
	assertXcodeReclaimCheckboxes(t, resp.Stdout)
	assertXcodeSafeToReclaim(t, resp.Stdout)
	assertXcodeCLIFirstPurge(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, req.TargetPath)
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
