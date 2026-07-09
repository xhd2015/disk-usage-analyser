package usagescan

// LiveTreeSource walks the filesystem to produce a TreeResult.
type LiveTreeSource struct {
	Path string
	Opts ScanOptions
}

func (s LiveTreeSource) Load() (TreeResult, error) {
	return Scan(s.Path, s.Opts)
}

// JSONTreeSource loads a previously captured TreeResult JSON (field "min").
// Path may be "-" for stdin.
type JSONTreeSource struct {
	Path string
}

func (s JSONTreeSource) Load() (TreeResult, error) {
	return LoadTreeResultFile(s.Path)
}
