package explain

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"disk-usage-analyser/usagescan"

	lessflags "github.com/xhd2015/less-flags"
)

const help = `Usage: disk-usage-analyser explain [PATH] [--kind KIND] [--all-kinds] [--json] [--color=...]

Explain reclaim kind, size breakdown, and safe-to-reclaim advice.

  PATH     Directory or file. Required unless --kind or --all-kinds is set
           (then PATH is optional scope; default: home directory).

Options:
  --kind KIND            Force pack/kind (xcode, grok, android-sdk, iterm2, codex).
                         PATH optional scope (default: home). For grok, measures
                         {scope}/.grok (or scope if it is already .grok). For
                         android-sdk, measures {scope}/Library/Android/sdk (or
                         scope if it is already an Android SDK root). For iterm2,
                         measures {scope}/Library/Application Support/iTerm2 (or
                         scope if it is already an iTerm2 Application Support root).
                         For codex, measures {scope}/.codex (or scope if it is
                         already .codex). Mutually exclusive with --all-kinds.
  --all-kinds            Analyse all registered packs (xcode, grok, android-sdk,
                         iterm2, codex) under optional scope (default: home). Prints
                         an INDEX of all kinds plus detail for each present pack.
                         Mutually exclusive with --kind.
  --json                 Emit one JSON object instead of human sections
                         (Explanation for single PATH/--kind; AllKindsResult for
                         --all-kinds)
  --color=always|never|auto
                         Color base command on human command lines (default: auto)
                         Bare --color means always. auto colors only when stdout is a TTY
                         and NO_COLOR is unset; --color=always forces on.
  -h, --help             Show help
`

// RunCLI implements: disk-usage-analyser explain [PATH] [--kind KIND] [--all-kinds] [--json] [--color=...] [-h]
// args are the tokens after the "explain" subcommand.
func RunCLI(args []string, opts CLIOptions) (int, error) {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	args = normalizeColorArgs(args)

	var jsonOut bool
	var allKinds bool
	var kindFlag string
	colorMode := "auto"
	remain, err := lessflags.
		Bool("--json", &jsonOut).
		Bool("--all-kinds", &allKinds).
		String("--kind", &kindFlag).
		String("--color", &colorMode).
		HelpFunc("-h,--help", func() {
			fmt.Fprint(stdout, help)
			if !strings.HasSuffix(help, "\n") {
				fmt.Fprintln(stdout)
			}
		}).
		HelpNoExit().
		Parse(args)
	if err == lessflags.ErrHelp {
		return 0, nil
	}
	if err != nil {
		return 2, err
	}

	switch strings.ToLower(strings.TrimSpace(colorMode)) {
	case "", "auto", "always", "never":
		// ok
	default:
		return 2, fmt.Errorf("invalid --color value %q (want always|never|auto)", colorMode)
	}
	if colorMode == "" {
		colorMode = "auto"
	}

	kindFlag = strings.TrimSpace(kindFlag)
	if kindFlag != "" {
		if !isKnownKindPack(kindFlag) {
			return 2, fmt.Errorf("unknown kind %q (supported: xcode, grok, android-sdk, iterm2, codex)", kindFlag)
		}
	}

	// --all-kinds and --kind are mutually exclusive.
	if allKinds && kindFlag != "" {
		return 2, fmt.Errorf("--all-kinds and --kind are mutually exclusive; use one or the other")
	}

	if len(remain) == 0 && kindFlag == "" && !allKinds {
		return 2, fmt.Errorf("PATH is required unless --kind or --all-kinds is set; usage: disk-usage-analyser explain [PATH] [--kind KIND] [--all-kinds]")
	}
	if len(remain) > 1 {
		return 2, fmt.Errorf("unexpected extra argument: %s", remain[1])
	}

	// Multi-pack mode: --all-kinds
	if allKinds {
		scope, err := resolveKindScope(remain, opts.HomeDir)
		if err != nil {
			return 1, err
		}
		res, err := buildAllKindsResult(scope)
		if err != nil {
			return 1, err
		}
		if jsonOut {
			if err := writeAllKindsJSON(stdout, res); err != nil {
				return 1, err
			}
			return 0, nil
		}
		useColor := shouldColor(colorMode, stdout)
		if err := writeAllKindsHuman(stdout, res, useColor); err != nil {
			return 1, err
		}
		return 0, nil
	}

	var exp Explanation

	if kindFlag != "" {
		// Forced pack/kind: PATH is optional scope (default HomeDir / user home).
		scope, err := resolveKindScope(remain, opts.HomeDir)
		if err != nil {
			return 1, err
		}
		switch strings.ToLower(kindFlag) {
		case "grok":
			// Alias: measure {scope}/.grok (or scope if already .grok); kind id grok-home.
			exp, err = buildGrokKindExplanation(scope)
		case "codex":
			// Alias: measure {scope}/.codex (or scope if already .codex); kind id codex-home.
			exp, err = buildCodexKindExplanation(scope)
		case "android-sdk":
			// Measure {scope}/Library/Android/sdk (or scope if already an SDK root).
			exp, err = buildAndroidSDKKindExplanation(scope)
		case "iterm2":
			// Measure {scope}/Library/Application Support/iTerm2 (or scope if already iTerm2 root).
			exp, err = buildITerm2KindExplanation(scope)
		default:
			exp, err = buildKindPackExplanation(kindFlag, scope)
		}
		if err != nil {
			return 1, err
		}
	} else {
		target := remain[0]
		absPath, err := filepath.Abs(target)
		if err != nil {
			return 1, fmt.Errorf("invalid path %s: %w", target, err)
		}

		info, err := os.Stat(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				return 1, fmt.Errorf("path does not exist: %s", absPath)
			}
			return 1, fmt.Errorf("cannot access %s: %w", absPath, err)
		}

		det := detectKind(absPath, info)
		exp, err = buildExplanation(absPath, info, det)
		if err != nil {
			return 1, err
		}
	}

	if jsonOut {
		if err := writeJSON(stdout, exp); err != nil {
			return 1, err
		}
		return 0, nil
	}

	useColor := shouldColor(colorMode, stdout)
	if err := writeHuman(stdout, exp, useColor); err != nil {
		return 1, err
	}
	return 0, nil
}

