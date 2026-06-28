## Expected

- Output with no byte total fails parsing (`ParseFailed` true).

## Side Effects

- None (pure function).

## Errors

- Run returns `ParseFailed` without fatal Run error.

## Exit Code

- Test passes when parse failure is reported.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	if !resp.ParseFailed {
		t.Fatal("expected parse failure for garbage du output")
	}
	if resp.ParsedBytes != 0 {
		t.Fatalf("expected zero bytes on failure, got %d", resp.ParsedBytes)
	}
}
```