# Scenario

**Feature**: repository scans layout pure functions (sort + filter)

```
# SSE events accumulate in React state -> repositoryScansLayout filter/sort -> render
TmpFilesAnalyse <- repo/worktree/binary SSE -> repositoryScansLayout -> sorted filtered rows
```

## Preconditions

- Node.js is installed and on PATH.
- `disk-usage-analyser-react/src/repositoryScansLayout.ts` exports layout helpers.
- Each leaf provides a JSON fixture under `testdata/`.

## Steps

1. Root Setup verifies `node` and `npx` are available.
2. Leaf Setup sets `req.Op` and `req.FixtureFile`.

## Context

- Harness: `layout-harness.ts` invoked via `npx --yes tsx`.
- Constants: `ONE_MIB = 1048576`, `TEN_MIB = 10485760`.

```go
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if _, err := exec.LookPath("node"); err != nil {
		return err
	}
	if _, err := exec.LookPath("npx"); err != nil {
		return err
	}
	if req.FixtureFile == "" {
		req.FixtureFile = "testdata/fixture.json"
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
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func jsonStringSlice(t *testing.T, v interface{}) []string {
	t.Helper()
	arr, ok := v.([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", v)
	}
	out := make([]string, len(arr))
	for i, item := range arr {
		s, ok := item.(string)
		if !ok {
			t.Fatalf("expected string at index %d, got %T", i, item)
		}
		out[i] = s
	}
	return out
}

func jsonIntSlice(t *testing.T, v interface{}) []int64 {
	t.Helper()
	arr, ok := v.([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", v)
	}
	out := make([]int64, len(arr))
	for i, item := range arr {
		switch n := item.(type) {
		case float64:
			out[i] = int64(n)
		case int64:
			out[i] = n
		default:
			t.Fatalf("expected number at index %d, got %T", i, item)
		}
	}
	return out
}

func assertMonotonicDesc(t *testing.T, sizes []int64) {
	t.Helper()
	for i := 1; i < len(sizes); i++ {
		if sizes[i] > sizes[i-1] {
			t.Fatalf("sizes not DESC at index %d: %d > %d", i, sizes[i], sizes[i-1])
		}
	}
}

func assertOrder(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("order length %d, want %d\ngot:  %s\nwant: %s", len(got), len(want), strings.Join(got, ", "), strings.Join(want, ", "))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %q want %q\nfull got:  %s\nfull want: %s", i, got[i], want[i], strings.Join(got, ", "), strings.Join(want, ", "))
		}
	}
}
```