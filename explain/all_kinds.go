package explain

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"disk-usage-analyser/usagescan"
)

// allKindsCLIRegistry is the locked v1 --all-kinds pack order (cli kind ids).
var allKindsCLIRegistry = []string{"xcode", "grok", "android-sdk", "iterm2", "codex"}

// buildAllKindsResult measures every registered pack under absolute scope.
// Missing pack roots become status "missing" with totalSize 0; the overall
// result always succeeds (caller exits 0 even when all packs are missing).
func buildAllKindsResult(scope string) (AllKindsResult, error) {
	scope = strings.TrimSpace(scope)
	if scope == "" {
		return AllKindsResult{}, fmt.Errorf("scope is required for --all-kinds")
	}
	abs, err := filepath.Abs(scope)
	if err != nil {
		return AllKindsResult{}, fmt.Errorf("invalid scope %s: %w", scope, err)
	}

	entries := make([]KindEntry, 0, len(allKindsCLIRegistry))
	var totalPresent int64

	for _, cliKind := range allKindsCLIRegistry {
		entry := measureOneAllKindsEntry(cliKind, abs)
		if entry.Status == "present" {
			totalPresent += entry.TotalSize
		}
		entries = append(entries, entry)
	}

	// Present size DESC, then missing/error (stable by original registry order within status).
	sort.SliceStable(entries, func(i, j int) bool {
		pi := entries[i].Status == "present"
		pj := entries[j].Status == "present"
		if pi != pj {
			return pi // present before missing/error
		}
		if pi && pj && entries[i].TotalSize != entries[j].TotalSize {
			return entries[i].TotalSize > entries[j].TotalSize
		}
		return false // keep registry order
	})

	return AllKindsResult{
		Scope:     abs,
		TotalSize: totalPresent,
		Kinds:     entries,
	}, nil
}

// measureOneAllKindsEntry builds one KindEntry for a cli pack under scope.
func measureOneAllKindsEntry(cliKind, scope string) KindEntry {
	cliKind = strings.ToLower(strings.TrimSpace(cliKind))
	expectedPath := expectedPackPath(cliKind, scope)

	switch cliKind {
	case "xcode":
		// Multi-root pack: present only if at least one root exists under scope.
		exp, err := buildKindPackExplanation("xcode", scope)
		if err != nil {
			return KindEntry{
				Kind:      "xcode",
				CLIKind:   "xcode",
				Path:      expectedPath,
				Status:    "error",
				TotalSize: 0,
				Error:     err.Error(),
			}
		}
		if len(exp.Breakdown) == 0 && exp.TotalSize == 0 {
			return KindEntry{
				Kind:      "xcode",
				CLIKind:   "xcode",
				Path:      expectedPath,
				Status:    "missing",
				TotalSize: 0,
			}
		}
		return kindEntryFromExplanation(cliKind, exp)

	case "grok":
		exp, err := buildGrokKindExplanation(scope)
		if err != nil {
			return missingOrErrorEntry(cliKind, "grok-home", expectedPath, err)
		}
		return kindEntryFromExplanation(cliKind, exp)

	case "android-sdk":
		exp, err := buildAndroidSDKKindExplanation(scope)
		if err != nil {
			return missingOrErrorEntry(cliKind, "android-sdk", expectedPath, err)
		}
		return kindEntryFromExplanation(cliKind, exp)

	case "iterm2":
		exp, err := buildITerm2KindExplanation(scope)
		if err != nil {
			return missingOrErrorEntry(cliKind, "iterm2", expectedPath, err)
		}
		return kindEntryFromExplanation(cliKind, exp)

	case "codex":
		exp, err := buildCodexKindExplanation(scope)
		if err != nil {
			return missingOrErrorEntry(cliKind, "codex-home", expectedPath, err)
		}
		return kindEntryFromExplanation(cliKind, exp)

	default:
		return KindEntry{
			Kind:      cliKind,
			CLIKind:   cliKind,
			Path:      expectedPath,
			Status:    "error",
			TotalSize: 0,
			Error:     fmt.Sprintf("unknown kind pack %q", cliKind),
		}
	}
}

