## Expected
- There are exactly 17 software locations
- Each software location has RebootSafe=true
- Each has a non-empty Path, Label, Category distinct from core categories
- All home-relative paths use `~` prefix (e.g., `~/go/pkg/mod` instead of `/Users/testuser/go/pkg/mod`)
- Go location has Path `~/go/pkg/mod` and ExtraPaths containing `~/Library/Caches/go-build`
- Xcode location has Path `~/Library/Developer/Xcode/DerivedData` and ExtraPaths containing `~/Library/Developer/CoreSimulator/Devices`
- Nginx path remains absolute (not under home dir): `/usr/local/var/log/nginx`
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

	expectedMap := map[string]struct{
		pathSnippet string
		label       string
		hasExtras   bool
		extraCheck  string // substring expected in ExtraPaths[0] (empty = no check)
	}{
		"go":       {"~/go/pkg/mod", "Go", true, "~/Library/Caches/go-build"},
		"npm":      {"~/.npm", "npm", false, ""},
		"bun":      {"~/.bun/install/cache", "Bun", false, ""},
		"yarn":     {"~/Library/Caches/Yarn", "Yarn", false, ""},
		"pnpm":     {"~/Library/pnpm/store", "pnpm", false, ""},
		"pip":      {"~/Library/Caches/pip", "pip", false, ""},
		"cargo":    {"~/.cargo/registry/cache", "Cargo", false, ""},
		"ruby":     {"~/.gem", "Ruby Gems", false, ""},
		"docker":   {"~/Library/Containers/com.docker.docker", "Docker", false, ""},
		"podman":   {"~/.local/share/containers", "Podman", false, ""},
		"nginx":    {"/usr/local/var/log/nginx", "Nginx", false, ""},
		"gradle":   {"~/.gradle/caches", "Gradle", false, ""},
		"maven":    {"~/.m2/repository", "Maven", false, ""},
		"android":  {"~/Library/Android/sdk", "Android", false, ""},
		"brew":     {"~/Library/Caches/Homebrew", "Homebrew", false, ""},
		"xcode":    {"~/Library/Developer/Xcode/DerivedData", "Xcode", true, "~/Library/Developer/CoreSimulator/Devices"},
		"composer": {"~/.composer/cache", "Composer", false, ""},
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
			if exp.extraCheck != "" {
				if found.ExtraPaths[0] != exp.extraCheck {
					t.Fatalf("location %s: expected ExtraPaths[0]=%s, got %s", cat, exp.extraCheck, found.ExtraPaths[0])
				}
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
