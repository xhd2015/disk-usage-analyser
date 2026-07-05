# Scenario

**Feature**: gitTracked column from sibling package.json git index status

```
named hit node_modules -> parent dir package.json -> git ls-files --error-unmatch -> gitTracked on named / named_size / named_enriched
```

## Preconditions

- Real `git` is required; leaves skip when `git` is not on PATH.
- Fixtures use `git init` (not bare `.git/` mkdir) so `ls-files` is meaningful.
- `node_modules` is never tracked; only sibling `package.json` tracking matters.

## Context

- Sibling leaves are mutually exclusive on package.json presence and git index state.
- `nested-pnpm-no-sibling` uses a nested `.pnpm/pkg@ver/node_modules` hit whose parent has no `package.json`.

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if !gitAvailable(t) {
		return nil
	}
	req.Op = "named-scan"
	req.Name = "node_modules"
	return nil
}

func gitAvailable(t *testing.T) bool {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available on PATH")
		return false
	}
	return true
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, out)
	}
}

func gitInitRepo(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")
}

func gitRepoWithNodeModules(t *testing.T, home, rel string) string {
	t.Helper()
	app := filepath.Join(home, filepath.FromSlash(rel))
	gitInitRepo(t, app)
	writeFile(t, app, "node_modules/pkg/dep.txt", []byte("fixture"), 0644)
	return app
}

func assertNamedEnrichedGitTracked(t *testing.T, resp *Response, want bool) {
	t.Helper()
	if len(resp.NamedEnrichedJSON) == 0 {
		t.Fatal("expected at least one named_enriched hit")
	}
	got := resp.NamedEnrichedJSON[0].GitTracked
	if got != want {
		t.Fatalf("named_enriched gitTracked = %v, want %v", got, want)
	}
}
```