# Tmp Files Analyse — Frontend Test Cases

UI tests for the tmp-analyse page: structure, navigation, scan interaction, breakdown
display, runtime stats sections, live breakdown progress, and Repository Scans
(worktrees, binaries, node_modules, and vendor sections with independent scan
controls, filter/sort/select/delete UX).

## Version

0.0.2

# DSN (Domain Specific Notion)

The **React tmp-analyse page** loads location cards from the initial SSE **locations**
event, then starts a scan via **Start Scan**. During scanning, **progress** SSE events
update card header sizes and (for multi-path cards) individual **breakdown rows** live.
When each location completes, a **location** event finalizes sizes and may attach
**runtimeItems** for Docker/Podman, and **vmInternal** for Podman VM storage on macOS.
The page renders cards with **data-testid** hooks for headings, summary bar, breakdown
table rows, vm-internal section rows, runtime section rows, and scan badges.
Repository Scans sections (worktrees, binaries, node_modules, vendor) each have
independent scan buttons, SSE event streams, and delete endpoints. Tests drive the
page through **playwright-debug** scripts against a locally started Go dev server.
Repository Scans leaves use `ui-automation` and `slow` labels when they depend on
completing a subsection scan.

## Test Tree

```
tmp-analyse-frontend-test-cases/
├── [existing] verify-page-renders … verify-npm-breakdown
├── breakdown-live-progress/
│   └── go-multi-path/
├── runtime-section/
│   ├── docker-after-scan/
│   └── podman-after-scan/
├── podman-vm-internal/
│   └── after-scan/
├── repository-scans/
│   └── renders/
├── worktrees-section/
│   ├── after-scan/
│   ├── left-aligned/
│   ├── filter-under-10m-default/
│   ├── filter-show-under-10m/
│   ├── sort-by-size-desc/
│   └── running-total/
├── worktrees-live-stream/
├── binaries-section/
│   ├── after-scan/
│   ├── left-aligned/
│   ├── filter-under-1m-default/
│   ├── filter-show-under-1m/
│   ├── sort-by-size-desc/
│   ├── select-and-total/
│   ├── delete-selected/
│   ├── repo-select-all/
│   └── running-total/
├── named-section/
│   ├── after-scan/
│   ├── left-aligned/
│   ├── filter-under-1m-default/
│   ├── filter-show-under-1m/
│   ├── sort-by-size-desc/
│   ├── select-and-total/
│   ├── delete-selected/
│   ├── repo-select-all/
│   └── running-total/
├── vendor-section/
│   ├── after-scan/
│   └── running-total/
├── named-independent-scans/
└── independent-scan-controls/
```

## Test Index

| Leaf | Script |
|------|--------|
| verify-page-renders | page-renders.js |
| verify-breakdown-live-progress/go-multi-path | breakdown-live-progress.js |
| verify-runtime-section/docker-after-scan | docker-runtime-section.js |
| verify-runtime-section/podman-after-scan | podman-runtime-section.js |
| podman-vm-internal/after-scan | podman-vm-internal-section.js |
| repository-scans/renders | repository-scans-renders.js |
| worktrees-section/after-scan | worktrees-after-scan.js |
| worktrees-section/left-aligned | worktrees-left-aligned.js |
| worktrees-section/filter-under-10m-default | worktrees-filter-under-10m-default.js |
| worktrees-section/filter-show-under-10m | worktrees-filter-show-under-10m.js |
| worktrees-section/sort-by-size-desc | worktrees-sort-by-size-desc.js |
| worktrees-section/running-total | worktrees-running-total.js |
| worktrees-live-stream | worktrees-live-stream.js |
| binaries-section/after-scan | binaries-after-scan.js |
| binaries-section/left-aligned | binaries-left-aligned.js |
| binaries-section/filter-under-1m-default | binaries-filter-under-1m-default.js |
| binaries-section/filter-show-under-1m | binaries-filter-show-under-1m.js |
| binaries-section/sort-by-size-desc | binaries-sort-by-size-desc.js |
| binaries-section/select-and-total | binaries-select-total.js |
| binaries-section/delete-selected | binaries-delete-selected.js |
| binaries-section/repo-select-all | binaries-repo-select-all.js |
| binaries-section/running-total | binaries-running-total.js |
| named-section/after-scan | named-after-scan.js |
| named-section/left-aligned | named-left-aligned.js |
| named-section/filter-under-1m-default | named-filter-under-1m-default.js |
| named-section/filter-show-under-1m | named-filter-show-under-1m.js |
| named-section/sort-by-size-desc | named-sort-by-size-desc.js |
| named-section/select-and-total | named-select-total.js |
| named-section/delete-selected | named-delete-selected.js |
| named-section/repo-select-all | named-repo-select-all.js |
| named-section/running-total | node-modules-running-total.js |
| vendor-section/after-scan | vendor-after-scan.js |
| vendor-section/running-total | vendor-running-total.js |
| named-independent-scans | named-independent-scans.js |
| independent-scan-controls | independent-scan-controls.js |

## How to Run

```sh
doctest vet ./tests/tmp-analyse-frontend-test-cases
doctest test ./tests/tmp-analyse-frontend-test-cases
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