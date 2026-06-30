# Scenario

**Feature**: sort repository scan rows by size DESC

```
# repositoryScansLayout sort helpers reorder repos and children largest-first
repositoryScansLayout.sort* -> DESC by size bytes
```

## Preconditions

- Sort applies after filter on every state update.

## Steps

1. Leaf sets `req.Op` for the sort function under test.

## Context

- See leaf ASSERT.md for expected order.

```go
func Setup(t *testing.T, req *Request) error {
	if req.FixtureFile == "" {
		req.FixtureFile = "testdata/fixture.json"
	}
	return nil
}
```