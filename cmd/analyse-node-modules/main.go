// Command analyse-node-modules prints node_modules table columns for a single path.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"disk-usage-analyser/server"
)

func main() {
	homeFlag := flag.String("home", "", "HOME for ~/ path resolution (default: os.UserHomeDir())")
	jsonOut := flag.Bool("json", false, "emit full NamedHit JSON")
	quick := flag.Bool("quick", false, "skip shared metrics (faster; Shared shows 0 B)")
	verbose := flag.Bool("v", false, "log per-step timing to stderr")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: analyse-node-modules [-home DIR] [-json] [-quick] [-v] PATH\n")
		fmt.Fprintf(os.Stderr, "  PATH may be absolute or ~/.../node_modules\n")
		os.Exit(2)
	}

	homeDir := *homeFlag
	if homeDir == "" {
		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "home: %v\n", err)
			os.Exit(1)
		}
	}

	inputPath := flag.Arg(0)
	absPath := filepath.Clean(resolveTildePath(inputPath, homeDir))

	if filepath.Base(absPath) != "node_modules" {
		fmt.Fprintf(os.Stderr, "warn: basename is %q, expected node_modules\n", filepath.Base(absPath))
	}

	hit, err := server.AnalyseNodeModules(inputPath, homeDir, server.AnalyseNodeModulesOptions{
		SkipShared: *quick,
		Verbose:    *verbose,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "analyse: %v\n", err)
		os.Exit(1)
	}

	if *jsonOut {
		printJSON(hit)
		return
	}

	fmt.Println("=== node_modules (UI columns) ===")
	fmt.Printf("Path:           %s\n", hit.Path)
	fmt.Printf("package.json:   %s\n", boolLabel(hit.HasPackageJSON))
	fmt.Printf("Git:            %s\n", boolLabel(hit.GitTracked))
	fmt.Printf("PM:             %s\n", uiPackageManager(hit.PackageManager))
	if *quick {
		fmt.Printf("Shared:         %s (skipped; use without -quick)\n", hit.SharedHuman)
	} else {
		fmt.Printf("Shared:         %s\n", hit.SharedHuman)
	}
	fmt.Printf("Size:           %s\n", hit.SizeHuman)

	if !*quick && (hit.PnpmSharedSize > 0 || hit.BunSharedSize > 0) {
		fmt.Println()
		fmt.Println("--- shared breakdown ---")
		fmt.Printf("pnpm_shared:    %s (%d B)\n", hit.PnpmSharedHuman, hit.PnpmSharedSize)
		fmt.Printf("bun_shared:     %s (%d B)\n", hit.BunSharedHuman, hit.BunSharedSize)
	}
}

func resolveTildePath(path string, homeDir string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

func boolLabel(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

func uiPackageManager(pm string) string {
	if pm == "" {
		return "unknown"
	}
	return pm
}

func printJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "json: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}