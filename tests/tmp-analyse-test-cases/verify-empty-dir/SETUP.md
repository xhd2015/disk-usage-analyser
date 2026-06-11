## Preconditions
- An empty mock filesystem exists with no files

## Steps
1. Create an empty mock filesystem using testing/fstest.MapFS
2. Set req.FS to the empty filesystem
3. CalculateSize should return size=0, fileCount=0, and no error

```go
import (
	"io/fs"
	"testing/fstest"
)

func Setup(t *testing.T, req *Request) error {
	req.FS = fstest.MapFS{}
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	fsys := req.FS.(fs.FS)
	size, count, err := server.CalculateSize(fsys, ".")
	if err != nil {
		return nil, err
	}
	return &Response{Size: size, FileCount: count}, nil
}
```
