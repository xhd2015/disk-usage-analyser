package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PackageManagerTrace records how DetectPackageManager resolved a node_modules path.
type PackageManagerTrace struct {
	NodeModulesAbsPath string
	ProjectRoot        string
	PackageManager     string
	HasPackageJSON     bool
	Steps              []string
}

func traceStep(steps *[]string, format string, args ...interface{}) {
	*steps = append(*steps, fmt.Sprintf(format, args...))
}

// TracePackageManager runs DetectPackageManager with step-by-step diagnostics.
func TracePackageManager(nodeModulesAbsPath string) PackageManagerTrace {
	clean := filepath.Clean(nodeModulesAbsPath)
	projectRoot := filepath.Dir(clean)
	trace := PackageManagerTrace{
		NodeModulesAbsPath: clean,
		ProjectRoot:        projectRoot,
		Steps:              make([]string, 0, 12),
	}

	checkFile := func(label, path string) bool {
		ok := fileExists(path)
		traceStep(&trace.Steps, "check %s: %s -> %v", label, path, ok)
		return ok
	}
	checkDir := func(label, path string) bool {
		ok := dirExists(path)
		traceStep(&trace.Steps, "check %s: %s -> %v", label, path, ok)
		return ok
	}

	traceStep(&trace.Steps, "projectRoot = parent(node_modules) = %s", projectRoot)

	if checkFile("bun.lockb", filepath.Join(projectRoot, "bun.lockb")) ||
		checkFile("bun.lock", filepath.Join(projectRoot, "bun.lock")) {
		trace.PackageManager = "bun"
		trace.HasPackageJSON = fileExists(filepath.Join(projectRoot, "package.json"))
		return trace
	}
	if checkFile("pnpm-lock.yaml", filepath.Join(projectRoot, "pnpm-lock.yaml")) {
		trace.PackageManager = "pnpm"
		trace.HasPackageJSON = fileExists(filepath.Join(projectRoot, "package.json"))
		return trace
	}
	if checkFile("package-lock.json", filepath.Join(projectRoot, "package-lock.json")) {
		trace.PackageManager = "npm"
		trace.HasPackageJSON = fileExists(filepath.Join(projectRoot, "package.json"))
		return trace
	}
	if checkFile("yarn.lock", filepath.Join(projectRoot, "yarn.lock")) {
		trace.PackageManager = "yarn"
		trace.HasPackageJSON = fileExists(filepath.Join(projectRoot, "package.json"))
		return trace
	}
	if checkDir("node_modules/.pnpm", filepath.Join(clean, ".pnpm")) {
		trace.PackageManager = "pnpm"
		trace.HasPackageJSON = fileExists(filepath.Join(projectRoot, "package.json"))
		return trace
	}

	pkgPath := filepath.Join(projectRoot, "package.json")
	if checkFile("package.json", pkgPath) {
		trace.HasPackageJSON = true
		field := parsePackageManagerField(pkgPath)
		traceStep(&trace.Steps, "package.json packageManager field -> %q", field)
		if field != "" {
			trace.PackageManager = field
			return trace
		}
		traceStep(&trace.Steps, "package.json present, no packageManager field -> default npm")
		trace.PackageManager = "npm"
		return trace
	}

	traceStep(&trace.Steps, "no lockfiles, .pnpm, or package.json in projectRoot -> unknown")
	trace.PackageManager = "unknown"
	trace.HasPackageJSON = false
	return trace
}

// DetectPackageManager inspects the project root (parent of node_modules) for
// lockfile markers and returns npm, pnpm, yarn, bun, or unknown.
func DetectPackageManager(nodeModulesAbsPath string) string {
	return TracePackageManager(nodeModulesAbsPath).PackageManager
}

type packageJSONPackageManager struct {
	PackageManager string `json:"packageManager"`
}

func parsePackageManagerField(pkgPath string) string {
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return ""
	}
	var pkg packageJSONPackageManager
	if err := json.Unmarshal(data, &pkg); err != nil {
		return ""
	}
	tool, _, ok := strings.Cut(pkg.PackageManager, "@")
	if !ok || tool == "" {
		return ""
	}
	switch tool {
	case "npm", "pnpm", "yarn", "bun":
		return tool
	default:
		return ""
	}
}

// DetectHasPackageJSON reports whether package.json exists in the project root
// (parent of node_modules).
func DetectHasPackageJSON(nodeModulesAbsPath string) bool {
	projectRoot := filepath.Dir(nodeModulesAbsPath)
	return fileExists(filepath.Join(projectRoot, "package.json"))
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}