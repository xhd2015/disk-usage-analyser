package nminventory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// CleanNodeModulesPath expands ~/ and returns an absolute node_modules path.
func CleanNodeModulesPath(path string) (string, error) {
	p := strings.TrimSpace(path)
	if strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		p = filepath.Join(home, p[2:])
	}
	return filepath.Abs(p)
}

// EntryFromPath builds an inventory entry from a node_modules directory path.
func EntryFromPath(path string, index int) (Entry, bool) {
	abs, err := CleanNodeModulesPath(path)
	if err != nil {
		return Entry{}, false
	}
	info, err := os.Stat(abs)
	if err != nil || !info.IsDir() {
		return Entry{}, false
	}

	total, err := MeasureNodeModulesSize(abs)
	if err != nil {
		return Entry{}, false
	}
	totalHuman := FormatCompactHumanSize(total)
	raw, err := json.Marshal(map[string]any{
		"path":       abs,
		"total_size": totalHuman,
	})
	if err != nil {
		return Entry{}, false
	}
	return Entry{
		Index:      index,
		Raw:        raw,
		Path:       abs,
		TotalBytes: total,
		TotalSize:  totalHuman,
	}, true
}