# Xcode Card — Expanded Discovery and Simulator Runtimes

Backend tests for expanded Xcode ExtraPaths, simulator runtime JSON parsing,
`CollectSimulatorRuntimeStats`, and SSE `runtimeItems` on the Xcode location event.

## Version

0.0.2

# DSN (Domain Specific Notion)

**DiscoverLocations** registers the Xcode card with primary DerivedData and four
ordered ExtraPaths (CoreSimulator devices, iOS DeviceSupport, Archives,
DocumentationCache). **CollectSimulatorRuntimeStats** runs
`xcrun simctl runtime list -j` through an injectable command runner, parses the
UUID-keyed JSON (`version`, `identifier`, `mountPath`, `deletable`, `state`),
computes on-disk size per `mountPath`, and builds **TmpRuntimeItem** rows
(`Type`, `TotalCount` 1, `ActiveCount` 1 when `state == "Ready"`, `Reclaimable`
equals size when `deletable`). **HandleTmpAnalyse** attaches non-empty
`runtimeItems` on the Xcode **location** SSE event. Command failure or missing
`simctl` returns empty stats gracefully (nil error).

## Test Tree

```
xcode-card/
├── discover-extra-paths/
├── parse-simulator-runtime-json/
├── collect-simulator-runtimes/
│   ├── mock-success/
│   └── command-fails/
├── sse-xcode-runtime-items/
└── frontend/DOCTEST.md          # nested root — UI leaves
```

## How to Run

```sh
doctest vet ./tests/xcode-card
doctest test ./tests/xcode-card
doctest test ./tests/xcode-card/frontend
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
	Op           string
	HomeDir      string
	FixtureFile  string
	MockOutput   string
	MockFail     bool
	MockNotFound bool
	MountSizes   map[string]int64
}

type Response struct {
	Locations     []server.TmpLocation
	XcodeLoc      *server.TmpLocation
	RuntimeItems  []server.TmpRuntimeItem
	CollectFailed bool
	SSEOutput     string
	Err           error
}

func Run(t *testing.T, req *Request) (*Response, error) {
	switch req.Op {
	case "discover-extra-paths":
		return runDiscoverExtraPaths(req)
	case "parse-simulator-runtime-json":
		return runParseSimulatorRuntimeJSON(t, req)
	case "collect-simulator-mock-success":
		return runCollectSimulatorMockSuccess(t, req)
	case "collect-simulator-fails":
		return runCollectSimulatorFails(req)
	case "sse-xcode-runtime":
		return runSseXcodeRuntime(t, req)
	default:
		t.Fatalf("unknown xcode-card op: %q", req.Op)
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

func installMockSimulatorRunner(output string, fail bool, notFound bool) {
	server.SetSimulatorRuntimeCommandRunner(func(name string, args ...string) ([]byte, error) {
		if notFound {
			return nil, os.ErrNotExist
		}
		if fail {
			return nil, errors.New("simctl unavailable")
		}
		return []byte(output), nil
	})
}

func runDiscoverExtraPaths(req *Request) (*Response, error) {
	locations := server.DiscoverLocations(req.HomeDir)
	var xcodeLoc *server.TmpLocation
	for i := range locations {
		if locations[i].Category == "xcode" {
			xcodeLoc = &locations[i]
			break
		}
	}
	return &Response{Locations: locations, XcodeLoc: xcodeLoc}, nil
}

func runParseSimulatorRuntimeJSON(t *testing.T, req *Request) (*Response, error) {
	output := readFixture(t, req.FixtureFile)
	sizeFn := func(mountPath string) int64 {
		if req.MountSizes != nil {
			if size, ok := req.MountSizes[mountPath]; ok {
				return size
			}
		}
		return 0
	}
	items, err := server.ParseSimulatorRuntimeJSON(output, sizeFn)
	if err != nil {
		return &Response{Err: err}, err
	}
	return &Response{RuntimeItems: items}, nil
}

func runCollectSimulatorMockSuccess(t *testing.T, req *Request) (*Response, error) {
	output := req.MockOutput
	if output == "" {
		output = readFixture(t, req.FixtureFile)
	}
	installMockSimulatorRunner(output, false, false)
	items, err := server.CollectSimulatorRuntimeStats()
	if err != nil {
		return &Response{CollectFailed: true, Err: err}, nil
	}
	return &Response{RuntimeItems: items}, nil
}

func runCollectSimulatorFails(req *Request) (*Response, error) {
	installMockSimulatorRunner("", true, false)
	items, err := server.CollectSimulatorRuntimeStats()
	if err != nil {
		return &Response{CollectFailed: true, Err: err}, nil
	}
	return &Response{RuntimeItems: items}, nil
}

func runSseXcodeRuntime(t *testing.T, req *Request) (*Response, error) {
	output := req.MockOutput
	if output == "" {
		output = readFixture(t, req.FixtureFile)
	}
	installMockSimulatorRunner(output, false, false)

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

	var xcodeLoc *server.TmpLocation
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
			if err := json.Unmarshal([]byte(data), &loc); err == nil && loc.Category == "xcode" {
				xcodeLoc = &loc
			}
		}
	}
	if xcodeLoc == nil {
		return &Response{SSEOutput: string(body)}, nil
	}
	return &Response{RuntimeItems: xcodeLoc.RuntimeItems, XcodeLoc: xcodeLoc, SSEOutput: "found"}, nil
}
```