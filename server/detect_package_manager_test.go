package server

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectPackageManagerMixedLockfilesPreferPnpm(t *testing.T) {
	dir, err := os.MkdirTemp("", "detect-pm-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	write := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	write("package.json", `{"name":"demo","packageManager":"npm@10.0.0"}`)
	write("pnpm-lock.yaml", "lockfileVersion: '9.0'\n")
	write("bun.lock", "{}")
	write("package-lock.json", `{"lockfileVersion":3}`)
	if err := os.MkdirAll(filepath.Join(dir, "node_modules"), 0755); err != nil {
		t.Fatal(err)
	}

	got := DetectPackageManager(filepath.Join(dir, "node_modules"))
	if got != "pnpm" {
		t.Fatalf("DetectPackageManager() = %q, want pnpm", got)
	}
}