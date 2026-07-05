# Scenario

**Feature**: open-iterm2 shared fixture harness

```
POST /api/open-iterm2 {"path":"~/..."} -> resolve ~ -> optional node_modules parent -> iterm2.OpenConfig ModeForceNew
```

## Preconditions

- Tests must never use the real user home directory.
- Default leaves capture AppleScript via `KOOL_ITERM2_SCRIPT_OUT` (no live iTerm2).
- `KOOL_ITERM2_INSTALLED=1` and `KOOL_ITERM2_GOOS=darwin` unless testing platform guard.

## Steps

1. Create a fresh fixture home for each leaf.
2. Set `req.Op` to `open-iterm2` and configure path / env overrides.
3. POST to httptest handler and read captured script when applicable.

## Context

- `req.ScriptOutPath` is set per leaf for script-capture scenarios.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.HomeDir = filepath.Join(t.TempDir(), "home")
	if err := os.MkdirAll(req.HomeDir, 0755); err != nil {
		return err
	}
	if req.Op == "" {
		req.Op = "open-iterm2"
	}
	return nil
}

func writeDirFixture(t *testing.T, home, rel string) string {
	t.Helper()
	dir := filepath.Join(home, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Join(dir, "node_modules", "pkg"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "node_modules", "pkg", "dep.txt"), []byte("fixture\n"), 0644); err != nil {
		t.Fatal(err)
	}
	return dir
}
```