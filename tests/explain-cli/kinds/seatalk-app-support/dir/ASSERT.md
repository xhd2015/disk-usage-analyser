## Expected

- Exit code 0.
- Human sections in order: `PATH:`, `KIND:`, `TOTAL:`, `CONFIDENCE:`, `SUMMARY`, `BREAKDOWN`,
  `SAFE TO RECLAIM`, `HOW TO PURGE`, `RAW COMMANDS`.
- Exact kind line: `KIND: seatalk-app-support`.
- `PATH:` includes the absolute SeaTalk Application Support directory.
- BREAKDOWN/summary mentions SeaTalk roles or basenames (web-cache / Service Worker / Cache,
  chat-db / main_, search-index / search_, backup / sqlite-backup, config).
- `SAFE TO RECLAIM`: web caches usually safe; chat/search use caution (not usually-safe purge).
- `HOW TO PURGE`:
  - Official prep: `$ osascript -e 'quit app "SeaTalk"'` (or equivalent osascript quit).
  - Cache reclaim Removes list includes `Service Worker` and `Cache` (and related Chromium caches when listed).
  - Separate backup step mentions `sqlite-backup`.
  - Does not mark primary `main_` / `search_` sqlite as usually-safe Removes.
  - Runnable official lines use **`$ `** prefix.
- `RAW COMMANDS` includes `$ disk-usage-analyser scan` and the SeaTalk path.
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
	assertKindLine(t, resp.Stdout, "seatalk-app-support")
	if !strings.Contains(resp.Stdout, req.TargetPath) {
		t.Fatalf("PATH/output must include target %q:\n%s", req.TargetPath, resp.Stdout)
	}
	assertSeaTalkBreakdownMentions(t, resp.Stdout)
	assertSeaTalkReclaimTiers(t, resp.Stdout)
	assertSeaTalkHumanHowToPurge(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, req.TargetPath)
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
