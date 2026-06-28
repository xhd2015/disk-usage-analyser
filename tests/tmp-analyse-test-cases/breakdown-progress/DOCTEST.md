# Tmp Analyse — Breakdown Progress Tests

Live breakdown progress during multi-path scans: accumulated card totals, row-level updates, npm paths.

## Version

0.0.2

# DSN (Domain Specific Notion)

For multi-breakdown cards (Go, Xcode, npm), **ScanWithProgress** emits **progress** SSE
events shaped by **BuildBreakdownProgressPayload**. The card header **size** reflects
the sum of completed breakdown rows plus the active row partial. Fields **breakdownIndex**,
**breakdownSize**, **breakdownFileCount**, and **breakdownPath** (npm) let the frontend
update individual breakdown rows before the final **location** event.

## Test Tree

```
breakdown-progress/
├── build-payload/
│   ├── accumulated-card-total/
│   ├── first-row-active/
│   └── second-row-frozen-first/
├── multi-path-simulation/
│   ├── card-total-never-drops/
│   └── completed-rows-retained/
└── npm-breakdown-path/
    └── progress-includes-path/
```

## How to Run

```sh
doctest vet ./tests/tmp-analyse-test-cases/breakdown-progress
doctest test ./tests/tmp-analyse-test-cases/breakdown-progress
```

```go
import (
	"testing"

	"disk-usage-analyser/server"
)

type breakdownSimStep struct {
	CompletedSizes []int64
	ActiveIndex    int
	ActiveSize     int64
	ActiveCount    int64
}

type Request struct {
	Op                     string
	Label                  string
	Reclaimable            bool
	CompletedSizes         []int64
	CompletedCounts        []int64
	ActiveIndex            int
	ActiveSize             int64
	ActiveCount            int64
	BreakdownPath          string
	AccumulatedSize        int64
	AccumulatedReclaimable int64
	SimSequence            []breakdownSimStep
}

type Response struct {
	ProgressPayload  map[string]interface{}
	ProgressSequence []map[string]interface{}
	Size             int64
	FileCount        int64
}

func Run(t *testing.T, req *Request) (*Response, error) {
	switch req.Op {
	case "breakdown-payload-accumulated", "breakdown-payload-first-row", "breakdown-payload-second-row":
		return runBuildBreakdownPayload(req)
	case "breakdown-sim-never-drops", "breakdown-sim-rows-retained":
		return runBreakdownSimulation(req)
	case "breakdown-npm-path":
		return runBreakdownNpmPath(req)
	default:
		t.Fatalf("unknown breakdown-progress op: %q", req.Op)
		return nil, nil
	}
}

func runBuildBreakdownPayload(req *Request) (*Response, error) {
	payload := server.BuildBreakdownProgressPayload(
		req.Label,
		req.CompletedSizes,
		req.CompletedCounts,
		req.ActiveIndex,
		req.ActiveSize,
		req.ActiveCount,
		req.AccumulatedSize,
		req.AccumulatedReclaimable,
		req.Reclaimable,
		req.BreakdownPath,
	)
	return &Response{ProgressPayload: payload}, nil
}

func boolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func runBreakdownSimulation(req *Request) (*Response, error) {
	var seq []map[string]interface{}
	var prevCardTotal int64 = -1
	neverDrops := true

	for _, step := range req.SimSequence {
		payload := server.BuildBreakdownProgressPayload(
			req.Label,
			step.CompletedSizes,
			nil,
			step.ActiveIndex,
			step.ActiveSize,
			step.ActiveCount,
			req.AccumulatedSize,
			req.AccumulatedReclaimable,
			req.Reclaimable,
			"",
		)
		seq = append(seq, payload)
		cardTotal := payload["size"].(int64)
		if prevCardTotal >= 0 && cardTotal < prevCardTotal {
			neverDrops = false
		}
		prevCardTotal = cardTotal
	}

	return &Response{
		ProgressSequence: seq,
		Size:             prevCardTotal,
		FileCount:        boolToInt64(neverDrops),
	}, nil
}

func runBreakdownNpmPath(req *Request) (*Response, error) {
	payload := server.BuildBreakdownProgressPayload(
		"npm",
		nil,
		nil,
		0,
		1024,
		5,
		0,
		0,
		true,
		req.BreakdownPath,
	)
	return &Response{ProgressPayload: payload}, nil
}
```