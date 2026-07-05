# Scenario

**Leaf**: client closes SSE after `scan_complete` while enrichment workers still run

Reproduces production log noise when late `named_enriched` writes hit a closed
connection:

```
Error sending SSE event named_enriched: write tcp ...: write: broken pipe
```

## Steps

1. Create multiple git repos with sizable `node_modules` trees so enrichment
   continues after `scan_complete`.
2. Run `named-scan-disconnect`: read until `scan_complete`, close body, drain.

```go
import "fmt"

func pnpmLikeNodeModulesBacklog(t *testing.T, home, rel string, pkgCount, filesPerPkg int) {
	t.Helper()
	app := repo(t, home, rel)
	for i := 0; i < pkgCount; i++ {
		pkg := fmt.Sprintf("pkg%d", i)
		base := filepath.Join("node_modules", ".pnpm", pkg+"@1.0.0", "node_modules", pkg)
		for j := 0; j < filesPerPkg; j++ {
			writeSizedFile(t, app, filepath.Join(base, fmt.Sprintf("data%d.bin", j)), 65536)
		}
	}
}

func Setup(t *testing.T, req *Request) error {
	req.Op = "named-scan-disconnect"
	repos := []string{
		"Projects/disconnect-app-a", "Projects/disconnect-app-b",
		"Projects/disconnect-app-c", "Projects/disconnect-app-d",
	}
	for _, rel := range repos {
		pnpmLikeNodeModulesBacklog(t, req.HomeDir, rel, 36, 120)
	}
	return nil
}
```