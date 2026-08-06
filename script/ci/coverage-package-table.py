#!/usr/bin/env python3
"""Print a markdown package coverage table from a classic go coverprofile.

Usage:
  python3 script/ci/coverage-package-table.py coverage.out

Skips script/ and cmd/ paths under disk-usage-analyser/.
Rows sorted by coverage ascending, then package name.
"""

from __future__ import annotations

import sys
from collections import defaultdict
from pathlib import Path

# Module path from go.mod
MODULE = "disk-usage-analyser"
PREFIX = MODULE + "/"


def main() -> int:
    if len(sys.argv) != 2:
        print("usage: coverage-package-table.py <coverage.out>", file=sys.stderr)
        return 2
    path = Path(sys.argv[1])
    if not path.is_file():
        print(f"missing profile: {path}", file=sys.stderr)
        return 1

    pkg: dict[str, list[int]] = defaultdict(lambda: [0, 0])  # cov, tot stmts
    for line in path.read_text().splitlines():
        if line.startswith("mode:") or not line.strip():
            continue
        left, count_s = line.rsplit(" ", 1)
        key, num_s = left.rsplit(" ", 1)
        count, num = int(count_s), int(num_s)
        file = key.rsplit(":", 1)[0]
        if PREFIX not in file and not file.startswith(MODULE):
            continue
        if file.startswith(PREFIX):
            rel = file[len(PREFIX) :]
        elif file == MODULE or file.endswith("/" + MODULE):
            rel = ""
        elif PREFIX in file:
            rel = file.split(PREFIX, 1)[1]
        else:
            continue
        if rel.startswith("script/") or rel.startswith("cmd/"):
            continue
        p = "/".join(rel.split("/")[:-1]) if "/" in rel else (rel or ".")
        pkg[p][1] += num
        if count > 0:
            pkg[p][0] += num

    rows = sorted(
        ((100.0 * c / t if t else 0.0, p) for p, (c, t) in pkg.items() if t > 0),
        key=lambda r: (r[0], r[1]),
    )
    print("| Coverage | Package |")
    print("|----------|---------|")
    for pct, p in rows:
        print(f"| {pct:.1f}% | `{p}` |")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
