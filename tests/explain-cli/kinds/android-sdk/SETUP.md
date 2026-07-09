# Scenario

**Feature**: Android SDK reclaim kind (`android-sdk`) for `~/Library/Android/sdk`

```
# Auto-detect: path ends with Android/sdk (or basename sdk under Android), OR
# dir contains platform-tools + one of platforms/system-images/build-tools/cmdline-tools/emulator
explain "…/Library/Android/sdk" -> kind=android-sdk, ContextRoot=SDK root, confidence high
explain "…/Library/Android/sdk/platform-tools/f" -> kind=android-sdk (ContextRoot=SDK root)

# Force via kind id (PATH optional)
explain --kind android-sdk [SCOPE?]
  scope = PATH | CLIOptions.HomeDir | os.UserHomeDir()
  measure {scope}/Library/Android/sdk when scope is home-like;
  if scope is already an SDK root (signatures), use it

# Roles (top-level children): system-images ☑, sources ☑, skins ☑, tmp ☑,
# emulator ☐, build-tools ☐, platforms ☐, platform-tools ☐, cmdline-tools ☐, …
# SAFE TO RECLAIM: temp usually-safe; system-images/sources/skins usually-safe-with-caution;
#   keep platform-tools/build-tools/platforms/emulator bulk
# HOW TO PURGE: sdkmanager --list_installed / --uninstall; scan; never rm -rf;
#   UI (Android Studio SDK settings) only in Notes
```

## Preconditions

- Fixture from `writeAndroidSDKFixture`: `{parent}/Library/Android/sdk` with system-images,
  emulator, sources, build-tools, platform-tools, platforms, `.temp`.
- Content payload sum is `androidSDKContentBytes` (890).
- Detection runs before `generic-dir` / `generic-file`.
- Tests must never use the real user home for fixtures; use `req.FixtureDir` / `req.HomeDir`.

## Context

- Kind id in output/JSON is **`android-sdk`** (same as CLI force id).
- Breakdown assigns roles (system-images, emulator, sources, build-tools, platform-tools,
  platforms, tmp, …).
- SAFE TO RECLAIM must not treat platform-tools / adb / build-tools bulk as usually-safe-only.
- HOW TO PURGE is CLI-first (`sdkmanager`, `disk-usage-analyser scan`); never `rm -rf`.


```go
func Setup(t *testing.T, req *Request) error {
	// Mark mode for android-sdk leaves; concrete TargetPath / --kind is set by child leaves.
	req.Mode = "cli"
	return nil
}
```
