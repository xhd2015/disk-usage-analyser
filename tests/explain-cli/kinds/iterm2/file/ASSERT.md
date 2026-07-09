## Expected

- Exit code 0.
- Kind is `iterm2` (iTerm2 ContextRoot preferred over `generic-file`).
- Output includes the explained file path and/or the iTerm2 context root path.
- Human sections present; no `rm -rf`; trailing blank line.
- HOW TO PURGE still iterm2-oriented (`scan` and/or `du`); RAW COMMANDS use `$ `
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
	assertKindLine(t, resp.Stdout, "iterm2")
	assertHumanSectionsPresent(t, resp.Stdout)
	if !strings.Contains(resp.Stdout, req.TargetPath) {
		t.Fatalf("output must include target file path %q:\n%s", req.TargetPath, resp.Stdout)
	}
	// ContextRoot should be the iTerm2 directory (…/Application Support/iTerm2),
	// parent of iterm2env. iterm2env/f → Dir → iterm2env → Dir → iTerm2.
	iterm2Dir := filepath.Dir(filepath.Dir(req.TargetPath))
	if !strings.Contains(resp.Stdout, iterm2Dir) {
		t.Fatalf("output must include iTerm2 context root %q:\n%s", iterm2Dir, resp.Stdout)
	}
	// Must not mis-classify as generic-file when inside iTerm2 with signatures.
	for _, line := range strings.Split(resp.Stdout, "\n") {
		if strings.TrimSpace(line) == "KIND: generic-file" {
			t.Fatalf("file under iTerm2 must prefer iterm2, got generic-file:\n%s", resp.Stdout)
		}
	}
	how := sectionBody(resp.Stdout, "HOW TO PURGE")
	lowerHow := strings.ToLower(how)
	if !strings.Contains(lowerHow, "disk-usage-analyser scan") &&
		!strings.Contains(lowerHow, "du -sh") && !strings.Contains(lowerHow, "du -s") {
		t.Fatalf("file-under-iTerm2 explain should still include scan and/or du HOW TO PURGE:\n%s", how)
	}
	assertHowToPurgeOfficialDollarPrefix(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, "")
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
