## Expected

- Exit code 0.
- Stdout is one JSON object (plus trailing blank line).
- Required keys: `path`, `kind`, `totalSize`, `breakdown`, `reclaim`, `howToPurge`, `rawCommands`,
  **`logsDb`**.
- Recommended keys when present: `confidence`, `summary`.
- `kind` is `"codex-home"`; `path` equals the absolute `.codex` directory.
- `totalSize` is an integer ≥ non-DB content payload (`codexHomeNonDBContentBytes` = 1198).
- `breakdown` is a non-empty array; entries include roles covering at least:
  `app-logs-db`, `session-logs`, `cache`, `tmp`, `config` (from the fixture).
- `reclaimable` bools: true for app-logs-db, session-logs, cache, tmp; false for config.
- **`logsDb`**: `path` references `logs_*.sqlite`, `size` > 0, `rows` = 5, `samples` len 1–3
  (newest first; sample fields id/ts/level/target/body).
- `reclaim` is a non-empty array of advice objects.
- `howToPurge` is a non-empty array; entries include officialCommand and removes.
- `howToPurge[].officialCommand` is plain (no `$`, no ANSI); includes `disk-usage-analyser scan`;
  never `rm -rf`; must not treat auth/config as usually-safe primary purge.
- **`howToPurge` logs reclaim** (`assertJSONCodexHowToPurgeLogs`): quit Codex; `logs_2.sqlite`;
  both `mv`/backup (+ wal/shm) and `sqlite3` DELETE+VACUUM; caution on `state_5`/auth/config;
  diagnostic/recreate/regrow notes; plain (no `$`/ANSI/`rm -rf`).
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
	for _, key := range []string{"path", "kind", "totalSize", "breakdown", "reclaim", "howToPurge", "rawCommands", "logsDb"} {
		if payload[key] == nil {
			t.Fatalf("json missing required key %q: %s", key, firstJSONObjectLine(t, resp.Stdout))
		}
	}

	path := jsonStringField(t, payload, "path")
	if path != req.TargetPath {
		t.Fatalf("json path: expected %q, got %q", req.TargetPath, path)
	}
	kind := jsonStringField(t, payload, "kind")
	if kind != "codex-home" {
		t.Fatalf("json kind: expected codex-home, got %q", kind)
	}
	total := jsonInt64Field(t, payload, "totalSize")
	if total < codexHomeNonDBContentBytes {
		t.Fatalf("json totalSize: expected >= %d (non-DB file payload), got %d", codexHomeNonDBContentBytes, total)
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
	for _, want := range []string{"app-logs-db", "session-logs", "cache", "tmp", "config"} {
		if !roles[want] {
			t.Fatalf("breakdown roles missing %q (got %v):\n%s", want, roles, firstJSONObjectLine(t, resp.Stdout))
		}
	}
	assertJSONBreakdownRoleReclaimable(t, breakdown, "app-logs-db", true)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "session-logs", true)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "cache", true)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "tmp", true)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "config", false)

	assertCodexLogsDBJSON(t, payload, codexHomeFixtureLogRows)

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
	assertJSONCodexHowToPurgeLogs(t, howToPurge)

	joinedOC := ""
	joinedRemoves := ""
	for _, step := range howToPurge {
		oc, _ := step["officialCommand"].(string)
		rm, _ := step["removes"].(string)
		joinedOC += oc + "\n"
		joinedRemoves += rm + "\n"
	}
	jl := strings.ToLower(joinedOC + "\n" + joinedRemoves)
	if !strings.Contains(jl, "disk-usage-analyser scan") {
		t.Fatalf("JSON howToPurge must include disk-usage-analyser scan: %s", joinedOC)
	}
	if strings.Contains(jl, "rm -rf") || strings.Contains(jl, "rm -fr") {
		t.Fatalf("JSON howToPurge must not contain rm -rf: %s", joinedOC)
	}
	for _, step := range howToPurge {
		rm, _ := step["removes"].(string)
		rl := strings.ToLower(rm)
		if (strings.Contains(rl, "usually safe") || strings.Contains(rl, "usually-safe")) &&
			(strings.Contains(rl, "auth") || strings.Contains(rl, "config.toml")) {
			t.Fatalf("JSON howToPurge removes must not mark auth/config as usually-safe: %q", rm)
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
