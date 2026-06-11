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
	d, err := os.MkdirTemp("", "setup-check")
	if err != nil {
		return err
	}
	os.RemoveAll(d)
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	dir, err := os.MkdirTemp("", "scan-err-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Chmod(filepath.Join(dir, "locked"), 0700)
		os.RemoveAll(dir)
	}()

	os.WriteFile(filepath.Join(dir, "a.txt"), make([]byte, 100), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), make([]byte, 200), 0644)

	locked := filepath.Join(dir, "locked")
	os.Mkdir(locked, 0700)
	os.WriteFile(filepath.Join(locked, "c.txt"), make([]byte, 50), 0644)
	os.Chmod(locked, 0000)

	fsys := os.DirFS(dir)
	var progressSizes []int64
	size, count, scanErr := server.ScanWithProgress(fsys, ".", func(s int64, c int64) {
		progressSizes = append(progressSizes, s)
	})
	hadErr := int64(0)
	if scanErr != nil {
		hadErr = 1
	}
	return &Response{
		Size:            size,
		FileCount:       count,
		TotalSize:       int64(len(progressSizes)),
		ReclaimableSize: hadErr,
	}, nil
}
```
