// usage: go run ./script/github/release [--dry-run] [--skip-prebuild]
//
// Release multi-platform disk-usage-analyser binaries to GitHub Releases.
// Pattern matches github.com/xhd2015/browser-agent/script/github/release
// (and doctest).
//
// Live requirements:
//   - clean git worktree
//   - HEAD is a v* tag (git describe --tags HEAD)
//   - .upload-credentials.json {token, owner, repo}
//   - node/npm for frontend prebuild (unless --skip-prebuild with existing dist/)
package main

import (
	"fmt"
	"os"
	"strings"

	githublib "disk-usage-analyser/script/github/lib"

	"github.com/xhd2015/kool/pkgs/github"
	"github.com/xhd2015/kool/pkgs/release"
	"github.com/xhd2015/less-flags"
)

const help = `
Usage: go run ./script/github/release [options]

Release disk-usage-analyser multi-platform binaries to GitHub Releases.

Asset names:
  disk-usage-analyser-{tag}-{os}-{arch}
  e.g. disk-usage-analyser-v0.1.0-darwin-arm64

Prebuild:
  Runs go run ./script/build so //go:embed disk-usage-analyser-react/dist
  ships a usable UI. Use --skip-prebuild only when dist/ is already fat.

Options:
  --dry-run         print plan without building or uploading
  --skip-prebuild   skip frontend build (require existing dist/)
  -h, --help        show this help

Credentials (live only):
  .upload-credentials.json
    {"token":"ghp_...","owner":"xhd2015","repo":"disk-usage-analyser"}

Install published binaries:
  curl -fsSL https://raw.githubusercontent.com/xhd2015/disk-usage-analyser/master/install.sh | bash
`

func main() {
	if err := handle(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handle() error {
	var dryRun bool
	var skipPrebuild bool
	args, err := lessflags.
		Bool("--dry-run", &dryRun).
		Bool("--skip-prebuild", &skipPrebuild).
		Help("-h,--help", help).
		Parse(os.Args[1:])
	if err != nil {
		return err
	}
	if len(args) > 0 {
		return fmt.Errorf("unrecognized extra args: %s", strings.Join(args, " "))
	}

	tag, err := release.GetTag()
	if err != nil {
		if !dryRun {
			return err
		}
		fmt.Fprintf(os.Stderr, "[dry-run] warning: %v\n", err)
		tag = "(unknown)"
	}

	creds, err := release.LoadCredentials(".upload-credentials.json")
	if err != nil {
		if !dryRun {
			return err
		}
		fmt.Fprintf(os.Stderr, "[dry-run] warning: %v\n", err)
		creds = &release.Credentials{Owner: "xhd2015", Repo: "disk-usage-analyser"}
	}

	plannedBins := githublib.PlannedBinaryNames(tag)

	if dryRun {
		fmt.Printf("[dry-run] tag: %s\n", tag)
		if skipPrebuild {
			fmt.Println("[dry-run] would skip frontend prebuild (--skip-prebuild)")
		} else {
			fmt.Println("[dry-run] would prebuild: go run ./script/build")
		}
		for _, name := range plannedBins {
			fmt.Printf("[dry-run] would build: %s\n", name)
		}
		fmt.Printf("[dry-run] would upload to %s/%s release (creates if 404)\n", creds.Owner, creds.Repo)
		return nil
	}

	result, err := githublib.BuildRelease(skipPrebuild)
	if err != nil {
		return err
	}

	client := github.NewReleaseClient(creds.Token, creds.Owner, creds.Repo)
	rel, err := client.GetOrCreateRelease(result.Tag)
	if err != nil {
		return fmt.Errorf("failed to get or create release for tag %s: %v", result.Tag, err)
	}

	for _, file := range result.Files {
		if err := client.UploadReleaseAsset(rel.ID, file); err != nil {
			return fmt.Errorf("failed to upload %s: %v", file, err)
		}
		fmt.Printf("Uploaded %s\n", file)
	}

	fmt.Printf("Release %s: https://github.com/%s/%s/releases/tag/%s\n",
		result.Tag, creds.Owner, creds.Repo, result.Tag)
	return nil
}
