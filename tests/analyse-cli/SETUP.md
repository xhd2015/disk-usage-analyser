# Scenario

**Feature**: analyse CLI shared fixture harness

```
analyse [DIR] -> walk tree (no symlink follow) -> per-immediate-child DirResult rows -> summary -> TSV or JSON
```

## Preconditions

- Tests use isolated temporary fixture trees under `t.TempDir()`.
- Symlinks are created with `os.Symlink`; hard links with `os.Link`.
- Fixture sizes use exact byte counts so numeric assertions are deterministic.
- The `analyse` package provides `Analyse`, `RunCLI`, and `CLIOptions`.

## Steps

1. Create a fresh fixture root for each leaf.
2. Set `req.FixtureDir` to the fixture path (or `req.Args` for CLI leaves).
3. Run `analyse.Analyse` or `analyse.RunCLI` via the root `Run` harness.

## Context

- Default `req.Mode` is `analyse` (structured `Result` assertions).
- CLI leaves set `req.Mode` to `cli` or `dispatch`.
- Human-readable size strings follow du -h style; leaves assert numeric `int64` fields primarily.

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"disk-usage-analyser/analyse"
)

const file4K = 4096

func Setup(t *testing.T, req *Request) error {
	req.FixtureDir = filepath.Join(t.TempDir(), "fixture")
	if err := os.MkdirAll(req.FixtureDir, 0755); err != nil {
		return err
	}
	if req.Mode == "" {
		req.Mode = "analyse"
	}
	return nil
}

func mkdir(t *testing.T, base string, rel string) string {
	t.Helper()
	dir := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir %s: %v", rel, err)
	}
	return dir
}

func writeSizedFile(t *testing.T, base string, rel string, size int64) string {
	t.Helper()
	path := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir parent %s: %v", rel, err)
	}
	data := make([]byte, size)
	for i := range data {
		data[i] = 'x'
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
	return path
}

func symlinkTo(t *testing.T, base string, linkRel string, targetRel string) string {
	t.Helper()
	link := filepath.Join(base, filepath.FromSlash(linkRel))
	target := filepath.Join(base, filepath.FromSlash(targetRel))
	if err := os.MkdirAll(filepath.Dir(link), 0755); err != nil {
		t.Fatalf("mkdir symlink parent %s: %v", linkRel, err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("symlink %s -> %s: %v", linkRel, targetRel, err)
	}
	return link
}

func hardlinkFile(t *testing.T, base string, newRel string, existingRel string) string {
	t.Helper()
	newPath := filepath.Join(base, filepath.FromSlash(newRel))
	existing := filepath.Join(base, filepath.FromSlash(existingRel))
	if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
		t.Fatalf("mkdir hardlink parent %s: %v", newRel, err)
	}
	if err := os.Link(existing, newPath); err != nil {
		t.Fatalf("hardlink %s -> %s: %v", newRel, existingRel, err)
	}
	return newPath
}

func symlinkLabel(files, dirs int) string {
	return fmt.Sprintf("%df+%dd", files, dirs)
}

func rowByPath(result *analyse.Result, want string) *analyse.DirResult {
	for i := range result.Rows {
		if result.Rows[i].Path == want {
			return &result.Rows[i]
		}
	}
	return nil
}
```