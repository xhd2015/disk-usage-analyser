# Scenario

**Feature**: tmp-files scan shared CLI fixture harness

```
tmp-files scan -> roots -> scan_repo repo discovery -> repo walk -> detect/buildinfo -> BinaryHit stream -> summary
```

## Preconditions

- Tests must never scan the real user home directory.
- Each leaf creates a temporary fake home and places fixture repositories under it.
- A git repository is represented by a directory containing `.git/`.
- Mach-O and ELF fixtures use minimal magic bytes sufficient for file type detection.
- Go binary fixtures are built from a tiny `main` package so `debug/buildinfo.ReadFile` succeeds.

## Steps

1. Create a fresh fixture home for each leaf.
2. Set `req.Args` to the `tmp-files` subcommand arguments below `tmp-files`.
3. Run the CLI with injected stdout, stderr, and home directory.

## Context

- The root `Run` calls `tmpfiles.RunCLI` for CLI leaves.
- The dispatch leaf calls `run.RunWithOptions` to prove `tmp-files` bypasses web serving.
- Fixture helpers intentionally live in the doctest tree so implementation does not depend on testdata binaries.

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
	req.Args = []string{"scan"}
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

func mkdir(t *testing.T, base string, rel string) string {
	t.Helper()
	dir := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir %s: %v", rel, err)
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

func cloudStorageProvider(t *testing.T, home string, provider string) string {
	t.Helper()
	return mkdir(t, home, "Library/CloudStorage/"+provider)
}

func writeGoBinary(t *testing.T, base string, rel string) string {
	t.Helper()
	srcDir := mkdir(t, base, ".doctest-go-src")
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