// resolveKindScope returns the absolute scope directory for a --kind pack.
// remain may be empty (use HomeDir / user home) or one PATH argument.
func resolveKindScope(remain []string, homeDir string) (string, error) {
	if len(remain) == 0 {
		scope := strings.TrimSpace(homeDir)
		if scope == "" {
			var err error
			scope, err = os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("cannot resolve home directory for --kind scope: %w", err)
			}
		}
		abs, err := filepath.Abs(scope)
		if err != nil {
			return "", fmt.Errorf("invalid home/scope %s: %w", scope, err)
		}
		return abs, nil
	}

	abs, err := filepath.Abs(remain[0])
	if err != nil {
		return "", fmt.Errorf("invalid path %s: %w", remain[0], err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("path does not exist: %s", abs)
		}
		return "", fmt.Errorf("cannot access %s: %w", abs, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("kind pack scope must be a directory: %s", abs)
	}
	return abs, nil
}

// normalizeColorArgs turns bare --color into --color=always so less-flags
// string parsing does not consume the next positional argument.
func normalizeColorArgs(args []string) []string {
	out := make([]string, 0, len(args))
	for _, a := range args {
		if a == "--color" {
			out = append(out, "--color=always")
			continue
		}
		out = append(out, a)
	}
	return out
}

// shouldColor decides whether human command lines get green base-command ANSI.
// --color=always forces on (even with NO_COLOR); never forces off;
// auto colors only when stdout is a TTY and NO_COLOR is unset.
func shouldColor(mode string, stdout io.Writer) bool {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "always":
		return true
	case "never":
		return false
	default: // auto
		if os.Getenv("NO_COLOR") != "" {
			return false
		}
		return isTTY(stdout)
	}
}

