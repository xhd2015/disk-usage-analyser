package server

import (
	"os/exec"
	"strings"
)

// DetectBelongsToGit reports whether nodeModulesAbsPath lies inside a git work tree.
func DetectBelongsToGit(nodeModulesAbsPath string) bool {
	cmd := exec.Command("git", "-C", nodeModulesAbsPath, "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}