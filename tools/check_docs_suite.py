#!/usr/bin/env python3

from __future__ import annotations

import subprocess
import sys
from pathlib import Path


def run_check(repo_root: Path, script_name: str) -> int:
    return subprocess.run(
        [sys.executable, str(repo_root / "tools" / script_name), "--skip-prepare"],
        cwd=repo_root,
    ).returncode


def main() -> int:
    repo_root = Path(__file__).resolve().parents[1]

    prepare = subprocess.run(
        [sys.executable, str(repo_root / "tools" / "prepare_docs_site.py")],
        cwd=repo_root,
    )
    if prepare.returncode != 0:
        return prepare.returncode

    status = 0
    for script_name in ("check_markdown_links.py", "check_doc_continuity.py"):
        result = run_check(repo_root, script_name)
        if result != 0:
            status = result
    return status


if __name__ == "__main__":
    raise SystemExit(main())
