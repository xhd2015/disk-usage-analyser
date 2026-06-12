## Preconditions
- A home directory path of "/Users/testuser" is provided

## Steps
1. Set req.HomeDir to "/Users/testuser"
2. Call DiscoverLocations with the home directory
3. Return the discovered locations, counting categories and detected/not-detected status

## Context
- macOS has standard locations: ~/.Trash (user trash), ~/Library/Caches (user caches), ~/Library/Logs (user logs), $TMPDIR (user temp), /tmp (system temp)
- Additionally, 17 software cache/log locations are discovered: Go, npm, Bun, Yarn, pnpm, pip, Cargo, Ruby Gems, Docker, Podman, Nginx, Gradle, Maven, Android, Homebrew, Xcode, Composer
- Total locations should be at least 22
- Each location must have a Label, Category, Detected field, and RebootSafe flag
- Category must be a non-empty string
- Core locations always have Detected=true
- Software locations have Detected based on whether the path exists on disk

```go
import (
	"disk-usage-analyser/server"
)

func Setup(t *testing.T, req *Request) error {
	req.HomeDir = "/Users/testuser"
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	locations := server.DiscoverLocations(req.HomeDir)
	categoryCount := make(map[string]int)
	detectedCount := 0
	notDetectedCount := 0
	for _, loc := range locations {
		categoryCount[loc.Category]++
		if loc.Detected {
			detectedCount++
		} else {
			notDetectedCount++
		}
	}
	return &Response{
		Locations:        locations,
		CategoryCount:    categoryCount,
		DetectedCount:    detectedCount,
		NotDetectedCount: notDetectedCount,
	}, nil
}
```
