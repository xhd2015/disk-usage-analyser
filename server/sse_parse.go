package server

import (
	"encoding/json"
	"strings"
)

// BinariesSSEResult holds parsed tmp-binaries-scan SSE output.
type BinariesSSEResult struct {
	SSEOutput  string
	EventTypes []string
	Binaries   []BinaryHit
	Summary    *BinaryScanSummary
}

// WorktreesSSEResult holds parsed tmp-worktrees-scan SSE output.
type WorktreesSSEResult struct {
	SSEOutput  string
	EventTypes []string
	Worktrees  []WorktreeHit
	Repos      int
	Summary    *WorktreeScanSummary
}

// ParseBinariesSSE parses SSE body from tmp-binaries-scan.
func ParseBinariesSSE(body string) (*BinariesSSEResult, error) {
	resp := &BinariesSSEResult{SSEOutput: body}
	lines := strings.Split(body, "\n")
	var currentEvent string
	for _, line := range lines {
		if strings.HasPrefix(line, "event: ") {
			currentEvent = strings.TrimPrefix(line, "event: ")
			resp.EventTypes = append(resp.EventTypes, currentEvent)
		}
		if strings.HasPrefix(line, "data: ") && currentEvent != "" {
			data := strings.TrimPrefix(line, "data: ")
			switch currentEvent {
			case "binary":
				var hit BinaryHit
				if err := json.Unmarshal([]byte(data), &hit); err == nil {
					resp.Binaries = append(resp.Binaries, hit)
				}
			case "summary":
				var summary BinaryScanSummary
				if err := json.Unmarshal([]byte(data), &summary); err == nil {
					resp.Summary = &summary
				}
			}
		}
	}
	return resp, nil
}

// ParseWorktreesSSE parses SSE body from tmp-worktrees-scan.
func ParseWorktreesSSE(body string) (*WorktreesSSEResult, error) {
	resp := &WorktreesSSEResult{SSEOutput: body}
	lines := strings.Split(body, "\n")
	var currentEvent string
	for _, line := range lines {
		if strings.HasPrefix(line, "event: ") {
			currentEvent = strings.TrimPrefix(line, "event: ")
			resp.EventTypes = append(resp.EventTypes, currentEvent)
			if currentEvent == "repo" {
				resp.Repos++
			}
		}
		if strings.HasPrefix(line, "data: ") && currentEvent != "" {
			data := strings.TrimPrefix(line, "data: ")
			switch currentEvent {
			case "worktree":
				var hit WorktreeHit
				if err := json.Unmarshal([]byte(data), &hit); err == nil {
					resp.Worktrees = append(resp.Worktrees, hit)
				}
			case "summary":
				var summary WorktreeScanSummary
				if err := json.Unmarshal([]byte(data), &summary); err == nil {
					resp.Summary = &summary
				}
			}
		}
	}
	return resp, nil
}

// EventBefore reports whether event `before` appears before `after` in SSE event order.
func EventBefore(events []string, before, after string) bool {
	beforeIdx, afterIdx := -1, -1
	for i, e := range events {
		if e == before && beforeIdx < 0 {
			beforeIdx = i
		}
		if e == after {
			afterIdx = i
		}
	}
	return beforeIdx >= 0 && afterIdx >= 0 && beforeIdx < afterIdx
}