## Expected

- Xcode primary Path is `~/Library/Developer/Xcode/DerivedData`.
- Xcode has exactly 4 ExtraPaths in this order:
  1. `~/Library/Developer/CoreSimulator/Devices`
  2. `~/Library/Developer/Xcode/iOS DeviceSupport`
  3. `~/Library/Developer/Xcode/Archives`
  4. `~/Library/Developer/Xcode/DocumentationCache`
- All paths use `~` prefix (not absolute home).

## Side Effects

- None (pure discovery).

## Errors

- None.

## Exit Code

- Test passes when expectations match.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.XcodeLoc == nil {
		t.Fatal("missing xcode location")
	}
	loc := resp.XcodeLoc
	if loc.Path != "~/Library/Developer/Xcode/DerivedData" {
		t.Fatalf("expected primary DerivedData path, got %s", loc.Path)
	}
	expected := []string{
		"~/Library/Developer/CoreSimulator/Devices",
		"~/Library/Developer/Xcode/iOS DeviceSupport",
		"~/Library/Developer/Xcode/Archives",
		"~/Library/Developer/Xcode/DocumentationCache",
	}
	if len(loc.ExtraPaths) != len(expected) {
		t.Fatalf("expected %d ExtraPaths, got %d: %v", len(expected), len(loc.ExtraPaths), loc.ExtraPaths)
	}
	for i, want := range expected {
		if loc.ExtraPaths[i] != want {
			t.Fatalf("ExtraPaths[%d]: expected %q, got %q", i, want, loc.ExtraPaths[i])
		}
	}
}
```