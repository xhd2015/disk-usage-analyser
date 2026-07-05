## Expected

- At least one `named`, `named_size`, and `named_enriched` event.
- Each `named` event JSON object contains key `gitTracked`.
- Each `named_size` event JSON object contains key `gitTracked`.
- Each `named_enriched` event JSON object contains keys:
  `packageManager`, `gitTracked`, `pnpmSharedSize`, `pnpmSharedHuman`,
  `bunSharedSize`, `bunSharedHuman`, `sharedSize`, `sharedHuman`.

## Errors

- No harness error is returned.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	namedObjs := namedEventDataObjects(t, resp.SSEOutput)
	if len(namedObjs) == 0 {
		t.Fatal("expected at least one named SSE event")
	}
	for i, obj := range namedObjs {
		if _, ok := obj["gitTracked"]; !ok {
			t.Fatalf("named event %d missing JSON key gitTracked in %#v", i, obj)
		}
	}

	sizeObjs := namedSizeEventDataObjects(t, resp.SSEOutput)
	if len(sizeObjs) == 0 {
		t.Fatal("expected at least one named_size SSE event")
	}
	for i, obj := range sizeObjs {
		if _, ok := obj["gitTracked"]; !ok {
			t.Fatalf("named_size event %d missing JSON key gitTracked in %#v", i, obj)
		}
	}

	enrichedObjs := namedEnrichedEventDataObjects(t, resp.SSEOutput)
	if len(enrichedObjs) == 0 {
		t.Fatal("expected at least one named_enriched SSE event")
	}
	required := []string{
		"packageManager",
		"gitTracked",
		"pnpmSharedSize",
		"pnpmSharedHuman",
		"bunSharedSize",
		"bunSharedHuman",
		"sharedSize",
		"sharedHuman",
	}
	for i, obj := range enrichedObjs {
		for _, key := range required {
			if _, ok := obj[key]; !ok {
				t.Fatalf("named_enriched event %d missing JSON key %q in %#v", i, key, obj)
			}
		}
	}
}
```