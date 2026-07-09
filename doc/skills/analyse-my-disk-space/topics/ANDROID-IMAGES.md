# Android images and AVDs

Reference for inspecting and reclaiming Android emulator (AVD) disk usage discovered during a home scan. Complements the main [SKILL.md](../SKILL.md) playbook.

## Layout on disk

| Piece | Typical path | Role |
|-------|----------------|------|
| AVD registry | `~/.android/avd/<Name>.ini` | Points tools at the data directory |
| AVD data dir | `~/.android/avd/<Name>.avd/` | VM disks, snapshots, runtime state |
| System image | `~/Library/Android/sdk/system-images/android-<API>/…` | Shared OS image (not inside the AVD) |
| Broader SDK | `~/Library/Android/sdk/` | platforms, build-tools, emulator binary |

An AVD is registered by a small `.ini` and a sibling `.avd` directory. Example:

```text
~/.android/avd/MediumSize_Phone_SDK35.ini
~/.android/avd/Medium-Size_Phone_SDK_35.avd/
```

Registry `.ini` (short):

```ini
avd.ini.encoding=UTF-8
path=/Users/<user>/.android/avd/Medium-Size_Phone_SDK_35.avd
path.rel=avd/Medium-Size_Phone_SDK_35.avd
target=android-35
```

Hardware and image selection live in `<name>.avd/config.ini` (API level, ABI, RAM, display, `image.sysdir.1`, sdcard size, etc.).

## What usually consumes space

Prefer **offline** analysis of the home capture from [SKILL.md](../SKILL.md):

```bash
disk-usage-analyser scan ~ --json --max-depth 6 > /tmp/disk-scan.json
disk-usage-analyser inspect /tmp/disk-scan.json --find .android
disk-usage-analyser inspect /tmp/disk-scan.json --at ~/.android
disk-usage-analyser inspect /tmp/disk-scan.json --suffix .qcow2 --min-size 50M
disk-usage-analyser inspect /tmp/disk-scan.json --find Android/sdk
```

Only if that tree is missing or truncated, run **one** targeted scan (do not chain multiple shallow scans — each `scan` walks fully; `--max-depth` only limits tree emission):

```bash
disk-usage-analyser scan ~/.android --json --max-depth 6 > /tmp/android-scan.json
disk-usage-analyser inspect /tmp/android-scan.json --top 30
```

Typical large files **inside** one AVD directory:

| Size class | File / path | Meaning |
|------------|-------------|---------|
| Often largest | `userdata-qemu.img.qcow2` | Installed apps + app data (grows with use; base image may be only ~10 MB) |
| Large | `snapshots/**/ram.bin` | Quick-boot snapshot (mostly RAM dump) |
| Fixed by config | `sdcard.img` | Virtual SD card (`sdcard.size` in config.ini) |
| Smaller | `cache.img*`, `encryptionkey.img*`, `initrd` | Boot / runtime helpers |

**Example from a real medium-phone API 35 Play AVD (~17.6 GB total):**

| Size | Path |
|------|------|
| ~12.3 GB | `userdata-qemu.img.qcow2` |
| ~3.1 GB | `snapshots/default_boot/ram.bin` |
| ~2.0 GB | `sdcard.img` |
| ~170 MB | cache / encryptionkey / initrd / misc |

Config highlights from that device: API **35**, Google Play, **arm64-v8a**, 1080×2400 @ 420 dpi, **6 GB** RAM, data partition **20G**, sdcard **2G**. System image path:

`system-images/android-35/google_apis_playstore/arm64-v8a/`

The system image lives under the **SDK**, not under `~/.android/avd`. Deleting an AVD does not remove SDK system-images; deleting system-images forces a re-download when creating a new AVD.

## Relevant file suffixes

Treat these as Android emulator / disk-image artifacts when ranking findings:

| Suffix / pattern | Role |
|------------------|------|
| `.avd/` (directory) | AVD data root |
| AVD `*.ini` under `~/.android/avd/` | AVD registry |
| `.img` | Raw disk (userdata base, sdcard, cache, encryptionkey) |
| `.qcow2` | QEMU overlay (often the grown userdata) |
| `userdata-qemu.img`, `userdata-qemu.img.qcow2` | User partition |
| `sdcard.img` | Emulated SD card |
| `cache.img`, `cache.img.qcow2` | Emulator cache partition |
| `encryptionkey.img*` | Encryption key store for the VM |
| `ram.bin` | Snapshot RAM image under `snapshots/` |
| `hardware-qemu.ini` | Generated QEMU hardware config |

