package usagescan

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LoadTreeResult reads a TreeResult JSON object from r.
func LoadTreeResult(r io.Reader) (TreeResult, error) {
	dec := json.NewDecoder(r)
	var result TreeResult
	if err := dec.Decode(&result); err != nil {
		return TreeResult{}, fmt.Errorf("decode TreeResult JSON: %w", err)
	}
	if result.Path == "" && result.Tree.Path == "" && result.Tree.Name == "" {
		return TreeResult{}, fmt.Errorf("decode TreeResult JSON: empty or invalid payload")
	}
	return result, nil
}

// LoadTreeResultFile reads TreeResult JSON from path, or stdin when path is "-" or "".
func LoadTreeResultFile(path string) (TreeResult, error) {
	if path == "" || path == "-" {
		return LoadTreeResult(os.Stdin)
	}
	f, err := os.Open(path)
	if err != nil {
		return TreeResult{}, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	return LoadTreeResult(f)
}

// FlattenTree returns every node in depth-first order (root first).
func FlattenTree(root TreeNode) []TreeNode {
	var out []TreeNode
	var walk func(TreeNode)
	walk = func(n TreeNode) {
		out = append(out, n)
		for _, c := range n.Children {
			walk(c)
		}
	}
	walk(root)
	return out
}

// FindNodeByPath returns the node whose path equals want (exact), or nil.
func FindNodeByPath(root TreeNode, want string) *TreeNode {
	want = filepath.Clean(want)
	var found *TreeNode
	var walk func(*TreeNode)
	walk = func(n *TreeNode) {
		if found != nil {
			return
		}
		if filepath.Clean(n.Path) == want {
			cp := *n
			found = &cp
			return
		}
		for i := range n.Children {
			walk(&n.Children[i])
		}
	}
	walk(&root)
	return found
}

func matchesFind(n TreeNode, find string) bool {
	if find == "" {
		return true
	}
	f := strings.ToLower(find)
	return strings.Contains(strings.ToLower(n.Path), f) ||
		strings.Contains(strings.ToLower(n.Name), f)
}

func matchesSuffix(n TreeNode, suffix string) bool {
	if suffix == "" {
		return true
	}
	return strings.HasSuffix(n.Name, suffix) || strings.HasSuffix(n.Path, suffix)
}

func expandLeadingTilde(p string) string {
	if p == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			return home
		}
		return p
	}
	if len(p) >= 2 && p[0] == '~' && (p[1] == '/' || p[1] == filepath.Separator) {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, p[2:])
		}
	}
	return p
}
