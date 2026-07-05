## Expected

- At least one `named_enriched` hit for the nested pnpm `node_modules` path.
- Every enriched hit has `gitTracked=false`.

## Errors

- No harness error is returned.

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.NamedEnrichedJSON) == 0 {
		t.Fatal("expected at least one named_enriched hit")
	}
	var nested *namedEventJSON
	for i := range resp.NamedEnrichedJSON {
		if strings.Contains(resp.NamedEnrichedJSON[i].Path, ".pnpm") {
			nested = &resp.NamedEnrichedJSON[i]
			break
		}
	}
	if nested == nil {
		t.Fatalf("expected nested pnpm named_enriched hit, got %#v", resp.NamedEnrichedJSON)
	}
	if nested.GitTracked {
		t.Fatalf("nested pnpm gitTracked = true, want false for path %q", nested.Path)
	}
}
```