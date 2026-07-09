## Expected

- Exit code 0.
- Stdout is one JSON object (plus trailing blank line).
- Required keys: `path`, `kind`, `totalSize`, `breakdown`, `reclaim`, `howToPurge`, `rawCommands`.
- Recommended keys when present: `confidence`, `summary`.
- `kind` is `"iterm2"`; `path` equals the absolute iTerm2 directory.
- `totalSize` is an integer ≥ content payload (`iTerm2ContentBytes` = 674).
- `breakdown` is a non-empty array; entries include roles covering at least:
  `python-env`, `python-env-alias`, `logs`, `meta`, `user-config`.
- `reclaimable` bools: true for python-env, python-env-alias, logs; false for meta, user-config.
- `reclaim` is a non-empty array of advice objects.
- `summary` (when present) and/or reclaim advice mention hardlink/shared blocks / overcount /
  du parent (document-only; no hardlink-aware TOTAL required).
- `howToPurge` is a non-empty array; entries include officialCommand and removes.
- `howToPurge[].officialCommand` is plain (no `$`, no ANSI); includes
  `disk-usage-analyser scan` and/or `du`; never `rm -rf`.
- Entire JSON stdout has **no ANSI** escapes and must not contain `rm -rf`.
- `rawCommands` is non-empty; at least one command contains `disk-usage-analyser scan`;
  command strings have no `$` prefix and no ANSI.

## Exit Code

- 0

