package nminventory

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseCompactHumanSize parses compact binary sizes such as 10M, 10MB, 867.9MB, 0B.
func ParseCompactHumanSize(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty size string")
	}

	upper := strings.ToUpper(s)
	if upper == "0" || upper == "0B" {
		return 0, nil
	}

	units := []struct {
		suffix string
		mult   int64
	}{
		{"TB", 1024 * 1024 * 1024 * 1024},
		{"GB", 1024 * 1024 * 1024},
		{"MB", 1024 * 1024},
		{"M", 1024 * 1024},
		{"KB", 1024},
		{"K", 1024},
		{"B", 1},
	}

	for _, u := range units {
		if !strings.HasSuffix(upper, u.suffix) {
			continue
		}
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

	return 0, fmt.Errorf("invalid size string: %q", s)
}

// FormatCompactHumanSize formats bytes as compact binary sizes (e.g. 13MB, 0B).
func FormatCompactHumanSize(size int64) string {
	if size <= 0 {
		return "0B"
	}
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	}
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	value := float64(size)
	for _, unit := range units {
		value /= 1024
		if value < 1024 {
			if value == float64(int64(value)) {
				return fmt.Sprintf("%.0f%s", value, unit)
			}
			return fmt.Sprintf("%.1f%s", value, unit)
		}
	}
	return fmt.Sprintf("%.1fEB", value/1024)
}