## How to inspect

### CLI

```bash
# List AVDs
~/Library/Android/sdk/cmdline-tools/latest/bin/avdmanager list avd
~/Library/Android/sdk/emulator/emulator -list-avds

# Config
cat ~/.android/avd/*.ini
cat ~/.android/avd/<Name>.avd/config.ini

# Size — prefer offline from home JSON; else one targeted scan
disk-usage-analyser scan ~/.android --json --max-depth 6
# or classic:
du -sh ~/.android/avd/<Name>.avd/*

# Optional image metadata (if qemu-img installed)
qemu-img info ~/.android/avd/<Name>.avd/userdata-qemu.img.qcow2
```

### Live (emulator running)

```bash
emulator -avd <AvdId>
adb shell df -h
adb shell du -sh /data /sdcard 2>/dev/null
adb shell pm list packages
```

### Android Studio

**Device Manager** → ⋮ on the AVD → **Show on Disk** / details — opens the same folder and shows API / image / hardware.

## Wipe vs delete vs archive

| Goal | Action | Keeps | Frees (typical) |
|------|--------|--------|-----------------|
| Free most space, keep device definition | Wipe Data / cold boot (Device Manager) | AVD entry + system image | Userdata + snapshots (often 10s of GB) |
| Free only quick-boot | Remove `snapshots/` (emulator stopped) | Apps/data until wipe | ~snapshot size (e.g. 3 GB) |
| Full remove | Delete AVD in Device Manager, or remove `.avd` + matching `.ini` | Nothing of that VM | Full AVD tree |
| Might need data later | `tar` the `.ini` + `.avd` first | Everything if restored | Nothing until archive deleted |

**Partial cleanup (emulator stopped):**

```bash
# Drop quick-boot snapshot only
rm -rf ~/.android/avd/<Name>.avd/snapshots

# Prefer "Wipe Data" in Device Manager for a clean userdata without losing the AVD slot
```

## If you delete the AVD — recovery

### Instance state is not recoverable without a backup

Deleting the AVD removes installed apps, app data, sdcard contents, and snapshots. There is **no cloud restore**. Recovery of **that** VM requires:

- A manual archive you made, or
- Time Machine / other backup of `~/.android/avd/`

### Recreate an equivalent empty AVD (easy)

As long as the **system image** remains in the SDK, recreate in minutes:

**Android Studio:** Device Manager → Create Device → same form factor → same API / Google Play image.

**CLI (example package):**

```bash
sdkmanager --list_installed | grep 'system-images;android-35'

avdmanager create avd \
  -n MediumSize_Phone_SDK35 \
  -k "system-images;android-35;google_apis_playstore;arm64-v8a" \
  -d medium_phone
```

Then match RAM/sdcard via Device Manager or `config.ini` if needed.

### Backup / restore

```bash
# Backup
tar -czf ~/Desktop/avd-<Name>-backup.tgz \
  -C ~/.android/avd \
  <Registry>.ini \
  <Name>.avd

# Restore
tar -xzf ~/Desktop/avd-<Name>-backup.tgz -C ~/.android/avd
avdmanager list avd
```

## Recommendations for disk-analysis reports

1. Separate **AVD userdata/snapshots** (reclaimable per device) from **SDK system-images** (shared; re-download cost).
2. Prefer **Wipe Data** or **delete snapshots** over full AVD delete when the user still needs an API N emulator occasionally.
3. Only recommend full AVD delete when the emulator is unused; mention that recreate needs the system image package still installed (or a network download).
4. Before destructive steps, offer a short `tar` backup of the `.ini` + `.avd` if userdata might matter.

## Example scan commands

```bash
# Preferred: already covered by home capture
disk-usage-analyser scan ~ --json --max-depth 6

# Targeted only when needed
disk-usage-analyser scan ~/.android --json --max-depth 6
disk-usage-analyser scan ~/Library/Android/sdk --json --max-depth 6
```
