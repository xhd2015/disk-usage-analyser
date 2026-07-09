# Scenario

**Feature**: Android AVD reclaim kind (`android-avd`)

```
# Signals: dir named *.avd or config.ini + userdata-qemu.img* / sdcard.img / snapshots
explain MediumPhone.avd/ -> kind=android-avd
explain MediumPhone.avd/userdata-qemu.img.qcow2 -> kind=android-avd (parent context)
```

## Preconditions

- Fixture layout from `writeAVDFixture`: config.ini, userdata-qemu.img.qcow2, sdcard.img, snapshots/…/ram.bin.
- Content payload sum is `avdContentBytes` (732).

## Context

- Breakdown should assign roles (e.g. user-data, sdcard, snapshot) when practical.
- SAFE TO RECLAIM may mention snapshots as usually safer than wiping userdata.
- HOW TO PURGE official commands are CLI-first (`emulator`, `avdmanager`); UI only in Notes.


```go
func Setup(t *testing.T, req *Request) error {
	// Mark mode for AVD human-explain leaves; concrete TargetPath is set by dir/ or file/.
	req.Mode = "cli"
	return nil
}
```

