package server

import (
	"os/exec"
	"path/filepath"
)

// DetectGitTrackedPackageJSON reports whether package.json beside node_modules
// (same parent directory) is tracked in git. node_modules is typically ignored;
// package.json is not, so ls-files is a reliable signal.
func DetectGitTrackedPackageJSON(nodeModulesAbsPath string) bool {
	projectRoot := filepath.Dir(filepath.Clean(nodeModulesAbsPath))
	pkgPath := filepath.Join(projectRoot, "package.json")
	if !fileExists(pkgPath) {
		return false
	}
	cmd := exec.Command("git", "-C", projectRoot, "ls-files", "--error-unmatch", "--", "package.json")
	return cmd.Run() == nil
}