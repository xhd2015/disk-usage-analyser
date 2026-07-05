package nminventory

import "testing"

func TestParseCompactHumanSize(t *testing.T) {
	tests := []struct {
		in    string
		want  int64
		isErr bool
	}{
		{"0B", 0, false},
		{"10M", 10 * 1024 * 1024, false},
		{"10MB", 10 * 1024 * 1024, false},
		{"867.9MB", 910059110, false},
		{"4.6MB", 4823449, false},
		{"1.5GB", 1610612736, false},
		{"512B", 512, false},
		{"", 0, true},
		{"foo", 0, true},
	}

	for _, tc := range tests {
		got, err := ParseCompactHumanSize(tc.in)
		if tc.isErr {
			if err == nil {
				t.Fatalf("ParseCompactHumanSize(%q) expected error", tc.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("ParseCompactHumanSize(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("ParseCompactHumanSize(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}