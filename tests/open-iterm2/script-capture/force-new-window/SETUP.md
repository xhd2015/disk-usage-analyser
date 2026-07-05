# Scenario

**Leaf**: POST `~/.../node_modules` opens parent project dir in a new iTerm2 window

## Steps

1. Create `~/Projects/iterm-app/node_modules/` fixture.
2. POST tilde path ending with `node_modules`.
3. Capture AppleScript via `KOOL_ITERM2_SCRIPT_OUT`.

```go
func Setup(t *testing.T, req *Request) error {
	writeDirFixture(t, req.HomeDir, "Projects/iterm-app")
	req.Path = nodeModulesTildePath(req.HomeDir, "Projects/iterm-app")
	return nil
}
```