## Expected

- Exit code 0.
- Stdout is one JSON object (plus trailing blank line).
- Required keys: `path`, `kind`, `totalSize`, `breakdown`, `reclaim`, `howToPurge`, `rawCommands`.
- Recommended keys when present: `confidence`, `summary`.
- `kind` is `"android-sdk"`; `path` equals the absolute SDK directory.
- `totalSize` is an integer ≥ content payload (`androidSDKContentBytes` = 890).
- `breakdown` is a non-empty array; entries include roles covering at least:
  `system-images`, `emulator`, `sources`, `build-tools`, `platform-tools`, `platforms`, `tmp`.
- `reclaimable` bools: true for system-images, sources, tmp; false for emulator, build-tools,
  platform-tools, platforms.
- `reclaim` is a non-empty array of advice objects.
- `howToPurge` is a non-empty array; entries include officialCommand and removes.
- `howToPurge[].officialCommand` is plain (no `$`, no ANSI); includes `sdkmanager` and/or
  `disk-usage-analyser scan`; never `rm -rf`.
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
	if kind != "android-sdk" {
		t.Fatalf("json kind: expected android-sdk, got %q", kind)
	}
	total := jsonInt64Field(t, payload, "totalSize")
	if total < androidSDKContentBytes {
		t.Fatalf("json totalSize: expected >= %d (file payload), got %d", androidSDKContentBytes, total)
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
	for _, want := range []string{"system-images", "emulator", "sources", "build-tools", "platform-tools", "platforms", "tmp"} {
		if !roles[want] {
			t.Fatalf("breakdown roles missing %q (got %v):\n%s", want, roles, firstJSONObjectLine(t, resp.Stdout))
		}
	}
	assertJSONBreakdownRoleReclaimable(t, breakdown, "system-images", true)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "sources", true)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "tmp", true)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "emulator", false)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "build-tools", false)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "platform-tools", false)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "platforms", false)

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

	joinedOC := ""
	for _, step := range howToPurge {
		oc, _ := step["officialCommand"].(string)
		joinedOC += oc + "\n"
	}
	jl := strings.ToLower(joinedOC)
	if !strings.Contains(jl, "sdkmanager") && !strings.Contains(jl, "disk-usage-analyser scan") {
		t.Fatalf("JSON howToPurge must include sdkmanager and/or disk-usage-analyser scan: %s", joinedOC)
	}
	if strings.Contains(jl, "rm -rf") || strings.Contains(jl, "rm -fr") {
		t.Fatalf("JSON howToPurge must not contain rm -rf: %s", joinedOC)
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
