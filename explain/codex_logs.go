package explain

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"disk-usage-analyser/usagescan"

	_ "modernc.org/sqlite"
)

const logsBodySampleMaxRunes = 160

// findCodexLogsDBPath returns the first logs_*.sqlite file (not -wal/-shm) under codexDir.
func findCodexLogsDBPath(codexDir string) string {
	entries, err := os.ReadDir(codexDir)
	if err != nil {
		return ""
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, "-wal") || strings.HasSuffix(name, "-shm") {
			continue
		}
		if strings.HasPrefix(name, "logs_") && strings.HasSuffix(name, ".sqlite") {
			return filepath.Join(codexDir, name)
		}
	}
	return ""
}

// readCodexLogsDB opens logs_*.sqlite read-only and returns row count + last 3 samples.
// Returns (nil, nil) when no logs db is present; returns error only for unexpected failures
// that the caller may ignore (explain still succeeds).
func readCodexLogsDB(codexDir string) (*LogsDBInfo, error) {
	dbPath := findCodexLogsDBPath(codexDir)
	if dbPath == "" {
		return nil, nil
	}
	st, err := os.Stat(dbPath)
	if err != nil {
		return nil, err
	}
	if st.IsDir() || st.Size() == 0 {
		return nil, nil
	}

	// Pure Go SQLite driver; open read-only so we never mutate the user's logs DB.
	dsn := "file:" + filepath.ToSlash(dbPath) + "?mode=ro&_pragma=query_only(1)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var rows int64
	if err := db.QueryRow(`SELECT COUNT(*) FROM logs`).Scan(&rows); err != nil {
		return nil, err
	}

	samples, err := queryCodexLogSamples(db, 3)
	if err != nil {
		return nil, err
	}

	return &LogsDBInfo{
		Path:    dbPath,
		Size:    st.Size(),
		Rows:    rows,
		Samples: samples,
	}, nil
}

func queryCodexLogSamples(db *sql.DB, limit int) ([]LogSample, error) {
	if limit <= 0 {
		limit = 3
	}
	q := fmt.Sprintf(`
SELECT id, ts, level, target, COALESCE(feedback_log_body, '')
FROM logs
ORDER BY ts DESC, ts_nanos DESC, id DESC
LIMIT %d`, limit)
	rs, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var out []LogSample
	for rs.Next() {
		var s LogSample
		var body string
		if err := rs.Scan(&s.ID, &s.TS, &s.Level, &s.Target, &body); err != nil {
			return nil, err
		}
		s.Body = truncateRunes(body, logsBodySampleMaxRunes)
		out = append(out, s)
	}
	return out, rs.Err()
}

func truncateRunes(s string, max int) string {
	if max <= 0 || s == "" {
		return s
	}
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	if len(runes) > max {
		return string(runes[:max]) + "…"
	}
	return s
}

// writeLogsDBSection formats the human LOGS DB section (product shape A).
func writeLogsDBSection(b *strings.Builder, info *LogsDBInfo) {
	if info == nil {
		return
	}
	fmt.Fprintln(b, "LOGS DB")
	fmt.Fprintf(b, "  PATH: %s\n", info.Path)
	fmt.Fprintf(b, "  SIZE: %s\n", formatLogsDBSize(info.Size))
	fmt.Fprintf(b, "  ROWS: %d\n", info.Rows)
	fmt.Fprintln(b, "  SAMPLE (last 3, newest first):")
	if len(info.Samples) == 0 {
		fmt.Fprintln(b, "    (none)")
		return
	}
	for i, s := range info.Samples {
		ts := formatLogSampleTime(s.TS)
		fmt.Fprintf(b, "    %d) %s %s %s\n", i+1, ts, s.Level, s.Target)
		if s.Body != "" {
			fmt.Fprintf(b, "       %s\n", s.Body)
		}
	}
}

func formatLogsDBSize(n int64) string {
	return usagescan.FormatCompactHumanSize(n)
}

func formatLogSampleTime(ts int64) string {
	if ts <= 0 {
		return "0"
	}
	// Unix seconds → UTC RFC3339-ish short form; fixtures use small ints (1000…).
	return time.Unix(ts, 0).UTC().Format(time.RFC3339)
}
