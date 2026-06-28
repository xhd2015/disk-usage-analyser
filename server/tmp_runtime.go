package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xhd2015/xgo/support/cmd"
)

type TmpRuntimeItem struct {
	Type        string `json:"type"`
	TotalCount  int64  `json:"totalCount"`
	ActiveCount int64  `json:"activeCount"`
	Size        int64  `json:"size"`
	Reclaimable int64  `json:"reclaimable,omitempty"`
}

type RuntimeCommandRunner func(name string, args ...string) ([]byte, error)

var runtimeCommandRunner RuntimeCommandRunner

func SetRuntimeCommandRunner(runner RuntimeCommandRunner) {
	runtimeCommandRunner = runner
}

func defaultRuntimeCommandRunner(name string, args ...string) ([]byte, error) {
	out, err := cmd.Debug().Output(name, args...)
	return []byte(out), err
}

func getRuntimeCommandRunner() RuntimeCommandRunner {
	if runtimeCommandRunner != nil {
		return runtimeCommandRunner
	}
	return defaultRuntimeCommandRunner
}

type systemDFRecord struct {
	Active      string `json:"Active"`
	Reclaimable string `json:"Reclaimable"`
	Size        string `json:"Size"`
	TotalCount  string `json:"TotalCount"`
	Type        string `json:"Type"`
}

func ParseSystemDFJSON(output string) ([]TmpRuntimeItem, error) {
	var items []TmpRuntimeItem
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var rec systemDFRecord
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			return nil, fmt.Errorf("parse system df line: %w", err)
		}
		item, err := recordToRuntimeItem(rec)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func recordToRuntimeItem(rec systemDFRecord) (TmpRuntimeItem, error) {
	totalCount, err := strconv.ParseInt(rec.TotalCount, 10, 64)
	if err != nil {
		return TmpRuntimeItem{}, fmt.Errorf("parse TotalCount %q: %w", rec.TotalCount, err)
	}
	activeCount, err := strconv.ParseInt(rec.Active, 10, 64)
	if err != nil {
		return TmpRuntimeItem{}, fmt.Errorf("parse Active %q: %w", rec.Active, err)
	}
	size, err := ParseHumanSize(rec.Size)
	if err != nil {
		return TmpRuntimeItem{}, err
	}
	reclaimable, err := ParseHumanSize(rec.Reclaimable)
	if err != nil {
		return TmpRuntimeItem{}, err
	}
	itemType := rec.Type
	if itemType == "Image" {
		itemType = "Images"
	}
	return TmpRuntimeItem{
		Type:        itemType,
		TotalCount:  totalCount,
		ActiveCount: activeCount,
		Size:        size,
		Reclaimable: reclaimable,
	}, nil
}

func ParseHumanSize(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("empty size string")
	}

	upper := strings.ToUpper(s)
	if upper == "0" || upper == "0B" {
		return 0, nil
	}

	units := []struct {
		suffix string
		mult   int64
	}{
		{"TB", 1000000000000},
		{"GB", 1000000000},
		{"MB", 1000000},
		{"KB", 1000},
		{"B", 1},
	}

	for _, u := range units {
		if strings.HasSuffix(upper, u.suffix) {
			numStr := strings.TrimSpace(s[:len(s)-len(u.suffix)])
			if numStr == "" {
				return 0, fmt.Errorf("invalid size string: %q", s)
			}
			val, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid size string: %q", s)
			}
			return int64(val * float64(u.mult)), nil
		}
	}

	return 0, fmt.Errorf("invalid size string: %q", s)
}

func FilterRuntimeItems(items []TmpRuntimeItem, types ...string) []TmpRuntimeItem {
	allowed := make(map[string]bool, len(types))
	for _, t := range types {
		allowed[t] = true
	}
	var filtered []TmpRuntimeItem
	for _, item := range items {
		if allowed[item.Type] {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func CollectRuntimeStats(runtime string) ([]TmpRuntimeItem, error) {
	runner := getRuntimeCommandRunner()
	output, err := runner(runtime, "system", "df", "--format", "json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, nil
	}
	items, err := ParseSystemDFJSON(string(output))
	if err != nil {
		return nil, nil
	}
	return FilterRuntimeItems(items, "Images", "Build Cache"), nil
}