// expectedPackPath returns the default measured path for a pack under scope
// (used when the pack is missing or for INDEX PATH column).
func expectedPackPath(cliKind, scope string) string {
	switch strings.ToLower(cliKind) {
	case "xcode":
		return scope
	case "grok":
		if filepath.Base(scope) == ".grok" {
			return scope
		}
		return filepath.Join(scope, ".grok")
	case "android-sdk":
		if isAndroidSDKRoot(scope) {
			return scope
		}
		return filepath.Join(scope, filepath.FromSlash("Library/Android/sdk"))
	case "iterm2":
		if isITerm2Root(scope) {
			return scope
		}
		return filepath.Join(scope, filepath.FromSlash("Library/Application Support/iTerm2"))
	case "codex":
		if filepath.Base(scope) == ".codex" {
			return scope
		}
		return filepath.Join(scope, ".codex")
	default:
		return scope
	}
}

func missingOrErrorEntry(cliKind, outputKind, path string, err error) KindEntry {
	msg := err.Error()
	// Path-does-not-exist style errors → missing (overall run still exits 0).
	if os.IsNotExist(err) || strings.Contains(msg, "path does not exist") ||
		strings.Contains(msg, "no such file") {
		return KindEntry{
			Kind:      outputKind,
			CLIKind:   cliKind,
			Path:      path,
			Status:    "missing",
			TotalSize: 0,
		}
	}
	return KindEntry{
		Kind:      outputKind,
		CLIKind:   cliKind,
		Path:      path,
		Status:    "error",
		TotalSize: 0,
		Error:     msg,
	}
}

func kindEntryFromExplanation(cliKind string, exp Explanation) KindEntry {
	return KindEntry{
		Kind:       exp.Kind,
		CLIKind:    cliKind,
		Path:       exp.Path,
		Status:     "present",
		TotalSize:  exp.TotalSize,
		Confidence: exp.Confidence,
		Summary:    exp.Summary,
		Breakdown:  exp.Breakdown,
		Reclaim:    exp.Reclaim,
		HowToPurge: exp.HowToPurge,
	}
}

func writeAllKindsJSON(w io.Writer, res AllKindsResult) error {
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	// One JSON line, then trailing blank line (ends with \n\n).
	_, err = fmt.Fprintf(w, "%s\n\n", data)
	return err
}

func writeAllKindsHuman(w io.Writer, res AllKindsResult, useColor bool) error {
	var b strings.Builder

	// Header
	fmt.Fprintf(&b, "SCOPE: %s\n", res.Scope)
	fmt.Fprintf(&b, "MODE: all-kinds\n")
	fmt.Fprintf(&b, "TOTAL (present): %s\n", usagescan.FormatCompactHumanSize(res.TotalSize))

	presentN, missingN := 0, 0
	for _, e := range res.Kinds {
		switch e.Status {
		case "present":
			presentN++
		case "missing":
			missingN++
		}
	}
	fmt.Fprintf(&b, "PRESENT: %d  MISSING: %d\n", presentN, missingN)
	fmt.Fprintln(&b)

	// INDEX table
	fmt.Fprintln(&b, "INDEX")
	writeAllKindsIndexTable(&b, res.Kinds)
	fmt.Fprintln(&b)

	// Detail sections for present packs only
	for _, e := range res.Kinds {
		if e.Status != "present" {
			continue
		}
		writeAllKindsDetail(&b, e, useColor)
		fmt.Fprintln(&b)
	}

	// Optional MISSING list
	var missing []KindEntry
	for _, e := range res.Kinds {
		if e.Status == "missing" {
			missing = append(missing, e)
		}
	}
	if len(missing) > 0 {
		fmt.Fprintln(&b, "MISSING")
		for _, e := range missing {
			fmt.Fprintf(&b, "  - %s (%s): %s\n", e.CLIKind, e.Kind, e.Path)
		}
	}

	out := b.String()
	if !strings.HasSuffix(out, "\n") {
		out += "\n"
	}
	// Trailing blank line after last content.
	out += "\n"
	_, err := io.WriteString(w, out)
	return err
}

