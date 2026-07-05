package nmpipeline

import (
	"os"
	"path/filepath"
	"strings"

	"disk-usage-analyser/nmcacheshared"
	"disk-usage-analyser/nminventory"
	"disk-usage-analyser/nmmigrate"
)

// Row holds before/after measurements and migration outcome for one node_modules path.
type Row struct {
	Entry nminventory.Entry

	BeforeTotal int64
	BeforePnpm  int64
	BeforeBun   int64

	AfterTotal int64
	AfterPnpm  int64
	AfterBun   int64

	Migrate nmmigrate.MigrateResult
}

func (r Row) BeforeShared() int64 {
	return r.BeforePnpm + r.BeforeBun
}

func (r Row) AfterShared() int64 {
	return r.AfterPnpm + r.AfterBun
}

func (r Row) SharedAdded() int64 {
	return r.AfterShared() - r.BeforeShared()
}

func measureEntry(path string, calc nmcacheshared.CacheCalculator) (total, pnpm, bun int64) {
	pnpm, bun = nmcacheshared.MeasureShared(path, calc)
	total, err := nminventory.MeasureNodeModulesSize(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, pnpm, bun
		}
		return 0, pnpm, bun
	}
	return total, pnpm, bun
}

func processEntry(entry nminventory.Entry, dryRun bool, runner nmmigrate.CommandRunner, calc nmcacheshared.CacheCalculator) Row {
	beforeTotal, beforePnpm, beforeBun := measureEntry(entry.Path, calc)
	migrateResult := nmmigrate.Migrate(entry, dryRun, runner)
	afterTotal, afterPnpm, afterBun := measureEntry(entry.Path, calc)
	if dryRun {
		afterTotal, afterPnpm, afterBun = beforeTotal, beforePnpm, beforeBun
	}
	return Row{
		Entry:       entry,
		BeforeTotal: beforeTotal,
		BeforePnpm:  beforePnpm,
		BeforeBun:   beforeBun,
		AfterTotal:  afterTotal,
		AfterPnpm:   afterPnpm,
		AfterBun:    afterBun,
		Migrate:     migrateResult,
	}
}

func displayPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	rel, err := filepath.Rel(home, path)
	if err != nil || strings.HasPrefix(rel, "..") {
		return path
	}
	return filepath.Join("~", rel)
}