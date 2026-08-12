// Package lib holds shared helpers for script/github/release.
//
// Pre-build runs the React SPA into disk-usage-analyser-react/dist so
// //go:embed ships a usable UI in multi-platform release binaries.
package lib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xhd2015/kool/pkgs/release"
	"github.com/xhd2015/xgo/support/cmd"
)

// DefaultSpecs is the multi-platform matrix (darwin/linux × amd64/arm64).
var DefaultSpecs = release.DefaultSpecs

// BinaryName is the released CLI basename (matches install.sh download formula).
const BinaryName = "disk-usage-analyser"

// PackagePath is the go build package for the CLI (module root main).
const PackagePath = "./"

// BuildRelease runs pre-build (unless skipPreBuild) then multi-platform builds.
func BuildRelease(skipPreBuild bool) (*release.BuildReleaseResult, error) {
	var pre func() error
	if !skipPreBuild {
		pre = PreBuild
	}
	return release.BuildRelease(BinaryName, pre, DefaultSpecs, release.WithPackagePath(PackagePath))
}

// PreBuild stages the frontend into disk-usage-analyser-react/dist for //go:embed.
func PreBuild() error {
	root, err := findModuleRoot()
	if err != nil {
		return err
	}
	fmt.Println("==> prebuild: frontend (go run ./script/build)")
	if err := cmd.Debug().Dir(root).Run("go", "run", "./script/build"); err != nil {
		return fmt.Errorf("frontend build: %w\n  hint: need node + npm/pnpm in disk-usage-analyser-react/; or re-run with --skip-prebuild if dist/ is already fat", err)
	}
	dist := filepath.Join(root, "disk-usage-analyser-react", "dist")
	st, err := os.Stat(dist)
	if err != nil || !st.IsDir() {
		return fmt.Errorf("frontend build did not produce dist: %s", dist)
	}
	fmt.Printf("Staged frontend → %s\n", dist)
	return nil
}

func findModuleRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", wd)
		}
		dir = parent
	}
}

// PlannedBinaryNames returns artifact basenames for tag + DefaultSpecs
// (same formula as release.BuildRelease).
func PlannedBinaryNames(tag string) []string {
	var out []string
	for _, spec := range DefaultSpecs {
		out = append(out, fmt.Sprintf("%s-%s-%s-%s", BinaryName, tag, spec.OS, spec.Arch))
	}
	return out
}