```go
import (
	"encoding/json"
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
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)

	payload := parseJSONObject(t, resp.Stdout)
	for _, key := range []string{"path", "kind", "totalSize", "breakdown", "reclaim", "howToPurge", "rawCommands"} {
		if payload[key] == nil {
			t.Fatalf("json missing required key %q: %s", key, firstJSONObjectLine(t, resp.Stdout))
		}
	}

	path := jsonStringField(t, payload, "path")
	if path != req.TargetPath {
		t.Fatalf("json path: expected %q, got %q", req.TargetPath, path)
	}
	kind := jsonStringField(t, payload, "kind")
	if kind != "iterm2" {
		t.Fatalf("json kind: expected iterm2, got %q", kind)
	}
	total := jsonInt64Field(t, payload, "totalSize")
	if total < iTerm2ContentBytes {
		t.Fatalf("json totalSize: expected >= %d (file payload), got %d", iTerm2ContentBytes, total)
	}

	var breakdown []map[string]any
	if err := json.Unmarshal(payload["breakdown"], &breakdown); err != nil {
		t.Fatalf("breakdown: %v", err)
	}
	if len(breakdown) == 0 {
		t.Fatal("breakdown must be non-empty")
	}
	assertJSONBreakdownReclaimableBools(t, breakdown)
	if len(breakdown) >= 2 {
		assertJSONBreakdownSortedDesc(t, breakdown)
	}

	roles := map[string]bool{}
	for _, entry := range breakdown {
		if r, ok := entry["role"].(string); ok && r != "" {
			roles[r] = true
		}
	}
	for _, want := range []string{"python-env", "python-env-alias", "logs", "meta", "user-config"} {
		if !roles[want] {
			t.Fatalf("breakdown roles missing %q (got %v):\n%s", want, roles, firstJSONObjectLine(t, resp.Stdout))
		}
	}
	assertJSONBreakdownRoleReclaimable(t, breakdown, "python-env", true)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "python-env-alias", true)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "logs", true)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "meta", false)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "user-config", false)

	var reclaim []any
	if err := json.Unmarshal(payload["reclaim"], &reclaim); err != nil {
		t.Fatalf("reclaim: %v", err)
	}
	if len(reclaim) == 0 {
		t.Fatal("reclaim must be non-empty")
	}

	// Hardlink / shared-block / overcount / du parent wording in summary and/or reclaim.
	hardlinkBlob := firstJSONObjectLine(t, resp.Stdout)
	if rawSum, ok := payload["summary"]; ok && rawSum != nil {
		hardlinkBlob += " " + string(rawSum)
	}
	reclaimDump, _ := json.Marshal(reclaim)
	hardlinkBlob += " " + string(reclaimDump)
	hl := strings.ToLower(hardlinkBlob)
	hasHardlinkOrShare := strings.Contains(hl, "hardlink") || strings.Contains(hl, "hard link") ||
		strings.Contains(hl, "hard-link") || strings.Contains(hl, "shared inode") ||
		strings.Contains(hl, "shared block") || strings.Contains(hl, "share disk") ||
		strings.Contains(hl, "shared disk") ||
		(strings.Contains(hl, "share") && (strings.Contains(hl, "block") ||
			strings.Contains(hl, "inode") || strings.Contains(hl, "clone") ||
			strings.Contains(hl, "apfs") || strings.Contains(hl, "iterm2env")))
	if !hasHardlinkOrShare {
		t.Fatalf("JSON summary/reclaim must mention hardlink/shared blocks among iterm2env*: %s", hardlinkBlob)
	}
	hasOvercountOrDu := strings.Contains(hl, "overcount") || strings.Contains(hl, "overstate") ||
		strings.Contains(hl, "double-count") || strings.Contains(hl, "double count") ||
		strings.Contains(hl, "logical") || strings.Contains(hl, "du -sh") ||
		strings.Contains(hl, "du -s") || strings.Contains(hl, "freeable") ||
		strings.Contains(hl, "not sum") || strings.Contains(hl, "one env")
	if !hasOvercountOrDu {
		t.Fatalf("JSON summary/reclaim must warn logical overcount and/or du parent: %s", hardlinkBlob)
	}

	var howToPurge []map[string]any
	if err := json.Unmarshal(payload["howToPurge"], &howToPurge); err != nil {
		t.Fatalf("howToPurge: %v", err)
	}
	if len(howToPurge) == 0 {
		t.Fatal("howToPurge must be non-empty")
	}
	for i, step := range howToPurge {
		if step["officialCommand"] == nil || step["officialCommand"] == "" {
			t.Fatalf("howToPurge[%d] missing officialCommand", i)
		}
		if step["removes"] == nil || step["removes"] == "" {
			t.Fatalf("howToPurge[%d] missing removes", i)
		}
	}
	assertJSONHowToPurgePlainCLI(t, howToPurge, false)

	joinedOC := ""
	joinedRemoves := ""
	for _, step := range howToPurge {
		oc, _ := step["officialCommand"].(string)
		rm, _ := step["removes"].(string)
		joinedOC += oc + "\n"
		joinedRemoves += rm + "\n"
	}
	jl := strings.ToLower(joinedOC + "\n" + joinedRemoves)
	if !strings.Contains(jl, "disk-usage-analyser scan") &&
		!strings.Contains(jl, "du -sh") && !strings.Contains(jl, "du -s") {
		t.Fatalf("JSON howToPurge must include disk-usage-analyser scan and/or du: %s", joinedOC)
	}
	if !strings.Contains(jl, "disk-usage-analyser scan") {
		t.Fatalf("JSON howToPurge must include disk-usage-analyser scan: %s", joinedOC)
	}
	if strings.Contains(jl, "rm -rf") || strings.Contains(jl, "rm -fr") {
		t.Fatalf("JSON howToPurge must not contain rm -rf: %s", joinedOC)
	}
	// user-config must not be usually-safe Removes.
	for _, step := range howToPurge {
		rm, _ := step["removes"].(string)
		rl := strings.ToLower(rm)
		if (strings.Contains(rl, "usually safe") || strings.Contains(rl, "usually-safe")) &&
			(strings.Contains(rl, "dynamicprofiles") || strings.Contains(rl, "dynamic profiles") ||
				strings.Contains(rl, "scripts") || strings.Contains(rl, "user-config")) {
			t.Fatalf("JSON howToPurge removes must not mark user-config as usually-safe: %q", rm)
		}
	}

	var rawCommands []any
	if err := json.Unmarshal(payload["rawCommands"], &rawCommands); err != nil {
		t.Fatalf("rawCommands: %v", err)
	}
	if len(rawCommands) == 0 {
		t.Fatal("rawCommands must be non-empty")
	}
	rawDump, _ := json.Marshal(rawCommands)
	if !strings.Contains(string(rawDump), "disk-usage-analyser scan") {
		t.Fatalf("rawCommands must include disk-usage-analyser scan: %s", rawDump)
	}
	for i, rc := range rawCommands {
		switch v := rc.(type) {
		case map[string]any:
			if cmd, ok := v["command"].(string); ok {
				if containsANSI(cmd) {
					t.Fatalf("rawCommands[%d].command must not contain ANSI: %q", i, cmd)
				}
				if strings.HasPrefix(strings.TrimSpace(cmd), "$") {
					t.Fatalf("rawCommands[%d].command must not include $ prefix: %q", i, cmd)
				}
			}
		case string:
			if containsANSI(v) {
				t.Fatalf("rawCommands[%d] must not contain ANSI: %q", i, v)
			}
			if strings.HasPrefix(strings.TrimSpace(v), "$") {
				t.Fatalf("rawCommands[%d] must not include $ prefix: %q", i, v)
			}
		}
	}
}
```