## Expected
- Go location has exactly 1 ExtraPath: `~/Library/Caches/go-build` (tilde-prefixed, full path)
- Xcode location has exactly 1 ExtraPath: `~/Library/Developer/CoreSimulator/Devices` (tilde-prefixed, full path)
- All other single-path software locations have empty ExtraPaths (length 0); multi-path tools include OpenCode, Claude Code, Codex, and Cursor
- ExtraPaths use `~` prefix (not absolute home directory)

```go
import (
	"path/filepath"
	"strings"

	"disk-usage-analyser/server"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Locations) != 21 {
		t.Fatalf("expected 21 software locations, got %d", len(resp.Locations))
	}

	expectedExtraPaths := map[string]string{
		"go":       "~/Library/Caches/go-build",
		"xcode":    "~/Library/Developer/CoreSimulator/Devices",
		"opencode": "~/.local/share/opencode/project",
		"claude":   "~/.claude/telemetry",
		"codex":    "~/Library/Application Support/codex",
		"cursor":   "~/Library/Application Support/Caches/cursor-updater",
	}

	multiPathCats := map[string]int{"go": 1, "xcode": 1, "opencode": 6, "claude": 4, "codex": 1, "cursor": 2}
	for _, loc := range resp.Locations {
		expectedExtraCount, isMulti := multiPathCats[loc.Category]
		if isMulti {
			if len(loc.ExtraPaths) != expectedExtraCount {
				t.Fatalf("location %s: expected %d ExtraPaths, got %d", loc.Category, expectedExtraCount, len(loc.ExtraPaths))
			}
			expected := expectedExtraPaths[loc.Category]
			if loc.ExtraPaths[0] != expected {
				t.Fatalf("location %s: expected ExtraPaths[0]=%s, got %s", loc.Category, expected, loc.ExtraPaths[0])
			}
		} else {
			if len(loc.ExtraPaths) != 0 {
				t.Fatalf("location %s: expected 0 ExtraPaths, got %d", loc.Category, len(loc.ExtraPaths))
			}
		}
	}

	// Also verify ExtraPaths were captured in the response
	if len(resp.ExtraPaths) < 2 {
		t.Fatalf("expected at least 2 ExtraPaths (one for Go, one for Xcode), got %d", len(resp.ExtraPaths))
	}
}
```
