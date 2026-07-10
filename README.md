# disk-usage-analyser

A tool similar to the `du` command and macOS storage analyser. It helps users identify large files and eliminate them to save storage space.

## Agent skill

Use this skill in codex/claude code/grok:

```sh
# install first
disk-usage-analyser skill --install --global    # default: .agents/skills/analyse-my-disk-space/
```

Use:
```text
/analyse-my-disk-space help me figure out disk space eating monster
```

Other commands related to skill:

```sh
disk-usage-analyser skill --show             # print skill body
disk-usage-analyser skill --show android-images
```

# Usage
```sh
git clone https://github.com/xhd2015/disk-usage-analyser

cd disk-usage-analyser

# start dev mode
go run ./script/dev
```
