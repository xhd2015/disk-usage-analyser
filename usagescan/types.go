package usagescan

import (
	"context"
	"io"
)

type TreeNode struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	Size     int64      `json:"size"`
	IsDir    bool       `json:"isDir"`
	Depth    int        `json:"depth"`
	Children []TreeNode `json:"children,omitempty"`
}

type TreeResult struct {
	Path      string   `json:"path"`
	TotalSize int64    `json:"totalSize"`
	Min       int64    `json:"min"`
	MaxDepth  int      `json:"maxDepth"`
	Tree      TreeNode `json:"tree"`
}

type ScanOptions struct {
	Min      int64
	MaxDepth int // 0 = unlimited
}

// Match is one node selected for the optional match / TOP section.
type Match struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"isDir"`
	Depth int    `json:"depth"`
}

// ViewResult is the machine-readable output of a view (inspect or live query).
type ViewResult struct {
	ScanPath   string   `json:"scanPath"`
	TotalSize  int64    `json:"totalSize"`
	Min        int64    `json:"min"`
	MaxDepth   int      `json:"maxDepth"`
	SourceFile string   `json:"sourceFile,omitempty"`
	Tree       TreeNode `json:"tree"`
	Matches    []Match  `json:"matches,omitempty"`
}

// ViewOptions configures phase-2 rendering over a TreeResult.
type ViewOptions struct {
	Min         int64
	MaxDepth    int
	Top         int  // ranking cap; 0 means default 20 when match section is active
	TopSet      bool // true when --top was provided
	AtPath      string
	Find        string
	Suffix      string
	IncludeRoot bool
	SourceFile  string
}

// TreeSource is phase 1: produce a TreeResult without rendering.
type TreeSource interface {
	Load() (TreeResult, error)
}

type Item struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"isDir"`
}

type Result struct {
	Path      string `json:"path"`
	TotalSize int64  `json:"totalSize"`
	Items     []Item `json:"items"`
}

type Options struct {
	Context context.Context
	OnItem  func(Item)
}

type CLIOptions struct {
	Stdout io.Writer
	Stderr io.Writer
}
