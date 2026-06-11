## Preconditions
- A home directory path of "/Users/testuser" is provided

## Steps
1. Set req.HomeDir to "/Users/testuser"
2. Call DiscoverLocations with the home directory
3. Return the discovered locations in the response

## Context
- macOS has standard locations: ~/.Trash (user trash), ~/Library/Caches (user caches), ~/Library/Logs (user logs), $TMPDIR (user temp), /tmp (system temp)
- Each location must have a Label, Category, and RebootSafe flag
- Category must be one of: "trash", "temp", "cache", "log"

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
	for _, loc := range locations {
		categoryCount[loc.Category]++
	}
	return &Response{
		Locations:     locations,
		CategoryCount: categoryCount,
	}, nil
}
```
