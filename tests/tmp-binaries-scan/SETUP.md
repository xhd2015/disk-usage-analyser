# Scenario

**Feature**: tmp-binaries-scan shared fixture harness

```
GET /api/tmp-binaries-scan -> scan_repo discovery -> classify binaries -> SSE binary/summary/done
```

## Preconditions

- Tests must never scan the real user home directory.
- Each leaf creates a temporary fake home and sets `HOME` before httptest.
- Git repositories are directories containing `.git/`.
- Mach-O and ELF fixtures use minimal magic bytes; Go binaries are built from a tiny `main` package.

## Steps

1. Create a fresh fixture home for each leaf.
2. Place binaries inside git repos under the fake home.
3. Set `req.Op` and run the httptest SSE harness.

## Context

- Fixture helpers live in this doctest tree (same patterns as `tests/tmp-files-scan`).

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.HomeDir = filepath.Join(t.TempDir(), "home")
	if err := os.MkdirAll(req.HomeDir, 0755); err != nil {
		return err
	}
	if req.Op == "" {
		req.Op = "binaries-scan"
	}
	return nil
}

func repo(t *testing.T, home string, rel string) string {
	t.Helper()
	dir := filepath.Join(home, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0755); err != nil {
		t.Fatalf("create repo %s: %v", rel, err)
	}
	return dir
}

func writeFile(t *testing.T, base string, rel string, data []byte, mode os.FileMode) string {
	t.Helper()
	path := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir parent %s: %v", rel, err)
	}
	if err := os.WriteFile(path, data, mode); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
	return path
}

func writeMachO(t *testing.T, base string, rel string) string {
	t.Helper()
	return writeFile(t, base, rel, append([]byte{0xcf, 0xfa, 0xed, 0xfe, 0x0c, 0x00, 0x00, 0x01}, make([]byte, 96)...), 0755)
}

func writeELF(t *testing.T, base string, rel string) string {
	t.Helper()
	return writeFile(t, base, rel, append([]byte{0x7f, 'E', 'L', 'F', 2, 1, 1, 0}, make([]byte, 96)...), 0755)
}

func writeText(t *testing.T, base string, rel string, text string) string {
	t.Helper()
	return writeFile(t, base, rel, []byte(text), 0644)
}

func writeGoBinary(t *testing.T, base string, rel string) string {
	t.Helper()
	srcDir := filepath.Join(base, ".doctest-go-src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, srcDir, "main.go", []byte("package "+"main\n\nfunc "+"main() {}\n"), 0644)
	out := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
		t.Fatalf("mkdir go binary parent: %v", err)
	}
	cmd := exec.Command("go", "build", "-o", out, ".")
	cmd.Dir = srcDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build go binary: %v\n%s", err, string(output))
	}
	return out
}
```