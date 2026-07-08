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
	Threshold int64    `json:"threshold"`
	MaxDepth  int      `json:"maxDepth"`
	Tree      TreeNode `json:"tree"`
}

type ScanOptions struct {
	Threshold int64
	MaxDepth  int // 0 = unlimited
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