func writeAllKindsIndexTable(b *strings.Builder, kinds []KindEntry) {
	// Columns: SIZE, KIND, STATUS, PATH
	type row struct {
		size, kind, status, path string
	}
	rows := make([]row, len(kinds))
	maxSize := len("SIZE")
	maxKind := len("KIND")
	maxStatus := len("STATUS")
	maxPath := len("PATH")

	for i, e := range kinds {
		size := usagescan.FormatCompactHumanSize(e.TotalSize)
		// Prefer cliKind in INDEX so all 4 registry ids appear as tests expect.
		kind := e.CLIKind
		if kind == "" {
			kind = e.Kind
		}
		r := row{
			size:   size,
			kind:   kind,
			status: e.Status,
			path:   e.Path,
		}
		rows[i] = r
		if len(size) > maxSize {
			maxSize = len(size)
		}
		if len(kind) > maxKind {
			maxKind = len(kind)
		}
		if len(e.Status) > maxStatus {
			maxStatus = len(e.Status)
		}
		if len(e.Path) > maxPath {
			maxPath = len(e.Path)
		}
	}

	fmt.Fprintf(b, "  %-*s  %-*s  %-*s  %s\n",
		maxSize, "SIZE", maxKind, "KIND", maxStatus, "STATUS", "PATH")
	for _, r := range rows {
		fmt.Fprintf(b, "  %*s  %-*s  %-*s  %s\n",
			maxSize, r.size, maxKind, r.kind, maxStatus, r.status, r.path)
	}
}

// writeAllKindsDetail writes a mini-explain for one present KindEntry.
func writeAllKindsDetail(b *strings.Builder, e KindEntry, useColor bool) {
	fmt.Fprintf(b, "KIND: %s\n", e.Kind)
	fmt.Fprintf(b, "PATH: %s\n", e.Path)
	fmt.Fprintf(b, "TOTAL: %s\n", usagescan.FormatCompactHumanSize(e.TotalSize))
	if e.Confidence != "" {
		fmt.Fprintf(b, "CONFIDENCE: %s\n", e.Confidence)
	}
	fmt.Fprintln(b)

	fmt.Fprintln(b, "SUMMARY")
	if len(e.Summary) == 0 {
		fmt.Fprintln(b, "  (none)")
	} else {
		for _, line := range e.Summary {
			fmt.Fprintf(b, "  %s\n", line)
		}
	}
	fmt.Fprintln(b)

	fmt.Fprintln(b, "BREAKDOWN")
	writeBreakdownTable(b, e.Breakdown, useColor)

	// Short SAFE TO RECLAIM / HOW TO PURGE (reuse single-kind content when present).
	if len(e.Reclaim) > 0 {
		fmt.Fprintln(b)
		fmt.Fprintln(b, "SAFE TO RECLAIM")
		for _, r := range e.Reclaim {
			safe := "caution"
			if r.SafeToReclaim {
				safe = "usually safe"
			}
			fmt.Fprintf(b, "  - %s (%s): %s\n", r.Title, safe, r.Detail)
		}
	}
	if len(e.HowToPurge) > 0 {
		fmt.Fprintln(b)
		fmt.Fprintln(b, "HOW TO PURGE")
		for i, p := range e.HowToPurge {
			if i > 0 {
				fmt.Fprintln(b)
			}
			fmt.Fprintf(b, "  %d) %s\n", i+1, p.Title)
			fmt.Fprintln(b, "     Official command:")
			for _, line := range strings.Split(p.OfficialCommand, "\n") {
				fmt.Fprintf(b, "       %s\n", formatHumanCommandLine(line, useColor))
			}
			fmt.Fprintf(b, "     Removes: %s\n", p.Removes)
			if p.Notes != "" {
				fmt.Fprintf(b, "     Notes: %s\n", p.Notes)
			}
		}
	}
}
