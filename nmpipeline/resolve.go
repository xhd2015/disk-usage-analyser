package nmpipeline

import (
	"fmt"

	"disk-usage-analyser/nminventory"
	"disk-usage-analyser/nmmigrate"
)

func resolveEntries(input Input, opts nminventory.RunOptions, logf func(string, ...any)) ([]nminventory.Entry, error) {
	byPath := make(map[string]nminventory.Entry)
	order := make([]string, 0)
	index := 0

	addEntry := func(entry nminventory.Entry) {
		if _, seen := byPath[entry.Path]; !seen {
			order = append(order, entry.Path)
		}
		byPath[entry.Path] = entry
	}

	if input.RecordsFile != "" {
		file, err := nminventory.Load(input.RecordsFile)
		if err != nil {
			return nil, fmt.Errorf("read records %q: %w", input.RecordsFile, err)
		}
		parsed, _ := nminventory.ParseEntries(file.NodeModules, logf)
		filtered, _ := nminventory.FilterBySizeThreshold(parsed, opts.SizeThreshold, logf)
		for _, entry := range filtered {
			addEntry(entry)
		}
		index = len(file.NodeModules)
	}

	for _, argPath := range input.Paths {
		entry, ok := nminventory.EntryFromPath(argPath, index)
		index++
		if !ok {
			if logf != nil {
				logf("skip arg path=%q missing, not a directory, or unreadable", argPath)
			}
			continue
		}
		if opts.SizeThreshold > 0 && entry.TotalBytes < opts.SizeThreshold {
			if logf != nil {
				logf("skip arg path=%s total_size=%s below threshold %s",
					entry.Path, entry.TotalSize, nminventory.FormatCompactHumanSize(opts.SizeThreshold))
			}
			continue
		}
		addEntry(entry)
	}

	entries := make([]nminventory.Entry, 0, len(order))
	for _, path := range order {
		entries = append(entries, byPath[path])
	}

	eligible, ineligible := nmmigrate.FilterEligible(entries, logf)
	selected := nminventory.ApplyLimit(eligible, opts.Limit)

	if logf != nil {
		logf("resolved: records=%t arg_paths=%d merged=%d ineligible=%d eligible=%d selected=%d",
			input.RecordsFile != "", len(input.Paths), len(byPath), ineligible, len(eligible), len(selected))
	}

	return selected, nil
}