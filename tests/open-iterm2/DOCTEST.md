# Open iTerm2 API

Backend POST tests for `POST /api/open-iterm2`: tilde path resolution, parent-dir
open when path ends with `node_modules`, force-new-window via `iterm2.ModeForceNew`,
validation errors, and platform guard. Default tests capture AppleScript via
`KOOL_ITERM2_SCRIPT_OUT` without live osascript.

## Version

0.0.2

# DSN (Domain Specific Notion)

The **open-iterm2 handler** accepts JSON `{"path":"~/..."}`. It resolves `~` against
the process home directory, normalizes the path, and when the basename is
`node_modules` retargets to the parent project directory. On macOS it calls
`iterm2.OpenConfig` with `ModeForceNew` so iTerm2 always opens a new window and
`cd`s to the target directory. Non-darwin platforms return HTTP 501; missing or
empty path returns HTTP 400.

Test harness sets `KOOL_ITERM2_*` env vars so the shared `shell/iterm2` library
writes the generated AppleScript to a file instead of invoking real osascript.

## Decision Tree

```
open-iterm2/
├── script-capture/
│   └── force-new-window/
├── validation/
│   └── missing-path/
└── platform/
    └── non-darwin/
```

## Test Index

| Leaf | Op |
|------|-----|
| script-capture/force-new-window | open-iterm2 |
| validation/missing-path | open-iterm2 |
| platform/non-darwin | open-iterm2 |

## How to Run

```sh
doctest vet ./tests/open-iterm2
doctest test ./tests/open-iterm2
```

```go
import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"disk-usage-analyser/server"
)

const (
	koolIterm2InstalledEnv     = "KOOL_ITERM2_INSTALLED"
	koolIterm2ScriptOutEnv     = "KOOL_ITERM2_SCRIPT_OUT"
	koolIterm2OsascriptExitEnv = "KOOL_ITERM2_OSASCRIPT_EXIT"
	koolIterm2GOOSEnv          = "KOOL_ITERM2_GOOS"
)

type Request struct {
	Op            string
	HomeDir       string
	Path          string
	BodyJSON      string
	GoOS          string
	Installed     string
	OsascriptExit string
	ScriptOutPath string
}

type Response struct {
	StatusCode     int
	Body           string
	Opened         string
	CapturedScript string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Setenv("HOME", req.HomeDir)

	if req.Installed == "" {
		req.Installed = "1"
	}
	t.Setenv(koolIterm2InstalledEnv, req.Installed)

	goos := req.GoOS
	if goos == "" {
		goos = "darwin"
	}
	t.Setenv(koolIterm2GOOSEnv, goos)

	if req.ScriptOutPath != "" {
		t.Setenv(koolIterm2ScriptOutEnv, req.ScriptOutPath)
	}
	if req.OsascriptExit == "" {
		req.OsascriptExit = "0"
	}
	t.Setenv(koolIterm2OsascriptExitEnv, req.OsascriptExit)

	var payload []byte
	var err error
	if req.BodyJSON != "" {
		payload = []byte(req.BodyJSON)
	} else {
		payload, err = json.Marshal(server.OpenIterm2Request{Path: req.Path})
		if err != nil {
			return nil, err
		}
	}

	handler := http.HandlerFunc(server.HandleOpenIterm2)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	httpReq, err := http.NewRequest("POST", srv.URL+"/api/open-iterm2", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := &Response{
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}
	if resp.StatusCode == http.StatusOK {
		var okBody struct {
			Status string `json:"status"`
			Opened string `json:"opened"`
		}
		if err := json.Unmarshal(body, &okBody); err == nil {
			out.Opened = okBody.Opened
		}
	}
	if req.ScriptOutPath != "" {
		if b, readErr := os.ReadFile(req.ScriptOutPath); readErr == nil {
			out.CapturedScript = string(b)
		}
	}
	return out, nil
}

func tildePath(home, abs string) string {
	if abs == home {
		return "~"
	}
	sep := string(os.PathSeparator)
	if len(abs) > len(home) && abs[:len(home)] == home && abs[len(home)] == sep[0] {
		return "~" + abs[len(home):]
	}
	return abs
}

func projectParentAbs(home, rel string) string {
	return filepath.Join(home, filepath.FromSlash(rel))
}

func nodeModulesTildePath(home, rel string) string {
	parent := projectParentAbs(home, rel)
	return tildePath(home, filepath.Join(parent, "node_modules"))
}

func scriptContainsTargetDir(script, absDir string) bool {
	if strings.Contains(script, absDir) {
		return true
	}
	escaped := strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(absDir)
	return strings.Contains(script, escaped)
}
```