#!/usr/bin/env python3

from __future__ import annotations

import shutil
import time
from pathlib import Path


def copy_path(src: Path, dst: Path) -> None:
    if src.is_dir():
        shutil.copytree(src, dst, dirs_exist_ok=True)
    else:
        dst.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(src, dst)


def reset_site_root(site_root: Path) -> None:
    if not site_root.exists():
        return
    last_error: OSError | None = None
    for _ in range(3):
        try:
            shutil.rmtree(site_root)
            return
        except OSError as err:
            last_error = err
            time.sleep(0.1)
    if last_error is not None:
        raise last_error


def main() -> None:
    repo_root = Path(__file__).resolve().parents[1]
    site_root = repo_root / "site-docs"

    reset_site_root(site_root)
    site_root.mkdir(parents=True)

    # Keep the rendered site source close to the repository layout so the
    # existing tutorial links to README, script, and demo assets still work.
    copy_path(repo_root / "site-src", site_root)
    copy_path(repo_root / "README.zh-CN.md", site_root / "README.zh-CN.md")
    copy_path(repo_root / "ReadMe.md", site_root / "ReadMe.md")
    copy_path(repo_root / "embedded_assets.go", site_root / "embedded_assets.go")
    copy_path(repo_root / "static_server.go", site_root / "static_server.go")
    copy_path(repo_root / "docs", site_root / "docs")
    copy_path(repo_root / "script", site_root / "script")
    copy_path(repo_root / "demo", site_root / "demo")

    docs_index = site_root / "docs" / "index.md"
    if docs_index.exists():
        docs_index.unlink()


if __name__ == "__main__":
    main()
