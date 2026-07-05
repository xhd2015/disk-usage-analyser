package server

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDetectBelongsToGit(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	root := t.TempDir()
	runGit(t, root, "init")
	runGit(t, root, "config", "user.email", "test@example.com")
	runGit(t, root, "config", "user.name", "test")

	pkgDir := filepath.Join(root, "app")
	if err := os.MkdirAll(filepath.Join(pkgDir, "node_modules"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "package.json"), []byte(`{"name":"app"}`), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "app/package.json")
	runGit(t, root, "commit", "-m", "init")

	nmPath := filepath.Join(pkgDir, "node_modules")
	if !DetectBelongsToGit(nmPath) {
		t.Fatalf("DetectBelongsToGit(%q) = false, want true", nmPath)
	}

	outside := t.TempDir()
	if DetectBelongsToGit(outside) {
		t.Fatalf("DetectBelongsToGit(%q) = true, want false", outside)
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, out)
	}
}