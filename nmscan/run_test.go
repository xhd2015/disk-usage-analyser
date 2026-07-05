package nmscan

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"disk-usage-analyser/nmcacheshared"
	"disk-usage-analyser/nminventory"
)

type zeroCalculator struct{}

func (zeroCalculator) PnpmCacheShared(string) int64 { return 0 }
func (zeroCalculator) BunCacheShared(string) int64  { return 0 }

func TestBuildRecord(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	root := t.TempDir()
	runGit(t, root, "init")
	runGit(t, root, "config", "user.email", "test@example.com")
	runGit(t, root, "config", "user.name", "test")

	project := filepath.Join(root, "proj")
	nm := filepath.Join(project, "node_modules")
	if err := os.MkdirAll(nm, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, "package.json"), []byte(`{"name":"proj","packageManager":"npm@10.0.0"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, "package-lock.json"), []byte(`{"lockfileVersion":3}`), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "proj")
	runGit(t, root, "commit", "-m", "init")

	record := BuildRecord(nm, 20*1024*1024, zeroCalculator{}, true)
	if record.Path != nm {
		t.Fatalf("path = %q", record.Path)
	}
	if !record.HasPackageJSON {
		t.Fatal("has_package_json = false")
	}
	if record.PackageManager != "npm" {
		t.Fatalf("package_manager = %q", record.PackageManager)
	}
	if record.SharedSize != "0B" {
		t.Fatalf("shared_size = %q", record.SharedSize)
	}
	if record.TotalSize != "20MB" {
		t.Fatalf("total_size = %q", record.TotalSize)
	}
	if !record.BelongsToGit {
		t.Fatal("belongs_to_git = false")
	}
}

func TestWriteOutputFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")
	records := []nminventory.Record{
		{
			Path:           "/tmp/a/node_modules",
			HasPackageJSON: true,
			PackageManager: "npm",
			SharedSize:     "0B",
			TotalSize:      "10MB",
			BelongsToGit:   false,
		},
	}
	if err := writeOutputFile(path, records); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var out nminventory.OutputFile
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	if out.Version != "1.0" || len(out.NodeModules) != 1 {
		t.Fatalf("out = %#v", out)
	}
	if out.NodeModules[0].BelongsToGit {
		t.Fatal("belongs_to_git should be false")
	}
}

func TestRunCLIHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code, err := RunCLI([]string{"--help"}, CLIOptions{Stdout: &stdout, Stderr: &stderr})
	if err != nil || code != 0 {
		t.Fatalf("code=%d err=%v", code, err)
	}
	if !strings.Contains(stdout.String(), "node-modules-scan") {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestRunOnTempRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	root := t.TempDir()
	home := filepath.Join(root, "home")
	if err := os.MkdirAll(home, 0755); err != nil {
		t.Fatal(err)
	}

	repo := filepath.Join(home, "Projects", "demo")
	if err := os.MkdirAll(filepath.Join(repo, "node_modules"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "package.json"), []byte(`{"name":"demo"}`), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "init")
	runGit(t, repo, "config", "user.email", "test@example.com")
	runGit(t, repo, "config", "user.name", "test")
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-m", "init")

	var stdout, stderr bytes.Buffer
	outPath := filepath.Join(root, "inventory.json")
	code, err := RunCLI([]string{
		"--quick",
		"--workers", "1",
		"--verbose",
		"-o", outPath,
	}, CLIOptions{
		Stdout:     &stdout,
		Stderr:     &stderr,
		HomeDir:    home,
		Calculator: zeroCalculator{},
	})
	if err != nil || code != 0 {
		t.Fatalf("code=%d err=%v stderr=%q", code, err, stderr.String())
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) == 0 || lines[0] == "" {
		t.Fatalf("stdout empty, stderr=%q", stderr.String())
	}
	var row map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &row); err != nil {
		t.Fatal(err)
	}
	if row["belongs_to_git"] != true {
		t.Fatalf("row = %#v", row)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	var file nminventory.OutputFile
	if err := json.Unmarshal(data, &file); err != nil {
		t.Fatal(err)
	}
	if file.Version != "1.0" || len(file.NodeModules) != 1 {
		t.Fatalf("file = %#v", file)
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

var _ nmcacheshared.CacheCalculator = zeroCalculator{}