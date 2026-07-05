package nmcacheshared

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"disk-usage-analyser/nminventory"
)

type stubCalculator struct {
	pnpm int64
	bun  int64
}

func (s stubCalculator) PnpmCacheShared(path string) int64 { return s.pnpm }
func (s stubCalculator) BunCacheShared(path string) int64  { return s.bun }

func TestRunInventoryDryRun(t *testing.T) {
	path := writeTempJSON(t, sampleInventory())
	var stdout, stderr bytes.Buffer
	code, err := RunInventory(path, nminventory.RunOptions{
		SizeThreshold: 10 * 1024 * 1024,
		Limit:         1,
		DryRun:        true,
	}, nil, &stdout, &stderr)
	if err != nil || code != 0 {
		t.Fatalf("RunInventory: code=%d err=%v", code, err)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "would_scan=1") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunInventoryEmitsJSONL(t *testing.T) {
	path := writeTempJSON(t, sampleInventory())
	var stdout, stderr bytes.Buffer
	code, err := RunInventory(path, nminventory.RunOptions{
		Workers:       1,
		SizeThreshold: 10 * 1024 * 1024,
		Limit:         1,
	}, stubCalculator{pnpm: 1024, bun: 2048}, &stdout, &stderr)
	if err != nil || code != 0 {
		t.Fatalf("RunInventory: code=%d err=%v", code, err)
	}
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("lines = %d", len(lines))
	}
	var row map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &row); err != nil {
		t.Fatal(err)
	}
	if row["pnpm_cache_shared"] != "1KB" || row["bun_cache_shared"] != "2KB" {
		t.Fatalf("row = %#v", row)
	}
}

func sampleInventory() []byte {
	inv := map[string]any{
		"version": "1.0",
		"node_modules": []map[string]any{
			{"path": "/tmp/small/node_modules", "total_size": "3.4MB"},
			{"path": "/tmp/large/node_modules", "total_size": "20MB"},
		},
	}
	data, _ := json.Marshal(inv)
	return data
}

func writeTempJSON(t *testing.T, data []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "inventory.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	return path
}