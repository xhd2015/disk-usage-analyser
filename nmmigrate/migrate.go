package nmmigrate

import (
	"fmt"
	"os"
	"path/filepath"

	"disk-usage-analyser/nminventory"
	"disk-usage-analyser/server"
)

// MigrateResult captures the outcome of a single node_modules migration.
type MigrateResult struct {
	Path               string
	ProjectRoot        string
	DryRun             bool
	NodeModulesRemoved bool
	CorepackExitCode   int
	CorepackOutput     string
	Success            bool
	Error              string
}

// FilterEligible keeps entries with a sibling package.json that is git-tracked.
func FilterEligible(entries []nminventory.Entry, logf func(string, ...any)) ([]nminventory.Entry, int) {
	return filterEligible(entries, logf)
}

// Migrate removes node_modules and runs corepack use pnpm@latest in the project root.
func Migrate(entry nminventory.Entry, dryRun bool, runner CommandRunner) MigrateResult {
	result := migrateOne(entry, dryRun, runner)
	return MigrateResult{
		Path:               entry.Path,
		ProjectRoot:        filepath.Dir(filepath.Clean(entry.Path)),
		DryRun:             dryRun,
		NodeModulesRemoved: result["node_modules_removed"] == true,
		CorepackExitCode:   intFromAny(result["corepack_exit_code"]),
		CorepackOutput:     stringFromAny(result["corepack_output"]),
		Success:            result["success"] == true,
		Error:              stringFromAny(result["error"]),
	}
}

func intFromAny(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	default:
		return 0
	}
}

func stringFromAny(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func filterEligible(entries []nminventory.Entry, logf func(string, ...any)) ([]nminventory.Entry, int) {
	eligible := make([]nminventory.Entry, 0, len(entries))
	ineligible := 0
	for _, entry := range entries {
		if !server.DetectHasPackageJSON(entry.Path) {
			ineligible++
			if logf != nil {
				logf("skip index=%d path=%s no sibling package.json", entry.Index, entry.Path)
			}
			continue
		}
		if !server.DetectGitTrackedPackageJSON(entry.Path) {
			ineligible++
			if logf != nil {
				logf("skip index=%d path=%s package.json not git-tracked", entry.Index, entry.Path)
			}
			continue
		}
		eligible = append(eligible, entry)
	}
	return eligible, ineligible
}

func migrateOne(entry nminventory.Entry, dryRun bool, runner CommandRunner) map[string]any {
	obj, err := nminventory.DecodeRaw(entry.Raw)
	if err != nil {
		return map[string]any{
			"path":  entry.Path,
			"error": err.Error(),
		}
	}

	nodeModulesPath := filepath.Clean(entry.Path)
	projectRoot := filepath.Dir(nodeModulesPath)

	obj["project_root"] = projectRoot
	obj["corepack_command"] = corepackCommand
	obj["dry_run"] = dryRun
	obj["node_modules_removed"] = false
	obj["corepack_exit_code"] = 0
	obj["corepack_output"] = ""
	obj["success"] = false
	obj["error"] = ""

	if dryRun {
		obj["success"] = true
		return obj
	}

	if err := os.RemoveAll(nodeModulesPath); err != nil {
		obj["error"] = fmt.Sprintf("remove node_modules: %v", err)
		return obj
	}
	obj["node_modules_removed"] = true

	exitCode, output, err := runner.Run(projectRoot, "corepack", "use", "pnpm@latest")
	obj["corepack_exit_code"] = exitCode
	obj["corepack_output"] = output
	if err != nil {
		obj["error"] = fmt.Sprintf("corepack: %v", err)
		return obj
	}
	if exitCode != 0 {
		obj["error"] = fmt.Sprintf("corepack exited %d", exitCode)
		return obj
	}

	obj["success"] = true
	return obj
}