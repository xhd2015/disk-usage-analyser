# Scenario

**Leaf**: non-darwin platform returns HTTP 501 Not Implemented

## Steps

1. Set `KOOL_ITERM2_GOOS=linux`.
2. POST a valid tilde project path.

```go
func Setup(t *testing.T, req *Request) error {
	writeDirFixture(t, req.HomeDir, "Projects/linux-app")
	req.Path = tildePath(req.HomeDir, projectParentAbs(req.HomeDir, "Projects/linux-app"))
	return nil
}
```