## Expected
- There are exactly 17 software locations
- Each software location has RebootSafe=true
- Each has a non-empty Path, Label, Category distinct from core categories
- Go location has a Path containing "go/pkg/mod" and ExtraPaths containing "Caches/go-build"
- Xcode location has a Path containing "Xcode/DerivedData" and ExtraPaths containing "CoreSimulator"
- npm location has Path containing ".npm"
- Docker location has Path containing "com.docker.docker"
- Single-path tools like npm, Bun, Docker have empty or nil ExtraPaths

```go
import (
	"path/filepath"
	"strings"

	"disk-usage-analyser/server"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Locations) != 17 {
		t.Fatalf("expected 17 software locations, got %d", len(resp.Locations))
	}

	homeDir := req.HomeDir
	expectedMap := map[string]struct{
		pathSnippet string
		label       string
		hasExtras   bool
	}{
		"go":       {filepath.Join(homeDir, "go", "pkg", "mod"), "Go", true},
		"npm":      {filepath.Join(homeDir, ".npm"), "npm", false},
		"bun":      {filepath.Join(homeDir, ".bun", "install", "cache"), "Bun", false},
		"yarn":     {filepath.Join(homeDir, "Library", "Caches", "Yarn"), "Yarn", false},
		"pnpm":     {filepath.Join(homeDir, "Library", "pnpm", "store"), "pnpm", false},
		"pip":      {filepath.Join(homeDir, "Library", "Caches", "pip"), "pip", false},
		"cargo":    {filepath.Join(homeDir, ".cargo", "registry", "cache"), "Cargo", false},
		"ruby":     {filepath.Join(homeDir, ".gem"), "Ruby Gems", false},
		"docker":   {filepath.Join(homeDir, "Library", "Containers", "com.docker.docker"), "Docker", false},
		"podman":   {filepath.Join(homeDir, ".local", "share", "containers"), "Podman", false},
		"nginx":    {"/usr/local/var/log/nginx", "Nginx", false},
		"gradle":   {filepath.Join(homeDir, ".gradle", "caches"), "Gradle", false},
		"maven":    {filepath.Join(homeDir, ".m2", "repository"), "Maven", false},
		"android":  {filepath.Join(homeDir, "Library", "Android", "sdk"), "Android", false},
		"brew":     {filepath.Join(homeDir, "Library", "Caches", "Homebrew"), "Homebrew", false},
		"xcode":    {filepath.Join(homeDir, "Library", "Developer", "Xcode", "DerivedData"), "Xcode", true},
		"composer": {filepath.Join(homeDir, ".composer", "cache"), "Composer", false},
	}

	for cat, exp := range expectedMap {
		found := findLocation(resp.Locations, cat)
		if found == nil {
			t.Fatalf("missing software location: %s", cat)
		}
		if found.Path != exp.pathSnippet {
			t.Fatalf("location %s: expected Path=%s, got %s", cat, exp.pathSnippet, found.Path)
		}
		if found.Label != exp.label {
			t.Fatalf("location %s: expected Label=%s, got %s", cat, exp.label, found.Label)
		}
		if !found.RebootSafe {
			t.Fatalf("location %s: expected RebootSafe=true", cat)
		}
		if exp.hasExtras {
			if len(found.ExtraPaths) == 0 {
				t.Fatalf("location %s: expected ExtraPaths to be non-empty", cat)
			}
		}
	}
}

func findLocation(locations []server.TmpLocation, category string) *server.TmpLocation {
	for i := range locations {
		if locations[i].Category == category {
			return &locations[i]
		}
	}
	return nil
}
```
