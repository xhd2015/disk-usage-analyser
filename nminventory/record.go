package nminventory

// Record is one node_modules inventory entry (JSON object shape).
type Record struct {
	Path           string `json:"path"`
	HasPackageJSON bool   `json:"has_package_json"`
	PackageManager string `json:"package_manager"`
	SharedSize     string `json:"shared_size"`
	TotalSize      string `json:"total_size"`
	BelongsToGit   bool   `json:"belongs_to_git"`
}

// OutputFile is the versioned inventory JSON envelope.
type OutputFile struct {
	Version     string   `json:"version"`
	NodeModules []Record `json:"node_modules"`
}