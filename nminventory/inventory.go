package nminventory

import (
	"encoding/json"
	"fmt"
	"os"
)

type File struct {
	Version     string            `json:"version"`
	NodeModules []json.RawMessage `json:"node_modules"`
}

type Entry struct {
	Index      int
	Raw        json.RawMessage
	Path       string
	TotalBytes int64
	TotalSize  string
}

type RunOptions struct {
	Workers       int
	Limit         int
	SizeThreshold int64
	DryRun        bool
	Verbose       bool
}

func Load(path string) (File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return File{}, err
	}
	var input File
	if err := json.Unmarshal(data, &input); err != nil {
		return File{}, err
	}
	return input, nil
}

func ParseEntries(rawEntries []json.RawMessage, logf func(string, ...any)) ([]Entry, int) {
	var entries []Entry
	skipped := 0
	for i, raw := range rawEntries {
		entry, ok := ParseEntry(i, raw)
		if !ok {
			skipped++
			if logf != nil {
				logf("skip index=%d missing or invalid path/total_size", i)
			}
			continue
		}
		entries = append(entries, entry)
	}
	return entries, skipped
}

func ParseEntry(index int, raw json.RawMessage) (Entry, bool) {
	var fields struct {
		Path      string `json:"path"`
		TotalSize string `json:"total_size"`
	}
	if err := json.Unmarshal(raw, &fields); err != nil || fields.Path == "" {
		return Entry{}, false
	}
	totalBytes, err := ParseCompactHumanSize(fields.TotalSize)
	if err != nil {
		return Entry{}, false
	}
	return Entry{
		Index:      index,
		Raw:        raw,
		Path:       fields.Path,
		TotalBytes: totalBytes,
		TotalSize:  fields.TotalSize,
	}, true
}

func FilterBySizeThreshold(entries []Entry, threshold int64, logf func(string, ...any)) ([]Entry, int) {
	if threshold <= 0 {
		return entries, 0
	}
	filtered := make([]Entry, 0, len(entries))
	below := 0
	for _, entry := range entries {
		if entry.TotalBytes < threshold {
			below++
			if logf != nil {
				logf("skip index=%d total_size=%s below threshold %s",
					entry.Index, entry.TotalSize, FormatCompactHumanSize(threshold))
			}
			continue
		}
		filtered = append(filtered, entry)
	}
	return filtered, below
}

func ApplyLimit(entries []Entry, limit int) []Entry {
	if limit <= 0 || limit >= len(entries) {
		return entries
	}
	return entries[:limit]
}

func DecodeRaw(raw json.RawMessage) (map[string]any, error) {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, fmt.Errorf("decode entry: %w", err)
	}
	return obj, nil
}