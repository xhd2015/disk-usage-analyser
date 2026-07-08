# Scenario

**Feature**: Xcode DiscoverLocations has four ordered ExtraPaths

```
# DiscoverLocations registers Xcode with DerivedData primary + four extras
DiscoverLocations -> TmpLocation(category=xcode) -> ExtraPaths[0..3]
```

## Preconditions

- Primary path unchanged: `~/Library/Developer/Xcode/DerivedData`.

## Steps

1. Set `req.Op` to `discover-extra-paths`.
2. Call `DiscoverLocations` with test home dir.

## Context

- ExtraPaths order is part of the contract (breakdown row order).

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "discover-extra-paths"
	return nil
}
```