func isTTY(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func writeJSON(w io.Writer, exp Explanation) error {
	data, err := json.Marshal(exp)
	if err != nil {
		return err
	}
	// One JSON line, then trailing blank line (ends with \n\n).
	_, err = fmt.Fprintf(w, "%s\n\n", data)
	return err
}

func writeHuman(w io.Writer, exp Explanation, useColor bool) error {
	var b strings.Builder

	fmt.Fprintf(&b, "PATH: %s\n", exp.Path)
	fmt.Fprintf(&b, "KIND: %s\n", exp.Kind)
	fmt.Fprintf(&b, "TOTAL: %s\n", usagescan.FormatCompactHumanSize(exp.TotalSize))
	fmt.Fprintf(&b, "CONFIDENCE: %s\n", exp.Confidence)
	fmt.Fprintln(&b)

	fmt.Fprintln(&b, "SUMMARY")
	for _, line := range exp.Summary {
		fmt.Fprintf(&b, "  %s\n", line)
	}
	if len(exp.Summary) == 0 {
		fmt.Fprintln(&b, "  (none)")
	}
	fmt.Fprintln(&b)

	fmt.Fprintln(&b, "BREAKDOWN")
	writeBreakdownTable(&b, exp.Breakdown, useColor)
	fmt.Fprintln(&b)

	// Product shape A: LOGS DB after BREAKDOWN, before SAFE TO RECLAIM (codex-home).
	if exp.LogsDB != nil {
		writeLogsDBSection(&b, exp.LogsDB)
		fmt.Fprintln(&b)
	}

	fmt.Fprintln(&b, "SAFE TO RECLAIM")
	for _, r := range exp.Reclaim {
		safe := "caution"
		if r.SafeToReclaim {
			safe = "usually safe"
		}
		fmt.Fprintf(&b, "  - %s (%s): %s\n", r.Title, safe, r.Detail)
	}
	if len(exp.Reclaim) == 0 {
		fmt.Fprintln(&b, "  - No automated reclaim advice for this path.")
	}
	fmt.Fprintln(&b)

	fmt.Fprintln(&b, "HOW TO PURGE")
	if len(exp.HowToPurge) == 0 {
		fmt.Fprintln(&b, "  (no official purge recipe for this kind)")
	} else {
		for i, p := range exp.HowToPurge {
			if i > 0 {
				fmt.Fprintln(&b)
			}
			fmt.Fprintf(&b, "  %d) %s\n", i+1, p.Title)
			fmt.Fprintln(&b, "     Official command:")
			for _, line := range strings.Split(p.OfficialCommand, "\n") {
				fmt.Fprintf(&b, "       %s\n", formatHumanCommandLine(line, useColor))
			}
			fmt.Fprintf(&b, "     Removes: %s\n", p.Removes)
			if p.Notes != "" {
				fmt.Fprintf(&b, "     Notes: %s\n", p.Notes)
			}
		}
	}
	fmt.Fprintln(&b)

	fmt.Fprintln(&b, "RAW COMMANDS")
	var lastGroup string
	for _, c := range exp.RawCommands {
		if c.Group != lastGroup {
			fmt.Fprintf(&b, "  # %s\n", c.Group)
			lastGroup = c.Group
		}
		fmt.Fprintf(&b, "  %s\n", formatHumanCommandLine(c.Command, useColor))
	}
	if len(exp.RawCommands) == 0 {
		fmt.Fprintln(&b, "  (none)")
	}

	// Trailing blank line after last content line.
	out := b.String()
	if !strings.HasSuffix(out, "\n") {
		out += "\n"
	}
	out += "\n"
	_, err := io.WriteString(w, out)
	return err
}

// writeBreakdownTable formats BREAKDOWN as an aligned multi-column table:
// SIZE (right), NAME, ROLE (optional green/yellow), RECLAIMABLE (☑/☐; ☑ green when color on),
// optional NOTES. Empty breakdown prints "  (empty)". Indent is two spaces under the section header.
func writeBreakdownTable(b *strings.Builder, items []Breakdown, useColor bool) {
	if len(items) == 0 {
		fmt.Fprintln(b, "  (empty)")
		return
	}

	type row struct {
		size, name, role, reclaim, notes string
		reclaimable                      bool
	}
	rows := make([]row, len(items))
	hasNotes := false
	maxSize := len("SIZE")
	maxName := len("NAME")
	maxRole := len("ROLE")
	// Header "RECLAIMABLE" is wider than a single-cell glyph; pad to header width.
	maxReclaim := len("RECLAIMABLE")

	for i, item := range items {
		size := usagescan.FormatCompactHumanSize(item.Size)
		// Unicode option E: ☑ reclaimable (U+2611), ☐ not reclaimable (U+2610).
		reclaim := "☐"
		if item.Reclaimable {
			reclaim = "☑"
		}
		rows[i] = row{
			size:        size,
			name:        item.Name,
			role:        item.Role,
			reclaim:     reclaim,
			notes:       item.Notes,
			reclaimable: item.Reclaimable,
		}
		if len(size) > maxSize {
			maxSize = len(size)
		}
		if len(item.Name) > maxName {
			maxName = len(item.Name)
		}
		if len(item.Role) > maxRole {
			maxRole = len(item.Role)
		}
		if item.Notes != "" {
			hasNotes = true
		}
	}

	// Header: SIZE left-padded to same width as right-aligned values (end column shared).
	if hasNotes {
		fmt.Fprintf(b, "  %-*s  %-*s  %-*s  %-*s  %s\n",
			maxSize, "SIZE", maxName, "NAME", maxRole, "ROLE", maxReclaim, "RECLAIMABLE", "NOTES")
	} else {
		fmt.Fprintf(b, "  %-*s  %-*s  %-*s  %-*s\n",
			maxSize, "SIZE", maxName, "NAME", maxRole, "ROLE", maxReclaim, "RECLAIMABLE")
	}

	for _, r := range rows {
		// ROLE / RECLAIMABLE cells may include ANSI; pad on visible width so alignment holds.
		roleCell := colorRoleCell(r.role, useColor)
		if pad := maxRole - len(r.role); pad > 0 {
			roleCell += strings.Repeat(" ", pad)
		}
		reclaimCell := colorReclaimCheckbox(r.reclaim, r.reclaimable, useColor)
		// Visible width of ☑/☐ is one display cell (one rune each).
		if pad := maxReclaim - 1; pad > 0 {
			reclaimCell += strings.Repeat(" ", pad)
		}
		if hasNotes {
			fmt.Fprintf(b, "  %*s  %-*s  %s  %s  %s\n",
				maxSize, r.size, maxName, r.name, roleCell, reclaimCell, r.notes)
		} else {
			fmt.Fprintf(b, "  %*s  %-*s  %s  %s\n",
				maxSize, r.size, maxName, r.name, roleCell, reclaimCell)
		}
	}
}

// colorRoleCell wraps ROLE text in green (reclaimable) or yellow (caution) when color is on.
func colorRoleCell(role string, useColor bool) string {
	if !useColor || role == "" {
		return role
	}
	const (
		green  = "\x1b[32m"
		yellow = "\x1b[33m"
		reset  = "\x1b[0m"
	)
	switch roleTierOf(role) {
	case tierReclaimable:
		return green + role + reset
	case tierCaution:
		return yellow + role + reset
	default:
		return role
	}
}

// colorReclaimCheckbox greens reclaimable ☑ when color is on; ☐ is never colored.
func colorReclaimCheckbox(glyph string, reclaimable, useColor bool) string {
	if !useColor || !reclaimable {
		return glyph
	}
	const (
		green = "\x1b[32m"
		reset = "\x1b[0m"
	)
	return green + glyph + reset
}

// formatHumanCommandLine adds "$ " to runnable lines and optionally greens the base command.
// Comment lines (after trim, starting with #) are left without "$" and without color.
func formatHumanCommandLine(line string, useColor bool) string {
	trim := strings.TrimSpace(line)
	if trim == "" {
		return line
	}
	if strings.HasPrefix(trim, "#") {
		return trim
	}
	// Runnable command: human shell-prompt prefix; JSON path never calls this.
	if useColor {
		return "$ " + colorBaseCommand(trim)
	}
	return "$ " + trim
}

// colorBaseCommand wraps the first whitespace-separated token in green SGR.
// The "$" prompt is applied by the caller and is never colored.
func colorBaseCommand(cmd string) string {
	base, rest, ok := strings.Cut(cmd, " ")
	const green = "\x1b[32m"
	const reset = "\x1b[0m"
	if !ok {
		return green + base + reset
	}
	return green + base + reset + " " + rest
}
