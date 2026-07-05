# Scenario

**Decision**: AppleScript capture via KOOL_ITERM2 env (no live osascript)

```
POST node_modules path -> OpenConfig ModeForceNew -> script written to KOOL_ITERM2_SCRIPT_OUT
```

## Preconditions

- `KOOL_ITERM2_OSASCRIPT_EXIT=0` for success path.
- Captured script must use force-new-window mode (`create window`).

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
	req.Op = "open-iterm2"
	req.GoOS = "darwin"
	req.Installed = "1"
	req.OsascriptExit = "0"
	req.ScriptOutPath = filepath.Join(t.TempDir(), "captured.applescript")
	return nil
}
```