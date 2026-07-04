//go:build darwin

package analyse

import (
	"io/fs"
	"path/filepath"
)

func buildCloneStoreIndex(storeDir string) map[cloneGroupKey]struct{} {
	index := make(map[cloneGroupKey]struct{})
	_ = filepath.WalkDir(storeDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || path == storeDir {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}
		docID, cloneID, _, ok := fileCloneAttrs(path)
		if !ok {
			return nil
		}
		key, ok := cloneGroupKeyFrom(docID, cloneID)
		if !ok {
			return nil
		}
		index[key] = struct{}{}
		return nil
	})
	return index
}

func sharedBytesInNodeModules(nodeModulesPath string, storeIndex map[cloneGroupKey]struct{}) int64 {
	if len(storeIndex) == 0 {
		return 0
	}

	var total int64
	counted := make(map[cloneGroupKey]struct{})

	_ = filepath.WalkDir(nodeModulesPath, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || path == nodeModulesPath {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}
		docID, cloneID, _, ok := fileCloneAttrs(path)
		if !ok {
			return nil
		}
		key, ok := cloneGroupKeyFrom(docID, cloneID)
		if !ok {
			return nil
		}
		if _, inStore := storeIndex[key]; !inStore {
			return nil
		}
		if _, seen := counted[key]; seen {
			return nil
		}
		counted[key] = struct{}{}
		total += info.Size()
		return nil
	})
	return total
}