# Tmp Worktrees Scan

Backend SSE tests for `GET /api/tmp-worktrees-scan`: git repo discovery, worktree
listing, per-checkout sizing, streaming event order, tilde paths, and client disconnect.

## Version

0.0.2

# DSN (Domain Specific Notion)

The **tmp-worktrees-scan handler** discovers git repositories under the user home
directory, lists worktrees per main repo via `git worktree list --porcelain`, sizes
each checkout, and streams **SSE events** to the frontend. A **repo** event announces
each main repository and carries the **main checkout size** (`size`, `sizeHuman`,
`fileCount`). **worktree** events carry **WorktreeHit** records for **linked**
checkouts only (`isMain=false`); the main checkout is never emitted as a worktree
child. **summary** aggregates repo and linked-worktree counts plus total size;
**done** marks completion. The scan root defaults to `~` with unlimited depth in v1.
Client disconnect via `r.Context().Done()` aborts the scan.

## Decision Tree

```
tmp-worktrees-scan/
├── streaming/
│   ├── worktree-before-done/
│   └── repo-before-worktree/
├── discovery/
│   ├── main-plus-linked/
│   ├── main-only/
│   ├── omit-main-worktree-event/
│   ├── multi-repos/
│   └── no-repos/
├── sizing/
│   └── non-empty-checkout/
├── paths/
│   └── tilde-paths/
└── cancellation/
    └── client-disconnect/
```

## Test Index

| Leaf | Op |
|------|-----|
| streaming/worktree-before-done | worktrees-sse-order |
| streaming/repo-before-worktree | worktrees-sse-order |
| discovery/main-plus-linked | worktrees-scan |
| discovery/main-only | worktrees-scan |
| discovery/omit-main-worktree-event | worktrees-scan |
| discovery/multi-repos | worktrees-scan |
| discovery/no-repos | worktrees-scan |
| sizing/non-empty-checkout | worktrees-scan |
| paths/tilde-paths | worktrees-scan |
| cancellation/client-disconnect | worktrees-disconnect |

## How to Run

```sh
doctest vet ./tests/tmp-worktrees-scan
doctest test ./tests/tmp-worktrees-scan
```

```go
import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"disk-usage-analyser/server"
)

type Request struct {
	Op      string
	HomeDir string
}

type WorktreeRepoRow struct {
	RepoPath  string `json:"repoPath"`
	RepoName  string `json:"repoName"`
	Size      int64  `json:"size"`
	SizeHuman string `json:"sizeHuman"`
	FileCount int64  `json:"fileCount"`
}

type Response struct {
	SSEOutput              string
	EventTypes             []string
	Worktrees              []server.WorktreeHit
	RepoRows               []WorktreeRepoRow
	Repos                  int
	Summary                *server.WorktreeScanSummary
	WorktreeBeforeDone     bool
	RepoBeforeWorktree     bool
	DisconnectAborted      bool
}

func parseRepoRowsFromSSE(body string) []WorktreeRepoRow {
	var rows []WorktreeRepoRow
	lines := strings.Split(body, "\n")
	var currentEvent string
	for _, line := range lines {
		if strings.HasPrefix(line, "event: ") {
			currentEvent = strings.TrimPrefix(line, "event: ")
		}
		if strings.HasPrefix(line, "data: ") && currentEvent == "repo" {
			data := strings.TrimPrefix(line, "data: ")
			var row WorktreeRepoRow
			if err := json.Unmarshal([]byte(data), &row); err == nil && row.RepoPath != "" {
				rows = append(rows, row)
			}
		}
	}
	return rows
}

func buildWorktreesResponse(parsed *server.WorktreesSSEResult) *Response {
	return &Response{
		SSEOutput:  parsed.SSEOutput,
		EventTypes: parsed.EventTypes,
		Worktrees:  parsed.Worktrees,
		RepoRows:   parseRepoRowsFromSSE(parsed.SSEOutput),
		Repos:      parsed.Repos,
		Summary:    parsed.Summary,
	}
}

func runWorktreesScan(req *Request) (*Response, error) {
	handler := http.HandlerFunc(server.HandleTmpWorktreesScan)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	httpReq, err := http.NewRequest("GET", srv.URL+"/api/tmp-worktrees-scan", nil)
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
	parsed, err := server.ParseWorktreesSSE(string(body))
	if err != nil {
		return nil, err
	}
	return buildWorktreesResponse(parsed), nil
}

func runWorktreesSSEOrder(req *Request) (*Response, error) {
	handler := http.HandlerFunc(server.HandleTmpWorktreesScan)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	httpReq, err := http.NewRequest("GET", srv.URL+"/api/tmp-worktrees-scan", nil)
	if err != nil {
		return nil, err
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	parsed, err := server.ParseWorktreesSSE(string(body))
	if err != nil {
		return nil, err
	}
	out := buildWorktreesResponse(parsed)
	out.WorktreeBeforeDone = server.EventBefore(parsed.EventTypes, "worktree", "done")
	out.RepoBeforeWorktree = server.EventBefore(parsed.EventTypes, "repo", "worktree")
	return out, nil
}

func runWorktreesDisconnect(t *testing.T, req *Request) (*Response, error) {
	handler := http.HandlerFunc(server.HandleTmpWorktreesScan)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "GET", srv.URL+"/api/tmp-worktrees-scan", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(resp.Body)
	sawRepo := false
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "event: repo") {
			sawRepo = true
			break
		}
	}
	cancel()
	resp.Body.Close()

	if !sawRepo {
		return &Response{DisconnectAborted: false}, nil
	}
	return &Response{DisconnectAborted: true}, nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Setenv("HOME", req.HomeDir)

	switch req.Op {
	case "worktrees-scan":
		return runWorktreesScan(req)
	case "worktrees-sse-order":
		return runWorktreesSSEOrder(req)
	case "worktrees-disconnect":
		return runWorktreesDisconnect(t, req)
	default:
		t.Fatalf("unknown test op: %q", req.Op)
		return nil, nil
	}
}

func tilde(home, path string) string {
	if path == home {
		return "~"
	}
	if strings.HasPrefix(path, home+string(os.PathSeparator)) {
		return "~" + path[len(home):]
	}
	return path
}
```