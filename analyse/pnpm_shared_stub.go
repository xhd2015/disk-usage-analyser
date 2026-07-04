//go:build !darwin

package analyse

type pnpmSharedContext struct{}

func newPnpmSharedContext() *pnpmSharedContext {
	return &pnpmSharedContext{}
}

func (c *pnpmSharedContext) pnpmSharedForNodeModules(nodeModulesPath string) int64 {
	return 0
}