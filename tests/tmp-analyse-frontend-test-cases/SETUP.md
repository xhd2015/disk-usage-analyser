## Preconditions
- A Go server is running serving the React app (set `SERVER_URL` env var or default `http://localhost:8080`)
- `playwright-debug` CLI tool is installed and on PATH

## Steps
1. Read the playwright JS script from req.ScriptFile
2. Run `playwright-debug` with the script content
3. Pass `SERVER_URL` as env var to the playwright process
4. Return stdout and exit status in the response

## Context
- Each playwright script uses `data-testid` attributes to locate DOM elements
- Scripts log structured output: `ELEM <name>: <value or MISSING>`, `PASS`, `FAIL`
- Scripts throw on fatal errors (navigation failure, etc.)

```go
import (
	"os"
	"os/exec"
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

func Run(t *testing.T, req *Request) (*Response, error) {
	script, err := os.ReadFile(req.ScriptFile)
	if err != nil {
		return nil, err
	}

	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		serverURL = "http://localhost:8080"
	}

	cmd := exec.Command("playwright-debug", string(script))
	cmd.Env = append(os.Environ(), "SERVER_URL="+serverURL)
	out, runErr := cmd.CombinedOutput()

	passed := cmd.ProcessState != nil && cmd.ProcessState.ExitCode() == 0

	return &Response{
		Output: string(out),
		Passed: passed,
	}, runErr
}
```
