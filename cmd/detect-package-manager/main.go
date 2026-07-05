// Command detect-package-manager traces package manager detection for a node_modules path.
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
	walkUp := flag.Bool("walk-up", false, "also list lockfiles walking up from projectRoot (informational)")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: detect-package-manager [-home DIR] [-walk-up] ~/.../node_modules\n")
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

	displayPath := flag.Arg(0)
	absPath := filepath.Clean(resolveTildePath(displayPath, homeDir))

	fmt.Println("=== detect-package-manager trace ===")
	fmt.Printf("HOME=%s\n", homeDir)
	fmt.Printf("displayPath=%s\n", displayPath)
	fmt.Printf("absPath=%s\n", absPath)

	if info, err := os.Stat(absPath); err != nil {
		fmt.Printf("stat absPath: %v\n", err)
	} else if !info.IsDir() {
		fmt.Println("warn: absPath is not a directory")
	} else if filepath.Base(absPath) != "node_modules" {
		fmt.Printf("warn: basename is %q, expected node_modules\n", filepath.Base(absPath))
	}

	trace := server.TracePackageManager(absPath)
	fmt.Println("\n--- detection steps ---")
	for i, step := range trace.Steps {
		fmt.Printf("%2d. %s\n", i+1, step)
	}

	fmt.Println("\n--- result ---")
	fmt.Printf("projectRoot=%s\n", trace.ProjectRoot)
	fmt.Printf("packageManager=%s\n", trace.PackageManager)
	fmt.Printf("hasPackageJson=%v\n", trace.HasPackageJSON)

	hit := server.NamedHit{Path: displayPath}
	applyPackageManager(&hit, absPath)
	applyHasPackageJSON(&hit, absPath)

	fmt.Println("\n--- server NamedHit (pass1 / OnNamedHit) ---")
	printJSON(hit)

	namedSize := map[string]interface{}{
		"path":           hit.Path,
		"size":           hit.Size,
		"sizeHuman":      hit.SizeHuman,
		"hasPackageJson": hit.HasPackageJSON,
	}
	fmt.Println("\n--- SSE named_size (nested sizing rows — no packageManager field) ---")
	printJSON(namedSize)

	enriched := server.NamedEnrichedTraceFields(hit.Path, hit.PackageManager, absPath, hit.HasPackageJSON)
	fmt.Println("\n--- SSE named_enriched PM/pkgjson fields (shared metrics omitted in trace) ---")
	printJSON(enriched)

	fmt.Println("\n--- UI PM column (TmpFilesAnalyse.tsx) ---")
	// named_size handler builds a row without packageManager (see TmpFilesAnalyse.tsx).
	rowAfterNamedSize := map[string]string{
		"path": hit.Path,
		"packageManager": "",
	}
	fmt.Printf("after named_size only: hit.packageManager=%q -> UI %q\n",
		rowAfterNamedSize["packageManager"], uiLabel(rowAfterNamedSize["packageManager"]))
	mergedPM := enriched.PackageManager
	if mergedPM == "" {
		mergedPM = rowAfterNamedSize["packageManager"]
	}
	fmt.Printf("after named_enriched merge: enriched.packageManager=%q (server computed %q) -> UI %q\n",
		enriched.PackageManager, hit.PackageManager, uiLabel(mergedPM))
	if rowAfterNamedSize["packageManager"] == "" && enriched.PackageManager == "" {
		fmt.Println("note: UI shows unknown until named_enriched arrives with packageManager")
	}

	if *walkUp {
		fmt.Println("\n--- walk-up lockfiles (not used by detector today) ---")
		walkUpLockfiles(trace.ProjectRoot)
	}
}

func resolveTildePath(path string, homeDir string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

func applyPackageManager(hit *server.NamedHit, absPath string) {
	hit.PackageManager = server.DetectPackageManager(absPath)
}

func applyHasPackageJSON(hit *server.NamedHit, absPath string) {
	hit.HasPackageJSON = server.DetectHasPackageJSON(absPath)
}

func uiLabel(pm string) string {
	if pm == "" {
		return "unknown"
	}
	return pm
}

func printJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "json: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

func walkUpLockfiles(start string) {
	dir := filepath.Clean(start)
	for i := 0; i < 10 && dir != "/" && dir != "."; i++ {
		markers := []string{"bun.lockb", "bun.lock", "pnpm-lock.yaml", "package-lock.json", "yarn.lock", "package.json"}
		var found []string
		for _, name := range markers {
			if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
				found = append(found, name)
			}
		}
		fmt.Printf("  %s -> %v\n", dir, found)
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
}