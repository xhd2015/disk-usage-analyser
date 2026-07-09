## Expected

- Exit code 0.
- Stdout is one JSON object (optional pretty-print; trailing blank line).
- Required keys: `scope`, `totalSize`, `kinds`.
- `scope` equals injected HomeDir (absolute).
- `kinds` is an array of length 5 (v1 packs).
- Each entry has: `kind`, `cliKind`, `path`, `status`, `totalSize`.
- Present entries: `cliKind` android-sdk / grok / iterm2 / codex with `status=present` and
  `totalSize` ≥ fixture payloads (890 / 798 / 674 / non-DB 1198); output `kind` ids
  android-sdk / grok-home / iterm2 / codex-home.
- Missing entry: `cliKind=xcode`, `status=missing`, `totalSize=0`.
- Envelope `totalSize` equals sum of present entry totalSizes (or ≥ sum of present payloads
  and excludes missing).
- Entire JSON has **no ANSI** and must not contain `rm -rf`.

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
		t.Fatalf("expected exit 0 for --json --all-kinds, got %d (err=%v stderr=%q)",
			resp.ExitCode, resp.Err, resp.Stderr)
	}
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)

	scope := req.HomeDir
	if scope == "" {
		scope = req.TargetPath
	}
	payload := parseJSONStdoutFlexible(t, resp.Stdout)
	kinds := assertAllKindsJSONShape(t, payload, scope)
	byCLI := allKindsByCLIKind(t, kinds)

	// Required cli kinds present as keys.
	for _, cli := range allKindsCLIKinds {
		if byCLI[cli] == nil {
			t.Fatalf("kinds missing cliKind %q: %v", cli, kinds)
		}
	}

	// xcode missing, size 0
	xc := byCLI["xcode"]
	if st, _ := xc["status"].(string); st != "missing" {
		t.Fatalf("xcode status: want missing, got %q (entry=%v)", st, xc)
	}
	if sz := jsonEntryInt64(t, xc, "totalSize"); sz != 0 {
		t.Fatalf("xcode totalSize: want 0, got %d", sz)
	}

	// Present packs: status + minimum payload sizes + output kind ids.
	type presentExpect struct {
		cliKind    string
		outputKind string
		minSize    int64
	}
	for _, pe := range []presentExpect{
		{"android-sdk", "android-sdk", androidSDKContentBytes},
		{"grok", "grok-home", grokHomeContentBytes},
		{"iterm2", "iterm2", iTerm2ContentBytes},
		{"codex", "codex-home", codexHomeNonDBContentBytes},
	} {
		e := byCLI[pe.cliKind]
		if st, _ := e["status"].(string); st != "present" {
			t.Fatalf("%s status: want present, got %q (entry=%v)", pe.cliKind, st, e)
		}
		kind, _ := e["kind"].(string)
		if kind != pe.outputKind {
			t.Fatalf("%s kind: want %q, got %q", pe.cliKind, pe.outputKind, kind)
		}
		sz := jsonEntryInt64(t, e, "totalSize")
		if sz < pe.minSize {
			t.Fatalf("%s totalSize: expected >= %d (payload), got %d", pe.cliKind, pe.minSize, sz)
		}
		path, _ := e["path"].(string)
		if strings.TrimSpace(path) == "" {
			t.Fatalf("%s path must be non-empty", pe.cliKind)
		}
	}

	// Envelope totalSize is sum of present only (must match sum of present entry sizes).
	var sumPresent int64
	for _, cli := range []string{"android-sdk", "grok", "iterm2", "codex"} {
		sumPresent += jsonEntryInt64(t, byCLI[cli], "totalSize")
	}
	total := jsonInt64Field(t, payload, "totalSize")
	if total != sumPresent {
		t.Fatalf("envelope totalSize=%d, want sum of present kinds %d", total, sumPresent)
	}
	// Lower bound: at least sum of fixture payloads (codex uses non-DB floor).
	minPayload := androidSDKContentBytes + grokHomeContentBytes + iTerm2ContentBytes + codexHomeNonDBContentBytes
	if total < minPayload {
		t.Fatalf("envelope totalSize=%d < sum of present payloads %d", total, minPayload)
	}
}
```
