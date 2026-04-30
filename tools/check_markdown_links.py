#!/usr/bin/env python3

from __future__ import annotations

import re
import subprocess
import sys
from pathlib import Path


LINK_RE = re.compile(r"!\[[^\]]*\]\(([^)]+)\)|(?<!!)\[[^\]]+\]\(([^)]+)\)")
EXTERNAL_PREFIXES = ("http://", "https://", "mailto:", "tel:", "javascript:", "data:")


def normalize_target(raw: str) -> str:
    target = raw.strip()
    if target.startswith("<") and target.endswith(">"):
        target = target[1:-1].strip()
    target = target.split("#", 1)[0].split("?", 1)[0].strip()
    return target


def should_skip(target: str) -> bool:
    return not target or target.startswith("#") or target.startswith("/") or target.startswith(EXTERNAL_PREFIXES)


def iter_markdown_files(root: Path) -> list[Path]:
    return sorted(p for p in root.rglob("*.md") if p.is_file())


def check_files(files: list[Path]) -> list[tuple[Path, int, str]]:
    broken: list[tuple[Path, int, str]] = []
    for path in files:
        text = path.read_text(encoding="utf-8")
        for lineno, line in enumerate(text.splitlines(), 1):
            for match in LINK_RE.finditer(line):
                raw = match.group(1) or match.group(2) or ""
                target = normalize_target(raw)
                if should_skip(target):
                    continue
                resolved = (path.parent / target).resolve()
                if not resolved.exists():
                    broken.append((path, lineno, target))
    return broken


def print_report(title: str, broken: list[tuple[Path, int, str]], repo_root: Path) -> None:
    print(f"## {title}")
    if not broken:
        print("OK")
        return
    for path, lineno, target in broken:
        rel = path.relative_to(repo_root)
        print(f"{rel}:{lineno}: {target}")


def main() -> int:
    args = set(sys.argv[1:])
    unknown = args - {"--skip-prepare"}
    if unknown:
        print(f"unknown arguments: {sorted(unknown)}", file=sys.stderr)
        return 2
    skip_prepare = "--skip-prepare" in args

    repo_root = Path(__file__).resolve().parents[1]

    repo_docs_files = [
        repo_root / "README.zh-CN.md",
        repo_root / "ReadMe.md",
        repo_root / "getting-started.md",
    ]
    repo_docs_files.extend(iter_markdown_files(repo_root / "docs"))
    repo_broken = check_files(repo_docs_files)

    site_root = repo_root / "site-docs"
    if not skip_prepare:
        subprocess.run([sys.executable, str(repo_root / "tools" / "prepare_docs_site.py")], check=True, cwd=repo_root)
    elif not site_root.exists():
        print("site-docs is missing; run without --skip-prepare first", file=sys.stderr)
        return 2
    site_broken = check_files(iter_markdown_files(site_root))

    print_report("Repository Docs", repo_broken, repo_root)
    print()
    print_report("Generated Site Docs", site_broken, repo_root)

    if repo_broken or site_broken:
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
