package explain

import "io"

// CLIOptions configures explain CLI I/O.
type CLIOptions struct {
	Stdout io.Writer
	Stderr io.Writer
	// HomeDir is optional; when non-empty, used as default scope for --kind without PATH.
	// When empty, RunCLI falls back to os.UserHomeDir().
	HomeDir string
}

// Explanation is the in-memory / JSON model for explain output.
type Explanation struct {
	Path        string       `json:"path"`
	Kind        string       `json:"kind"`
	TotalSize   int64        `json:"totalSize"`
	Confidence  string       `json:"confidence"`
	Summary     []string     `json:"summary"`
	Breakdown   []Breakdown  `json:"breakdown"`
	Reclaim     []Reclaim    `json:"reclaim"`
	HowToPurge  []PurgeStep  `json:"howToPurge"`
	RawCommands []RawCommand `json:"rawCommands"`
	// LogsDB is set for codex-home when a readable logs_*.sqlite is present (shape A).
	LogsDB *LogsDBInfo `json:"logsDb,omitempty"`
}

// LogsDBInfo is the Codex logs_*.sqlite preview (row count + last samples).
type LogsDBInfo struct {
	Path    string      `json:"path"`
	Size    int64       `json:"size"`
	Rows    int64       `json:"rows"`
	Samples []LogSample `json:"samples"`
}

// LogSample is one row from the Codex logs table (newest-first sample window).
type LogSample struct {
	ID     int64  `json:"id"`
	TS     int64  `json:"ts"` // unix seconds
	Level  string `json:"level"`
	Target string `json:"target"`
	Body   string `json:"body"` // truncated
}

// Breakdown is one semantic size entry.
type Breakdown struct {
	Name        string `json:"name"`
	Path        string `json:"path,omitempty"`
	Size        int64  `json:"size"`
	Role        string `json:"role,omitempty"`
	Notes       string `json:"notes,omitempty"`
	Reclaimable bool   `json:"reclaimable"` // true when role is reclaimable tier
}

// Reclaim is safe-to-reclaim advice (never suggests rm -rf).
type Reclaim struct {
	Title         string `json:"title"`
	SafeToReclaim bool   `json:"safeToReclaim"`
	Detail        string `json:"detail"`
}

// PurgeStep is an official/recommended purge recipe for this kind.
// Commands are informational; never use rm -rf.
type PurgeStep struct {
	// Title short label for the purge option.
	Title string `json:"title"`
	// OfficialCommand is the preferred tool command (may be multi-line text).
	OfficialCommand string `json:"officialCommand"`
	// Removes describes which files/data this purge removes.
	Removes string `json:"removes"`
	// Notes optional extra caution (emulator stopped, etc.).
	Notes string `json:"notes,omitempty"`
}

// RawCommand is an informational command for further inspection.
type RawCommand struct {
	Group   string `json:"group"`
	Command string `json:"command"`
}

// AllKindsResult is the multi-pack report for --all-kinds (human + JSON envelope).
type AllKindsResult struct {
	Scope     string      `json:"scope"`
	TotalSize int64       `json:"totalSize"` // sum of present kinds only
	Kinds     []KindEntry `json:"kinds"`
}

// KindEntry is one pack result inside AllKindsResult.
type KindEntry struct {
	Kind       string      `json:"kind"`    // output kind id (grok-home, iterm2, …)
	CLIKind    string      `json:"cliKind"` // xcode, grok, android-sdk, iterm2
	Path       string      `json:"path"`
	Status     string      `json:"status"` // present | missing | error
	TotalSize  int64       `json:"totalSize"`
	Confidence string      `json:"confidence,omitempty"`
	Summary    []string    `json:"summary,omitempty"`
	Breakdown  []Breakdown `json:"breakdown,omitempty"`
	Reclaim    []Reclaim   `json:"reclaim,omitempty"`
	HowToPurge []PurgeStep `json:"howToPurge,omitempty"`
	Error      string      `json:"error,omitempty"`
}

// detectResult is the outcome of kind detection before measurement.
type detectResult struct {
	Kind       string
	Confidence string
	// ContextRoot is the path to measure for kind-specific breakdown
	// (e.g. AVD dir when the user pointed at a file inside it).
	ContextRoot string
	// TargetIsFile is true when the original path is a regular file.
	TargetIsFile bool
}
