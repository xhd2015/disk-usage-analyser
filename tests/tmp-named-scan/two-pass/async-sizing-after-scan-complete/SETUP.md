# Scenario

**Leaf**: async `emitNamedDirSized` goroutines must finish before `Scan` returns

Reproduces lifecycle bug: `OnNamedPreview` spawns async sizing; `tmpfiles.Scan`
returns while goroutines still call `OnNamedHit`. The handler then `close(jobCh)`
and late sends panic with `send on closed channel`.

## Steps

1. Create git repos with sizable `node_modules` trees and nested
   `node_modules/path-scurry/node_modules` (matches production stack trace).
2. Run `named-scan-lifecycle` against the fake home.

```go
func nodeModulesWithNested(t *testing.T, home, rel string, fileCount int) {
	t.Helper()
	app := repo(t, home, rel)
	nm := filepath.Join(app, "node_modules")
	for i := 0; i < fileCount; i++ {
		pkg := "pkg" + string(rune('a'+i%26)) + string(rune('0'+(i/26)%10))
		writeFile(t, nm, pkg+"/data.bin", make([]byte, 16384), 0644)
	}
	nested := filepath.Join(nm, "path-scurry", "node_modules")
	writeFile(t, nested, "inner/x.txt", []byte("x"), 0644)
}

func Setup(t *testing.T, req *Request) error {
	req.Op = "named-scan-lifecycle"
	// Single repo with sizable node_modules: walk returns in <1ms but async sizing
	// takes tens of ms, reproducing post-scan OnNamedHit reliably.
	repos := []string{"Projects/race-app-00"}
	for _, rel := range repos {
		nodeModulesWithNested(t, req.HomeDir, rel, 100)
	}
	return nil
}
```