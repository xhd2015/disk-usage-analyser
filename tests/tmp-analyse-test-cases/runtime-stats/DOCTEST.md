# Tmp Analyse — Runtime Stats Tests

Docker/Podman runtime stats: NDJSON parsing, human size parsing, filtering, CLI collection, SSE integration.

## Version

0.0.2

# DSN (Domain Specific Notion)

When scanning Docker or Podman cards, **CollectRuntimeStats** runs
`docker|podman system df --format json`. **ParseSystemDFJSON** reads NDJSON lines,
**ParseHumanSize** converts size strings to bytes, and **FilterRuntimeItems** keeps
only Images and Build Cache. Results attach to the Docker/Podman **location** SSE event
as **runtimeItems**. Missing CLI or daemon errors omit runtime stats gracefully.

## Test Tree

```
runtime-stats/
├── parse-system-df/
│   ├── docker-ndjson/
│   ├── podman-ndjson/
│   └── malformed-line-rejected/
├── parse-human-size/
│   ├── megabytes-decimal/
│   ├── gigabytes-compact/
│   ├── bytes-and-zero/
│   └── invalid-input/
├── filter-runtime-items/
│   └── keeps-images-and-build-cache/
├── collect-runtime-stats/
│   ├── docker-mock-success/
│   ├── podman-mock-success/
│   ├── binary-missing/
│   └── command-fails/
└── sse-runtime-items/
    └── docker-location-event/
```

## How to Run

```sh
doctest vet ./tests/tmp-analyse-test-cases/runtime-stats
doctest test ./tests/tmp-analyse-test-cases/runtime-stats
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
	FixtureFile  string
	HumanSize    string
	Runtime      string
	MockOutput   string
	MockFail     bool
	MockNotFound bool
}

type Response struct {
	RuntimeItems []server.TmpRuntimeItem
	ParsedBytes  int64
	ParseFailed  bool
	CollectFailed bool
	SSEOutput    string
	Err          error
}

func Run(t *testing.T, req *Request) (*Response, error) {
	switch req.Op {
	case "parse-system-df-docker", "parse-system-df-podman", "parse-system-df-malformed":
		return runParseSystemDF(t, req)
	case "parse-human-size-mb", "parse-human-size-gb", "parse-human-size-bytes", "parse-human-size-invalid":
		return runParseHumanSize(req)
	case "filter-runtime-items":
		return runFilterRuntimeItems(req)
	case "collect-runtime-docker", "collect-runtime-podman", "collect-runtime-missing", "collect-runtime-fails":
		return runCollectRuntimeStats(req)
	case "sse-runtime-docker":
		return runSSERuntimeDocker(t, req)
	default:
		t.Fatalf("unknown runtime-stats op: %q", req.Op)
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

func readFixtureFromPath(path string) string {
	data, _ := os.ReadFile(path)
	return string(data)
}

func runParseSystemDF(t *testing.T, req *Request) (*Response, error) {
	output := readFixture(t, req.FixtureFile)
	items, err := server.ParseSystemDFJSON(output)
	if err != nil {
		return &Response{ParseFailed: true, Err: err}, nil
	}
	filtered := server.FilterRuntimeItems(items, "Images", "Build Cache")
	return &Response{RuntimeItems: filtered}, nil
}

func runParseHumanSize(req *Request) (*Response, error) {
	bytes, err := server.ParseHumanSize(req.HumanSize)
	if err != nil {
		return &Response{ParseFailed: true, Err: err}, nil
	}
	return &Response{ParsedBytes: bytes}, nil
}

func runFilterRuntimeItems(req *Request) (*Response, error) {
	output := readFixtureFromPath(req.FixtureFile)
	items, err := server.ParseSystemDFJSON(output)
	if err != nil {
		return nil, err
	}
	filtered := server.FilterRuntimeItems(items, "Images", "Build Cache")
	return &Response{RuntimeItems: filtered}, nil
}

func installMockRuntimeRunner(output string, fail bool, notFound bool) {
	server.SetRuntimeCommandRunner(func(name string, args ...string) ([]byte, error) {
		if notFound {
			return nil, os.ErrNotExist
		}
		if fail {
			return nil, errors.New("daemon unavailable")
		}
		return []byte(output), nil
	})
}

func runCollectRuntimeStats(req *Request) (*Response, error) {
	installMockRuntimeRunner(req.MockOutput, req.MockFail, req.MockNotFound)
	items, err := server.CollectRuntimeStats(req.Runtime)
	if err != nil {
		return &Response{CollectFailed: true, Err: err}, nil
	}
	return &Response{RuntimeItems: items}, nil
}

func runSSERuntimeDocker(t *testing.T, req *Request) (*Response, error) {
	output := readFixture(t, req.FixtureFile)
	installMockRuntimeRunner(output, false, false)

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

	var dockerLoc *server.TmpLocation
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
			if err := json.Unmarshal([]byte(data), &loc); err == nil && loc.Label == "Docker" {
				dockerLoc = &loc
			}
		}
	}
	if dockerLoc == nil {
		return &Response{SSEOutput: string(body)}, nil
	}
	return &Response{RuntimeItems: dockerLoc.RuntimeItems, SSEOutput: "found"}, nil
}
```