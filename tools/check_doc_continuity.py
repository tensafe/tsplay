#!/usr/bin/env python3

from __future__ import annotations

import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path


@dataclass(frozen=True)
class ContinuityRule:
    path: str
    title: str
    required_snippets: tuple[str, ...]
    required_any_groups: tuple[tuple[str, ...], ...] = ()


REPO_RULES: tuple[ContinuityRule, ...] = (
    ContinuityRule(
        path="ReadMe.md",
        title="Root English README keeps major entry paths visible",
        required_snippets=(
            "docs/tutorials/README.md",
            "docs/training/ai-intent-to-flow.md",
            "docs/training/README.md",
            "5-Minute Quick Start",
        ),
    ),
    ContinuityRule(
        path="README.zh-CN.md",
        title="Root Chinese README keeps major entry paths visible",
        required_snippets=(
            "docs/tutorials/README.zh-CN.md",
            "docs/training/ai-intent-to-flow.md",
            "docs/training/README.md",
            "5 分钟跑起来",
        ),
    ),
    ContinuityRule(
        path="getting-started.md",
        title="Root quick-start alias keeps next-step continuity",
        required_snippets=(
            "docs/tutorials/README.zh-CN.md",
            "docs/tutorials/track-newbie.zh-CN.md",
            "docs/training/ai-intent-to-flow.md",
            "docs/tutorials/144-single-binary-delivery-flow.md",
            "docs/tutorials/87-build-handoff-artifact-manifest.md",
        ),
    ),
    ContinuityRule(
        path="site-src/getting-started.md",
        title="Docs-site quick start keeps default branches and common gaps",
        required_snippets=(
            "docs/tutorials/119-mcp-chain-overview.md",
            "docs/tutorials/142-list-assets-for-beginners.md",
            "docs/product/core-feature-roadmap.md",
        ),
        required_any_groups=(
            ("跑通后下一步做什么", "下一步"),
            ("常见断层提醒", "几个常见卡点", "常见卡点", "常见问题"),
        ),
    ),
    ContinuityRule(
        path="docs/tutorials/README.zh-CN.md",
        title="Tutorial hub still points to first-run and stable tracks",
        required_snippets=(
            "../../getting-started.md",
            "track-newbie.zh-CN.md",
            "Lesson 111",
            "30 分钟起步线",
        ),
    ),
    ContinuityRule(
        path="docs/tutorials/README.md",
        title="English tutorial hub still points to first-run and stable tracks",
        required_snippets=(
            "../../getting-started.md",
            "track-newbie.md",
            "Lesson 111",
            "30-Minute Starter Path",
        ),
    ),
    ContinuityRule(
        path="docs/README.md",
        title="Docs map exposes roadmap, execution board, and doc health",
        required_snippets=(
            "product/core-feature-roadmap.md",
            "product/core-feature-execution-board.md",
            "doc-health-audit.md",
        ),
    ),
)


SITE_RULES: tuple[ContinuityRule, ...] = (
    ContinuityRule(
        path="ReadMe.md",
        title="Generated English README keeps major entry paths visible",
        required_snippets=(
            "docs/tutorials/README.md",
            "docs/training/ai-intent-to-flow.md",
            "docs/training/README.md",
            "5-Minute Quick Start",
        ),
    ),
    ContinuityRule(
        path="getting-started.md",
        title="Generated quick-start page keeps default branches and common gaps",
        required_snippets=(
            "docs/tutorials/119-mcp-chain-overview.md",
            "docs/tutorials/142-list-assets-for-beginners.md",
        ),
        required_any_groups=(
            ("跑通后下一步做什么", "下一步"),
            ("常见断层提醒", "几个常见卡点", "常见卡点", "常见问题"),
        ),
    ),
    ContinuityRule(
        path="docs/tutorials/README.md",
        title="Generated English tutorial hub still points to first-run and stable tracks",
        required_snippets=(
            "../../getting-started.md",
            "track-newbie.md",
            "Lesson 111",
            "30-Minute Starter Path",
        ),
    ),
    ContinuityRule(
        path="docs/README.md",
        title="Generated docs map exposes roadmap, execution board, and doc health",
        required_snippets=(
            "product/core-feature-roadmap.md",
            "product/core-feature-execution-board.md",
            "doc-health-audit.md",
        ),
    ),
)


def check_rules(root: Path, rules: tuple[ContinuityRule, ...]) -> list[str]:
    errors: list[str] = []
    for rule in rules:
        path = root / rule.path
        if not path.exists():
            errors.append(f"{rule.path}: missing file for rule: {rule.title}")
            continue
        text = path.read_text(encoding="utf-8")
        for snippet in rule.required_snippets:
            if snippet not in text:
                errors.append(f"{rule.path}: missing snippet {snippet!r} for rule: {rule.title}")
        for options in rule.required_any_groups:
            if not any(option in text for option in options):
                errors.append(
                    f"{rule.path}: missing any snippet from {options!r} for rule: {rule.title}"
                )
    return errors


def iter_markdown_files(root: Path) -> list[Path]:
    files: list[Path] = []
    if root.name == "site-docs":
        files.extend(root.rglob("*.md"))
        return sorted(set(files))

    for name in ("README.md", "README.zh-CN.md", "ReadMe.md", "getting-started.md"):
        path = root / name
        if path.exists():
            files.append(path)
    for folder in ("docs", "site-src"):
        path = root / folder
        if path.exists():
            files.extend(path.rglob("*.md"))
    return sorted(set(files))


def check_details_markdown(root: Path) -> list[str]:
    errors: list[str] = []
    for path in iter_markdown_files(root):
        text = path.read_text(encoding="utf-8")
        for lineno, line in enumerate(text.splitlines(), start=1):
            if "<details" in line and "markdown=" not in line:
                rel = path.relative_to(root)
                errors.append(
                    f'{rel}:{lineno}: <details> is missing markdown="1"; nested Markdown may render as plain text in MkDocs'
                )
    return errors


def print_report(title: str, errors: list[str]) -> None:
    print(f"## {title}")
    if not errors:
        print("OK")
        return
    for item in errors:
        print(item)


def main() -> int:
    args = set(sys.argv[1:])
    unknown = args - {"--skip-prepare"}
    if unknown:
        print(f"unknown arguments: {sorted(unknown)}", file=sys.stderr)
        return 2
    skip_prepare = "--skip-prepare" in args

    repo_root = Path(__file__).resolve().parents[1]

    repo_errors = check_rules(repo_root, REPO_RULES)
    repo_errors.extend(check_details_markdown(repo_root))

    site_root = repo_root / "site-docs"
    if not skip_prepare:
        subprocess.run([sys.executable, str(repo_root / "tools" / "prepare_docs_site.py")], check=True, cwd=repo_root)
    elif not site_root.exists():
        print("site-docs is missing; run without --skip-prepare first", file=sys.stderr)
        return 2
    site_errors = check_rules(site_root, SITE_RULES)
    site_errors.extend(check_details_markdown(site_root))

    print_report("Repository Continuity", repo_errors)
    print()
    print_report("Generated Site Continuity", site_errors)

    if repo_errors or site_errors:
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
