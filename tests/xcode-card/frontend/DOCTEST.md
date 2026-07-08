# Xcode Card — Frontend UI Tests

Playwright UI tests for expanded Xcode breakdown rows, simulator runtime section,
and extended cleanup tooltip on `/tmp-analyse`.

## Version

0.0.2

# DSN (Domain Specific Notion)

The **React tmp-analyse page** renders the Xcode card from SSE **locations** and
**location** events. Five **breakdown-items** rows show DerivedData plus four
ExtraPaths. When **runtimeItems** are present on the xcode location, the existing
**runtime-section** renders per-runtime rows. The cleanup popover for **xcode**
lists DerivedData, iOS DeviceSupport, Archives, DocumentationCache, simulator
device commands, and simulator runtime delete-by-UUID guidance.

## Test Tree

```
frontend/
├── breakdown-five-rows/
├── runtime-section/
└── cleanup-popover/
```

## How to Run

```sh
doctest vet ./tests/xcode-card/frontend
doctest test ./tests/xcode-card/frontend
doctest test --label ui-automation ./tests/xcode-card/frontend
```

```go
import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

type Request struct {
	ScriptFile string
}

type Response struct {
	Output string
	Passed bool
}

func Run(t *testing.T, req *Request) (*Response, error) {
	script, err := os.ReadFile(req.ScriptFile)
	if err != nil {
		return nil, err
	}

	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to find project root: %v", err)
	}

	serverCmd := exec.Command("go", "run", ".", "--dev")
	serverCmd.Dir = projectRoot
	serverCmd.Env = append(os.Environ(), "NO_BROWSER=1")

	stdout, err := serverCmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %v", err)
	}

	if err := serverCmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start server: %v", err)
	}

	defer func() {
		serverCmd.Process.Signal(syscall.SIGTERM)
		serverCmd.Wait()
	}()

	portCh := make(chan string, 1)
	errCh := make(chan error, 1)

	reader := bufio.NewReader(stdout)
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					errCh <- err
				}
				return
			}
			if p, ok := parsePort(line); ok {
				portCh <- p
				return
			}
		}
	}()

	var port string
	select {
	case port = <-portCh:
	case err := <-errCh:
		serverCmd.Process.Signal(syscall.SIGTERM)
		serverCmd.Wait()
		return nil, fmt.Errorf("failed to read server output: %v", err)
	case <-time.After(60 * time.Second):
		serverCmd.Process.Signal(syscall.SIGTERM)
		serverCmd.Wait()
		return nil, fmt.Errorf("server did not start within 60 seconds")
	}

	go func() {
		io.Copy(io.Discard, stdout)
	}()

	serverURL := fmt.Sprintf("http://localhost:%s", port)

	cmd := exec.Command("playwright-debug", string(script))
	cmd.Env = append(os.Environ(), "SERVER_URL="+serverURL, "NO_BROWSER=1")
	out, runErr := cmd.CombinedOutput()

	passed := cmd.ProcessState != nil && cmd.ProcessState.ExitCode() == 0

	return &Response{
		Output: string(out),
		Passed: passed,
	}, runErr
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found in any parent directory")
		}
		dir = parent
	}
}

func parsePort(line string) (string, bool) {
	prefix := "Serving directory preview at http://localhost:"
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	rest := strings.TrimPrefix(line, prefix)
	return strings.TrimSpace(rest), true
}
```