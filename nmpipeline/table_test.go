package nmpipeline

import (
	"bytes"
	"strings"
	"testing"

	"disk-usage-analyser/nminventory"
	"disk-usage-analyser/nmmigrate"
)

func TestFormatTable(t *testing.T) {
	rows := []Row{
		{
			Entry:       nminventory.Entry{Path: "/tmp/big/node_modules"},
			BeforeTotal: 200 * 1024 * 1024,
			BeforePnpm:  0,
			BeforeBun:   100 * 1024 * 1024,
			AfterTotal:  180 * 1024 * 1024,
			AfterPnpm:   150 * 1024 * 1024,
			AfterBun:    0,
			Migrate:     nmmigrate.MigrateResult{Success: true},
		},
	}
	var buf bytes.Buffer
	if err := FormatTable(&buf, rows); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{
		"path",
		"before_size",
		"before_pnpm",
		"before_bun",
		"after_size",
		"after_pnpm",
		"after_bun",
		"shared_size_added",
		"TOTAL shared_size_added:",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
}