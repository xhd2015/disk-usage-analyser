package usagescan

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

func Scan(path string, opts ScanOptions) (TreeResult, error) {
	return ScanTree(path, opts)
}

func ScanTree(path string, opts ScanOptions) (TreeResult, error) {
	ctx := context.Background()

	absPath, err := filepath.Abs(path)
	if err != nil {
		return TreeResult{}, fmt.Errorf("invalid path %s: %w", path, err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return TreeResult{}, fmt.Errorf("path does not exist: %s", absPath)
		}
		return TreeResult{}, fmt.Errorf("cannot access %s: %w", absPath, err)
	}
	if !info.IsDir() {
		return TreeResult{}, fmt.Errorf("not a directory: %s", absPath)
	}

	totalSize := getDirSizeWithCache(ctx, absPath, nil)
	tree := buildTreeNode(ctx, absPath, ".", 0, opts)

	return TreeResult{
		Path:      absPath,
		TotalSize: totalSize,
		Threshold: opts.Threshold,
		MaxDepth:  opts.MaxDepth,
		Tree:      tree,
	}, nil
}

func buildTreeNode(ctx context.Context, absPath, name string, depth int, opts ScanOptions) TreeNode {
	size := getDirSizeWithCache(ctx, absPath, nil)

	node := TreeNode{
		Name:  name,
		Path:  absPath,
		Size:  size,
		IsDir: true,
		Depth: depth,
	}

	canExpand := opts.MaxDepth == 0 || depth < opts.MaxDepth
	if !canExpand {
		return node
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return node
	}

	var subDirs []fs.DirEntry
	var files []fs.DirEntry
	for _, entry := range entries {
		if entry.IsDir() {
			subDirs = append(subDirs, entry)
		} else {
			files = append(files, entry)
		}
	}

	children := make([]TreeNode, 0, len(entries))

	for _, entry := range files {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		fileSize := info.Size()
		if !passesThreshold(fileSize, opts.Threshold) {
			continue
		}
		childPath := filepath.Join(absPath, entry.Name())
		children = append(children, TreeNode{
			Name:  entry.Name(),
			Path:  childPath,
			Size:  fileSize,
			IsDir: false,
			Depth: depth + 1,
		})
	}

	for _, entry := range subDirs {
		childPath := filepath.Join(absPath, entry.Name())
		childSize := getDirSizeWithCache(ctx, childPath, nil)
		if !passesThreshold(childSize, opts.Threshold) {
			continue
		}
		child := buildTreeNode(ctx, childPath, entry.Name(), depth+1, opts)
		children = append(children, child)
	}

	sortTreeChildren(children)
	node.Children = children
	return node
}

func passesThreshold(size, threshold int64) bool {
	return threshold == 0 || size >= threshold
}

func sortTreeChildren(children []TreeNode) {
	sort.Slice(children, func(i, j int) bool {
		if children[i].Size != children[j].Size {
			return children[i].Size > children[j].Size
		}
		if children[i].IsDir != children[j].IsDir {
			return children[i].IsDir
		}
		return children[i].Name < children[j].Name
	})
}