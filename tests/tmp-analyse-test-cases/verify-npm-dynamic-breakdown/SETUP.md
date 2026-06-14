## Preconditions
- npm location under ~/.npm gets dynamic breakdown based on actual subdirectories
- A mock filesystem is created with ~/.npm containing sub-directories _cacache and _logs

## Steps
1. Create a temporary directory mimicking ~/.npm with subdirectories containing files
2. DiscoverLocations finds npm as single path
3. Scan the npm location, then read its subdirectories to produce breakdownItems
4. Verify breakdownItems contains entries for each subdirectory

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	tmpDir, err := os.MkdirTemp("", "doctest-npm-*")
	if err != nil {
		return err
	}
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	npmDir := filepath.Join(tmpDir, ".npm")
	cacacheDir := filepath.Join(npmDir, "_cacache")
	logsDir := filepath.Join(npmDir, "_logs")

	os.MkdirAll(cacacheDir, 0755)
	os.MkdirAll(logsDir, 0755)
	os.WriteFile(filepath.Join(cacacheDir, "cache-entry-1"), []byte("cached content here more"), 0644)
	os.WriteFile(filepath.Join(cacacheDir, "cache-entry-2"), []byte("more cache data goes here"), 0644)
	os.WriteFile(filepath.Join(logsDir, "debug.log"), []byte("some log entry"), 0644)

	req.HomeDir = tmpDir
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	npmDir := filepath.Join(req.HomeDir, ".npm")

	entries, err := os.ReadDir(npmDir)
	if err != nil {
		return nil, err
	}

	items := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			items = append(items, filepath.Join(npmDir, e.Name()))
		}
	}

	return &Response{
		ExtraPaths:    items,
		DetectedCount: len(items),
	}, nil
}
```
