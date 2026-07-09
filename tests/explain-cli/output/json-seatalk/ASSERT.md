## Expected

- Exit code 0.
- Stdout is one JSON object (plus trailing blank line).
- Required keys: `path`, `kind`, `totalSize`, `breakdown`, `reclaim`, `howToPurge`, `rawCommands`.
- Recommended keys when present: `confidence`, `summary`.
- `kind` is `"seatalk-app-support"`; `path` equals the absolute SeaTalk directory.
- `totalSize` is an integer ≥ content payload (`seatalkContentBytes` = 510).
- `breakdown` is a non-empty array; entries include roles covering at least:
  `web-cache`, `chat-db`, `search-index`, `backup`, `config` (from the fixture).
- `reclaim` is a non-empty array of advice objects (title / safeToReclaim / detail).
- `howToPurge` is a non-empty array; entries include officialCommand and removes.
- `howToPurge[].officialCommand` is plain (no `$`, no ANSI); includes `osascript` quit
  for SeaTalk; overall steps cover web-cache basenames and `sqlite-backup`.
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
	if kind != "seatalk-app-support" {
		t.Fatalf("json kind: expected seatalk-app-support, got %q", kind)
	}
	total := jsonInt64Field(t, payload, "totalSize")
	if total < seatalkContentBytes {
		t.Fatalf("json totalSize: expected >= %d (file payload), got %d", seatalkContentBytes, total)
	}

	var breakdown []map[string]any
	if err := json.Unmarshal(payload["breakdown"], &breakdown); err != nil {
		t.Fatalf("breakdown: %v", err)
	}
	if len(breakdown) == 0 {
		t.Fatal("breakdown must be non-empty")
	}
	roles := map[string]bool{}
	for _, entry := range breakdown {
		if r, ok := entry["role"].(string); ok && r != "" {
			roles[r] = true
		}
	}
	for _, want := range []string{"web-cache", "chat-db", "search-index", "backup", "config"} {
		if !roles[want] {
			t.Fatalf("breakdown roles missing %q (got %v):\n%s", want, roles, firstJSONObjectLine(t, resp.Stdout))
		}
	}

	var reclaim []any
	if err := json.Unmarshal(payload["reclaim"], &reclaim); err != nil {
		t.Fatalf("reclaim: %v", err)
	}
	if len(reclaim) == 0 {
		t.Fatal("reclaim must be non-empty")
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

	// SeaTalk-specific howToPurge content (plain strings, no $).
	joinedOC := ""
	joinedRemoves := ""
	for _, step := range howToPurge {
		oc, _ := step["officialCommand"].(string)
		rm, _ := step["removes"].(string)
		joinedOC += oc + "\n"
		joinedRemoves += rm + "\n"
	}
	if !strings.Contains(joinedOC, "osascript") {
		t.Fatalf("JSON howToPurge officialCommand must include osascript quit: %s", joinedOC)
	}
	if !strings.Contains(joinedOC, "SeaTalk") {
		t.Fatalf("JSON howToPurge officialCommand must mention SeaTalk: %s", joinedOC)
	}
	if !strings.Contains(joinedRemoves, "Service Worker") && !strings.Contains(joinedRemoves, "Cache") {
		t.Fatalf("JSON howToPurge removes must list web-cache basenames (Service Worker / Cache): %s", joinedRemoves)
	}
	if !strings.Contains(joinedRemoves, "sqlite-backup") && !strings.Contains(joinedOC, "sqlite-backup") {
		t.Fatalf("JSON howToPurge must cover sqlite-backup: oc=%s removes=%s", joinedOC, joinedRemoves)
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
