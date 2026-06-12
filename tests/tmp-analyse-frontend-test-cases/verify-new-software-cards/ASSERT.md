## Expected
- Cards appear for opencode, claude, codex (all detected on this macOS system)
- Cursor card may not appear (AppSupport/Cursor not present on this system)
- Each detected card has: card-label, card-size, reboot-safe-badge
- Multi-path tools (opencode, claude, codex, cursor when detected) show breakdown-items wrapper
- Multi-path tools do NOT have standalone card-path
- OpenCode has at least 4 breakdown rows (snapshot, project, tool-output, storage, log, cache, state)
- Claude Code has at least 3 breakdown rows (plugins, telemetry, todos, cache, backups)
- Existing 17 software cards remain unaffected

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("playwright-debug failed: %v\nOutput:\n%s", err, resp.Output)
	}

	// Tools known to be installed on this macOS system (verified by path existence)
	expectedDetected := []string{"opencode", "claude", "codex"}
	// Cursor may not be detected (AppSupport/Cursor dir missing)

	for _, cat := range expectedDetected {
		foundLine := findLine(resp.Output, "CARD_FOUND "+cat)
		if !strings.Contains(foundLine, "true") {
			t.Fatalf("expected card-%s to be detected", cat)
		}

		// Must have label
		labelLine := findLine(resp.Output, "ELEM card-"+cat+"-label")
		if labelLine == "" || strings.Contains(labelLine, "MISSING") {
			t.Fatalf("expected card-%s to have label", cat)
		}

		// Must have size
		sizeLine := findLine(resp.Output, "ELEM card-"+cat+"-size")
		if sizeLine == "" || strings.Contains(sizeLine, "MISSING") {
			t.Fatalf("expected card-%s to have size", cat)
		}

		// Must have reboot-safe badge
		rebootLine := findLine(resp.Output, "REBOOT_SAFE "+cat)
		if strings.Contains(rebootLine, "MISSING") {
			t.Fatalf("expected card-%s to have reboot-safe badge", cat)
		}

		// Must have breakdown-items (all new tools are multi-path)
		bdLine := findLine(resp.Output, "HAS_BREAKDOWN "+cat)
		if !strings.Contains(bdLine, "true") {
			t.Fatalf("expected card-%s to have breakdown-items", cat)
		}

		// Must NOT have standalone card-path
		cpLine := findLine(resp.Output, "HAS_CARD_PATH "+cat)
		if !strings.Contains(cpLine, "false") {
			t.Fatalf("expected card-%s to NOT have standalone card-path", cat)
		}
	}

	// OpenCode: at least 4 breakdown rows
	opencodeFound := findLine(resp.Output, "CARD_FOUND opencode")
	if strings.Contains(opencodeFound, "true") {
		if !strings.Contains(resp.Output, "BREAKDOWN_ROW opencode-row-0: EXISTS") {
			t.Fatal("expected opencode breakdown-row-0 to exist")
		}
		if !strings.Contains(resp.Output, "BREAKDOWN_ROW opencode-row-3: EXISTS") {
			t.Fatal("expected opencode breakdown-row-3 to exist (at least 4 rows)")
		}
	}

	// Claude: at least 3 breakdown rows
	claudeFound := findLine(resp.Output, "CARD_FOUND claude")
	if strings.Contains(claudeFound, "true") {
		if !strings.Contains(resp.Output, "BREAKDOWN_ROW claude-row-0: EXISTS") {
			t.Fatal("expected claude breakdown-row-0 to exist")
		}
		if !strings.Contains(resp.Output, "BREAKDOWN_ROW claude-row-2: EXISTS") {
			t.Fatal("expected claude breakdown-row-2 to exist (at least 3 rows)")
		}
	}

	// Codex: must be detected
	codexFound := findLine(resp.Output, "CARD_FOUND codex")
	if !strings.Contains(codexFound, "true") {
		t.Fatal("expected codex card to be detected")
	}
}

func findLine(output, prefix string) string {
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			return line
		}
	}
	return ""
}
```
