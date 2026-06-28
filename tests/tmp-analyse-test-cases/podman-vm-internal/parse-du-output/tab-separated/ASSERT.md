## Expected

- `1735229168\t/path` parses to 1735229168 bytes.

## Side Effects

- None (pure function).

## Errors

- Parse must succeed; no error returned from Run.

## Exit Code

- Test passes when parsed bytes match.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ParseFailed {
		t.Fatalf("expected successful parse, got error: %v", resp.Err)
	}
	if resp.ParsedBytes != 1735229168 {
		t.Fatalf("expected 1735229168 bytes, got %d", resp.ParsedBytes)
	}
}
```