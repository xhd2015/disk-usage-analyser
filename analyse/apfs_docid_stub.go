//go:build !darwin

package analyse

type cloneGroupTracker struct{}

func newCloneGroupTracker() *cloneGroupTracker {
	return &cloneGroupTracker{}
}

func (t *cloneGroupTracker) CountSize(path string, size int64) int64 {
	return size
}

func (t *cloneGroupTracker) Add(path string, inode inodeKey, size int64) {}

func (t *cloneGroupTracker) TotalSharedCloneSize() int64 {
	return 0
}