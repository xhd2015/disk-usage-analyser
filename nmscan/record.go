package nmscan

import (
	"disk-usage-analyser/nmcacheshared"
	"disk-usage-analyser/nminventory"
	"disk-usage-analyser/server"
)

// BuildRecord assembles one inventory record for absPath.
func BuildRecord(absPath string, totalBytes int64, calc nmcacheshared.CacheCalculator, quick bool) nminventory.Record {
	pm := server.DetectPackageManager(absPath)
	if pm == "" {
		pm = "unknown"
	}

	shared := int64(0)
	if !quick && calc != nil {
		pnpm, bun := nmcacheshared.MeasureShared(absPath, calc)
		shared = pnpm + bun
	}

	return nminventory.Record{
		Path:           absPath,
		HasPackageJSON: server.DetectHasPackageJSON(absPath),
		PackageManager: pm,
		SharedSize:     nminventory.FormatCompactHumanSize(shared),
		TotalSize:      nminventory.FormatCompactHumanSize(totalBytes),
		BelongsToGit:   server.DetectBelongsToGit(absPath),
	}
}