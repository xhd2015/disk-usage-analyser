package nminventory

import "disk-usage-analyser/tmpfiles"

// MeasureNodeModulesSize returns recursive byte size using the same nested skip
// rules as the node_modules inventory scan.
func MeasureNodeModulesSize(nodeModulesPath string) (int64, error) {
	names := map[string]bool{"node_modules": true}
	return tmpfiles.DirSize(nodeModulesPath, names)
}