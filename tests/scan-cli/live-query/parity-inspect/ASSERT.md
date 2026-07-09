## Expected

- Live exit code 0 with `TOP 2` ranking `huge.bin` then `mid.bin` (root skipped).
- Inspect of the equivalent JSON with `--top 2` yields the same two match paths in the same order.
- Trailing blank line on live stdout.

## Exit Code

- 0

```go
import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"disk-usage-analyser/usagescan"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("live: expected exit 0, got %d (err=%v)\n%s", resp.ExitCode, resp.Err, resp.Stdout)
	}
	live := resp.Stdout
	if !strings.Contains(live, "TOP 2") {
		t.Fatalf("live missing TOP 2:\n%s", live)
	}
	hugePath := filepath.Join(req.FixtureDir, "huge.bin")
	midPath := filepath.Join(req.FixtureDir, "mid.bin")
	liveTop := live[strings.Index(live, "TOP 2"):]
	if !strings.Contains(liveTop, hugePath) || !strings.Contains(liveTop, midPath) {
		t.Fatalf("live TOP missing expected paths:\n%s", live)
	}
	if strings.Index(liveTop, hugePath) > strings.Index(liveTop, midPath) {
		t.Fatalf("live: expected huge.bin before mid.bin:\n%s", liveTop)
	}

	jsonPath := filepath.Join(filepath.Dir(req.FixtureDir), "capture.json")
	var stdout, stderr bytes.Buffer
	code, runErr := usagescan.RunCLI([]string{"--inspect", jsonPath, "--top", "2"}, usagescan.CLIOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if runErr != nil || code != 0 {
		t.Fatalf("inspect parity run: code=%d err=%v stderr=%s stdout=%s", code, runErr, stderr.String(), stdout.String())
	}
	ins := strings.ReplaceAll(stdout.String(), "\r\n", "\n")
	if !strings.Contains(ins, "TOP 2") {
		t.Fatalf("inspect missing TOP 2:\n%s", ins)
	}
	insTop := ins[strings.Index(ins, "TOP 2"):]
	if !strings.Contains(insTop, hugePath) || !strings.Contains(insTop, midPath) {
		t.Fatalf("inspect TOP missing expected paths:\n%s", ins)
	}
	if strings.Index(insTop, hugePath) > strings.Index(insTop, midPath) {
		t.Fatalf("inspect: expected huge.bin before mid.bin:\n%s", insTop)
	}
	if (strings.Index(liveTop, hugePath) < strings.Index(liveTop, midPath)) !=
		(strings.Index(insTop, hugePath) < strings.Index(insTop, midPath)) {
		t.Fatalf("live vs inspect ranking mismatch\nlive:\n%s\ninspect:\n%s", liveTop, insTop)
	}
	stdoutEndsWithBlankLine(t, live)
}
```
