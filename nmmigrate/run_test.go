package nmmigrate

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

type stubRunner struct {
	exitCode int
	output   string
	err      error
	calls    []string
}

func (s *stubRunner) Run(dir string, name string, args ...string) (int, string, error) {
	s.calls = append(s.calls, filepath.Join(dir, strings.Join(append([]string{name}, args...), " ")))
	return s.exitCode, s.output, s.err
}

func TestMigrateOneDryRun(t *testing.T) {
	entry := nminventory.Entry{
		Raw:  []byte(`{"path":"/proj/node_modules","total_size":"20MB","package_manager":"bun"}`),
		Path: "/proj/node_modules",
	}
	result := migrateOne(entry, true, &stubRunner{})
	if result["dry_run"] != true || result["success"] != true {
		t.Fatalf("result = %#v", result)
	}
	if result["node_modules_removed"] != false {
		t.Fatalf("should not remove in dry-run")
	}
}

func TestMigrateOneRemovesAndRunsCorepack(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "app")
	nodeModules := filepath.Join(project, "node_modules")
	if err := os.MkdirAll(filepath.Join(nodeModules, "pkg"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, "package.json"), []byte(`{"name":"app"}`), 0644); err != nil {
		t.Fatal(err)
	}

	runner := &stubRunner{}
	entry := nminventory.Entry{
		Raw:  mustJSON(t, map[string]any{"path": nodeModules, "total_size": "20MB"}),
		Path: nodeModules,
	}
	result := migrateOne(entry, false, runner)
	if result["success"] != true {
		t.Fatalf("result = %#v", result)
	}
	if _, err := os.Stat(nodeModules); !os.IsNotExist(err) {
		t.Fatalf("node_modules still exists")
	}
	if len(runner.calls) != 1 || !strings.Contains(runner.calls[0], "corepack use pnpm@latest") {
		t.Fatalf("calls = %#v", runner.calls)
	}
}

func TestRunInventoryDryRunEmitsJSONL(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "tracked")
	nodeModules := filepath.Join(project, "node_modules")
	if err := os.MkdirAll(nodeModules, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, "package.json"), []byte(`{"name":"tracked"}`), 0644); err != nil {
		t.Fatal(err)
	}
	initGitRepo(t, project)

	inv := map[string]any{
		"version": "1.0",
		"node_modules": []map[string]any{
			{
				"path":             nodeModules,
				"has_package_json": true,
				"package_manager":  "bun",
				"total_size":       "20MB",
			},
		},
	}
	path := writeJSON(t, inv)

	var stdout, stderr bytes.Buffer
	code, err := RunInventory(path, nminventory.RunOptions{DryRun: true, Workers: 1}, &stubRunner{}, &stdout, &stderr)
	if err != nil || code != 0 {
		t.Fatalf("code=%d err=%v", code, err)
	}
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("stdout = %q", stdout.String())
	}
	var row map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &row); err != nil {
		t.Fatal(err)
	}
	if row["dry_run"] != true || row["success"] != true {
		t.Fatalf("row = %#v", row)
	}
	if _, err := os.Stat(nodeModules); err != nil {
		t.Fatalf("node_modules should still exist in dry-run")
	}
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	runGit(t, dir, "init")
	runGit(t, dir, "add", "package.json")
	runGit(t, dir, "commit", "-m", "init", "--author", "test <test@example.com>")
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v in %s: %v\n%s", args, dir, err, out)
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func writeJSON(t *testing.T, v any) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "inventory.json")
	if err := os.WriteFile(path, mustJSON(t, v), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}