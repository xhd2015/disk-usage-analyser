# Scenario

**Feature**: filterNamedReposByColumnFilters pure function for node_modules column filters

```
# byRepo map + NamedColumnFilters -> visible hits per repo, drop empty repos
filterNamedReposByColumnFilters(byRepo, filters) -> filtered Map<repoPath, NamedHit[]>
```

## Preconditions

- Node.js is installed and on PATH.
- `disk-usage-analyser-react/src/repositoryScansLayout.ts` exports `filterNamedReposByColumnFilters`.
- Each leaf provides a JSON fixture under `testdata/`.

## Steps

1. Root Setup verifies `node` and `npx` are available.
2. Leaf Setup sets `req.Op` and `req.FixtureFile`.

```go
import (
	"os/exec"
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
			t.Fatalf("item %d: expected string, got %T", i, item)
		}
		out[i] = s
	}
	return out
}

func assertOrder(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("order len %d want %d\ngot:  %v\nwant: %v", len(got), len(want), got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %q want %q\nfull got:  %v\nfull want: %v", i, got[i], want[i], got, want)
		}
	}
}
```