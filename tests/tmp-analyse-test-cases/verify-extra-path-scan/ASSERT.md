## Expected
- Primary size = 300 (100 + 200)
- Primary file count = 2
- ExtraSizes has length 1 with value 800 (500 + 300)
- ExtraCounts has length 1 with value 2

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Size != 300 {
		t.Fatalf("expected primary Size=300, got %d", resp.Size)
	}
	if resp.FileCount != 2 {
		t.Fatalf("expected primary FileCount=2, got %d", resp.FileCount)
	}
	if len(resp.ExtraSizes) != 1 {
		t.Fatalf("expected 1 ExtraSize, got %d", len(resp.ExtraSizes))
	}
	if resp.ExtraSizes[0] != 800 {
		t.Fatalf("expected ExtraSize[0]=800, got %d", resp.ExtraSizes[0])
	}
	if len(resp.ExtraCounts) != 1 {
		t.Fatalf("expected 1 ExtraCount, got %d", len(resp.ExtraCounts))
	}
	if resp.ExtraCounts[0] != 2 {
		t.Fatalf("expected ExtraCount[0]=2, got %d", resp.ExtraCounts[0])
	}
}
```
