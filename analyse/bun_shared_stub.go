//go:build !darwin

package analyse

type bunSharedContext struct{}

func newBunSharedContext() *bunSharedContext {
	return &bunSharedContext{}
}

func (c *bunSharedContext) bunSharedForNodeModules(nodeModulesPath string) int64 {
	return 0
}