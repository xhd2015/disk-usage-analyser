package nmpipeline

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"disk-usage-analyser/nminventory"
)

type stubCalc struct{}

func (stubCalc) PnpmCacheShared(path string) int64 { return 1024 }
func (stubCalc) BunCacheShared(path string) int64  { return 2048 }

type stubRunner struct{}

func (stubRunner) Run(dir string, name string, args ...string) (int, string, error) {
	_ = os.MkdirAll(filepath.Join(dir, "node_modules"), 0755)
	return 0, "", nil
}

func TestRunReportDryRunTable(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "app")
	nodeModules := filepath.Join(project, "node_modules")
	if err := os.MkdirAll(nodeModules, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, "package.json"), []byte(`{"name":"app"}`), 0644); err != nil {
		t.Fatal(err)
	}
	initGitRepo(t, project)

	inv := map[string]any{
		"version": "1.0",
		"node_modules": []map[string]any{{
			"path": nodeModules, "total_size": "20MB", "package_manager": "bun",
		}},
	}
	path := filepath.Join(t.TempDir(), "inventory.json")
	data, _ := json.Marshal(inv)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code, err := RunReport(Input{RecordsFile: path}, nminventory.RunOptions{DryRun: true, Workers: 1}, stubRunner{}, stubCalc{}, &stdout, &stderr)
	if err != nil || code != 0 {
		t.Fatalf("code=%d err=%v stderr=%s", code, err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "shared_size_added") {
		t.Fatalf("stdout=%s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "TOTAL shared_size_added: 0B") {
		t.Fatalf("expected zero added in dry-run: %s", stdout.String())
	}
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	cmd = exec.Command("git", "add", "package.json")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add: %v\n%s", err, out)
	}
	cmd = exec.Command("git", "commit", "-m", "init", "--author", "test <test@example.com>")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit: %v\n%s", err, out)
	}
}