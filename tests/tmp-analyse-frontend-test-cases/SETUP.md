## Preconditions
- The test harness starts its own Go server, no external server needed
- `playwright-debug` CLI tool is installed and on PATH

## Steps
1. Read the playwright JS script from req.ScriptFile
2. Locate project root by walking up to find `go.mod`
3. Start the Go server in `--dev` mode with `NO_BROWSER=1`
4. Parse the port from "Serving directory preview at http://localhost:{port}"
5. Run `playwright-debug` with `SERVER_URL=http://localhost:{port}`
6. Kill the server with SIGTERM on cleanup via defer
7. Return stdout and exit status in the response

## Context
- Each playwright script uses `data-testid` attributes to locate DOM elements
- Scripts log structured output: `ELEM <name>: <value or MISSING>`, `PASS`, `FAIL`
- Scripts throw on fatal errors (navigation failure, etc.)
- The server startup may take up to 60 seconds (Vite dev server initialization)

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
	"time"
)

type Request struct {
	ScriptFile string
}

type Response struct {
	Output string
	Passed bool
}

func Setup(t *testing.T, req *Request) error {
	_, err := exec.LookPath("playwright-debug")
	if err != nil {
		return err
	}
	return nil
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
```
