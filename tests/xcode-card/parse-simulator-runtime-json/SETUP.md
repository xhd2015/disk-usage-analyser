# Scenario

**Feature**: parse simulator runtime list JSON into TmpRuntimeItem slice

```
# ParseSimulatorRuntimeJSON reads UUID-keyed simctl JSON
xcrun simctl runtime list -j -> ParseSimulatorRuntimeJSON -> []TmpRuntimeItem
```

## Preconditions

- Fixture contains two runtimes: Ready+deletable and Mounting+non-deletable.
- Size lookup returns known bytes per mountPath.

## Steps

1. Set `req.Op` to `parse-simulator-runtime-json`.
2. Load fixture and supply mount-path sizes.

## Context

- `Type` uses `version` field; `ActiveCount` is 1 when `state == "Ready"`.
- `Reclaimable` equals size when `deletable == true`, else 0.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "parse-simulator-runtime-json"
	req.FixtureFile = "testdata/simulator-runtimes.json"
	req.MountSizes = map[string]int64{
		"/Library/Developer/CoreSimulator/Volumes/iOS_22F77":  90000000000,
		"/Library/Developer/CoreSimulator/Volumes/iOS_21E213": 45000000000,
	}
	return nil
}
```