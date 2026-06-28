# Tmp Analyse — Podman VM Internal Storage Tests

Podman machine VM storage: `du -sb` parsing, SSH collection, inner-vm runtime stats, SSE integration.

## Version

0.0.2

# DSN (Domain Specific Notion)

On macOS, the **Podman card** scans host `~/.local/share/containers` for the card
header total. Separately, **CollectPodmanVmInternal** checks whether the **Podman
machine** is running via `podman machine list`, then SSHes into the VM to run
`du -sb` on container storage paths. **ParseDuSBOutput** extracts byte totals from
`du` stdout (tab-separated, concatenated, or noisy). Results populate **TmpVmInternal**
with labeled storage rows. **CollectPodmanRuntimeViaSSH** runs `podman system df`
inside the VM; existing **ParseSystemDFJSON** and **FilterRuntimeItems** produce
**runtimeItems**. Both attach to the Podman **location** SSE event. Collection never
fails the overall scan — missing machine, SSH errors, or non-darwin platforms omit
VM stats gracefully.

## Test Tree

```
podman-vm-internal/
├── parse-du-output/
│   ├── tab-separated/
│   ├── concatenated/
│   ├── with-warnings/
│   └── invalid-input/
├── collect-vm-internal/
│   ├── machine-running/
│   │   ├── storage-du-mock-success/
│   │   └── ssh-fails/
│   ├── machine-stopped/
│   └── non-darwin-skipped/
├── collect-runtime-via-ssh/
│   ├── mock-success/
│   └── ssh-fails-graceful/
└── sse-podman-location/
    └── location-event-includes-vm-internal/
```

## Test Index

| Leaf | Op |
|------|-----|
| parse-du-output/tab-separated | parse-du-tab |
| parse-du-output/concatenated | parse-du-concat |
| parse-du-output/with-warnings | parse-du-warnings |
| parse-du-output/invalid-input | parse-du-invalid |
| collect-vm-internal/machine-running/storage-du-mock-success | collect-vm-storage-success |
| collect-vm-internal/machine-running/ssh-fails | collect-vm-ssh-fails |
| collect-vm-internal/machine-stopped | collect-vm-machine-stopped |
| collect-vm-internal/non-darwin-skipped | collect-vm-non-darwin |
| collect-runtime-via-ssh/mock-success | collect-runtime-ssh-success |
| collect-runtime-via-ssh/ssh-fails-graceful | collect-runtime-ssh-fails |
| sse-podman-location/location-event-includes-vm-internal | sse-podman-location |

## How to Run

```sh
doctest vet ./tests/tmp-analyse-test-cases/podman-vm-internal
doctest test ./tests/tmp-analyse-test-cases/podman-vm-internal
```

```go
import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"disk-usage-analyser/server"
)

type Request struct {
	Op                string
	DuOutput          string
	FixtureFile       string
	MachineListJSON   string
	MachineRunning    bool
	MockSSHOutputs    map[string]string
	MockSSHFail       bool
	MockRuntimeOutput string
	ForceGOOS         string
}

type Response struct {
	ParsedBytes   int64
	ParseFailed   bool
	VmInternal    *server.TmpVmInternal
	RuntimeItems  []server.TmpRuntimeItem
	CollectFailed bool
	SSEOutput     string
	ScanHadError  bool
	Err           error
}

func Run(t *testing.T, req *Request) (*Response, error) {
	switch req.Op {
	case "parse-du-tab", "parse-du-concat", "parse-du-warnings", "parse-du-invalid":
		return runParseDuSB(req)
	case "collect-vm-storage-success", "collect-vm-ssh-fails", "collect-vm-machine-stopped", "collect-vm-non-darwin":
		return runCollectVmInternal(req)
	case "collect-runtime-ssh-success", "collect-runtime-ssh-fails":
		return runCollectRuntimeViaSSH(req)
	case "sse-podman-location":
		return runSSEPodmanLocation(t, req)
	default:
		t.Fatalf("unknown podman-vm-internal op: %q", req.Op)
		return nil, nil
	}
}

func readFixture(t *testing.T, name string) string {
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return string(data)
}

func runParseDuSB(req *Request) (*Response, error) {
	output := req.DuOutput
	if output == "" && req.FixtureFile != "" {
		data, err := os.ReadFile(req.FixtureFile)
		if err != nil {
			return nil, err
		}
		output = string(data)
	}
	bytes, err := server.ParseDuSBOutput(output)
	if err != nil {
		return &Response{ParseFailed: true, Err: err}, nil
	}
	return &Response{ParsedBytes: bytes}, nil
}

func installMockPodmanMachineRunner(req *Request) {
	server.SetPodmanMachineRunner(func(args ...string) ([]byte, error) {
		if req.MockSSHFail {
			return nil, errors.New("ssh connection refused")
		}
		joined := strings.Join(args, " ")
		if strings.Contains(joined, "machine list") {
			return []byte(req.MachineListJSON), nil
		}
		if strings.Contains(joined, "system df") {
			return []byte(req.MockRuntimeOutput), nil
		}
		var matchedPath string
		var matchedOut string
		for path, out := range req.MockSSHOutputs {
			if strings.Contains(joined, path) && len(path) > len(matchedPath) {
				matchedPath = path
				matchedOut = out
			}
		}
		if matchedPath != "" {
			return []byte(matchedOut), nil
		}
		return nil, errors.New("unexpected podman machine command: " + joined)
	})
}

func runCollectVmInternal(req *Request) (*Response, error) {
	if req.ForceGOOS != "" {
		server.SetPodmanVmGOOSOverride(req.ForceGOOS)
	}
	installMockPodmanMachineRunner(req)
	vm, err := server.CollectPodmanVmInternal()
	if err != nil {
		return &Response{CollectFailed: true, Err: err}, nil
	}
	return &Response{VmInternal: vm}, nil
}

func runCollectRuntimeViaSSH(req *Request) (*Response, error) {
	installMockPodmanMachineRunner(req)
	items, err := server.CollectPodmanRuntimeViaSSH()
	if err != nil {
		return &Response{CollectFailed: true, Err: err}, nil
	}
	return &Response{RuntimeItems: items}, nil
}

func runSSEPodmanLocation(t *testing.T, req *Request) (*Response, error) {
	installMockPodmanMachineRunner(req)
	server.SetRuntimeCommandRunner(func(name string, args ...string) ([]byte, error) {
		return nil, errors.New("host runtime should not be used for podman on darwin")
	})

	handler := http.HandlerFunc(server.HandleTmpAnalyse)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	httpReq, err := http.NewRequest("GET", srv.URL+"/api/tmp-analyse", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var podmanLoc *server.TmpLocation
	scanner := bufio.NewScanner(strings.NewReader(string(body)))
	var eventType string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "event: ") {
			eventType = strings.TrimPrefix(line, "event: ")
		}
		if strings.HasPrefix(line, "data: ") && eventType == "location" {
			var loc server.TmpLocation
			data := strings.TrimPrefix(line, "data: ")
			if err := json.Unmarshal([]byte(data), &loc); err == nil && loc.Category == "podman" {
				podmanLoc = &loc
			}
		}
	}
	if podmanLoc == nil {
		return &Response{SSEOutput: string(body)}, nil
	}
	return &Response{
		VmInternal:   podmanLoc.VmInternal,
		RuntimeItems: podmanLoc.RuntimeItems,
		SSEOutput:    "found",
	}, nil
}
```