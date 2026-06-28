# Scenario

**Feature**: Error handling during scan

```
# handler discovers locations, scans paths, streams SSE
Client -> HandleTmpAnalyse -> DiscoverLocations -> ScanWithProgress -> SSE events
```

## Preconditions
- A real temp directory with accessible files and a locked subdirectory (chmod 0000)

## Steps
1. Create temp dir with 2 accessible files: a.txt (100 bytes), b.txt (200 bytes)
2. Create locked subdirectory with chmod 0000 containing c.txt
3. Run ScanWithProgress — should accumulate a.txt+b.txt, then error on locked dir
4. Verify accumulated size=300, count=2 (not 0) and error is not nil

```go
import (
	"io/fs"
	"os"
	"path/filepath"

	"disk-usage-analyser/server"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "scan-with-partial-error"
	d, err := os.MkdirTemp("", "setup-check")
	if err != nil {
		return err
	}
	os.RemoveAll(d)
	return nil
}

```
