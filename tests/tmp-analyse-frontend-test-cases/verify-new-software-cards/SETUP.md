# Scenario

**Feature**: new-software-cards UI verification

```
# React page renders tmp-analyse cards and scan UI
User -> /tmp-analyse -> playwright-debug checks data-testid elements
```

## Preconditions
- Four new software locations added: OpenCode, Claude Code, Codex (OpenAI), Cursor
- Each has specific sub-path breakdowns (multi-path via breakdownItems)
- Cards appear based on path existence (detected field)

## Steps
1. Set req.ScriptFile to "new-software-cards.js"
2. The script navigates to /tmp-analyse and checks the new cards
3. Verifies multi-path breakdowns for opencode, claude, codex, cursor
4. Verifies single-path detected tools still appear

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "new-software-cards.js"
	return nil
}
```
