package server

import (
	"github.com/xhd2015/dot-pkgs/go-pkgs/npm"
)

// PackageManagerTrace records how DetectPackageManager resolved a node_modules path.
type PackageManagerTrace struct {
	NodeModulesAbsPath string
	ProjectRoot        string
	PackageManager     string
	HasPackageJSON     bool
	Steps              []string
}

// TracePackageManager runs DetectPackageManager with step-by-step diagnostics.
func TracePackageManager(nodeModulesAbsPath string) PackageManagerTrace {
	trace := npm.DetectFromNodeModules(nodeModulesAbsPath)
	return PackageManagerTrace{
		NodeModulesAbsPath: trace.NodeModulesAbsPath,
		ProjectRoot:        trace.ProjectRoot,
		PackageManager:     string(trace.Manager),
		HasPackageJSON:     trace.HasPackageJSON,
		Steps:              trace.Steps,
	}
}

// DetectPackageManager inspects the project root (parent of node_modules) for
// lockfile markers and returns npm, pnpm, yarn, bun, or unknown.
func DetectPackageManager(nodeModulesAbsPath string) string {
	return npm.ManagerFromNodeModules(nodeModulesAbsPath)
}

// DetectHasPackageJSON reports whether package.json exists in the project root
// (parent of node_modules).
func DetectHasPackageJSON(nodeModulesAbsPath string) bool {
	return npm.HasPackageJSON(nodeModulesAbsPath)
}