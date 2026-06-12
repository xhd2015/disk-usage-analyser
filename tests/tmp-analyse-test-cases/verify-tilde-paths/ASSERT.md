## Expected
- All paths that are under the home directory use `~` prefix (e.g., `~/Library/Caches` NOT `/Users/testuser/Library/Caches`)
- Paths NOT under the home directory (e.g., `/tmp`, `/usr/local/var/log/nginx`) remain as absolute paths
- No path contains the raw home directory string `/Users/testuser`
- All ExtraPaths also use `~` prefix

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Locations) < 22 {
		t.Fatalf("expected at least 22 locations, got %d", len(resp.Locations))
	}

	homeDir := req.HomeDir
	absPathDirs := map[string]bool{
		"/tmp": true,
		"/usr/local/var/log/nginx": true,
	}

	for _, loc := range resp.Locations {
		// Check primary path
		if loc.Path == "" {
			t.Fatal("location has empty Path")
		}
		if strings.Contains(loc.Path, homeDir) {
			t.Fatalf("path %q should not contain raw home directory %q — use ~ prefix", loc.Path, homeDir)
		}
		if !absPathDirs[loc.Path] && !strings.HasPrefix(loc.Path, "~/") {
			t.Fatalf("home-relative path %q should start with ~/", loc.Path)
		}
		if loc.Path != "/tmp" && strings.HasPrefix(loc.Path, "/") && !strings.HasPrefix(loc.Path, "~/") {
			// Check it's one of the known absolute paths
			if !absPathDirs[loc.Path] {
				t.Fatalf("path %q should use ~ prefix or be a known absolute path", loc.Path)
			}
		}

		// Check ExtraPaths
		for i, ep := range loc.ExtraPaths {
			if strings.Contains(ep, homeDir) {
				t.Fatalf("ExtraPaths[%d]=%q should not contain raw home directory %q — use ~ prefix", i, ep, homeDir)
			}
			if !strings.HasPrefix(ep, "~/") {
				t.Fatalf("ExtraPaths[%d]=%q should start with ~/", i, ep)
			}
		}
	}
}
```
