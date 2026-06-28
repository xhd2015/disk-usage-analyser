# Scenario

**Feature**: Each software location has correct Path, Label, Category, ExtraPaths

```
# handler discovers locations, scans paths, streams SSE
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events
```

## Preconditions
- A home directory path of "/Users/testuser" is provided
- 17 software cache/log locations are defined for macOS

## Steps
1. Set req.HomeDir to "/Users/testuser"
2. Call DiscoverLocations with the home directory
3. Filter to software-only locations (those with categories other than trash/cache/log/temp)
4. Verify each has correct Path, Label, Category, RebootSafe=true, and optional ExtraPaths

## Context
- Software locations are: Go, npm, Bun, Yarn, pnpm, pip, Cargo, Ruby Gems, Docker, Podman, Nginx, Gradle, Maven, Android, Homebrew, Xcode, Composer
- Go has ExtraPaths: ~/Library/Caches/go-build
- Xcode has ExtraPaths: ~/Library/Developer/CoreSimulator/Devices
- All software locations have RebootSafe=true (they persist across reboots)

```go
import (
	"path/filepath"
	"strings"

	"disk-usage-analyser/server"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "software-locations"
	req.HomeDir = "/Users/testuser"
	return nil
}

```
