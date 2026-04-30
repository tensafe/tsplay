#!/bin/sh
set -eu

usage() {
  cat <<'EOF'
Install the tsplay-flow-authoring skill into the local Codex skills directory.

Usage:
  sh ./install.sh [--codex-home PATH] [--install-dir PATH] [--tsplay PATH] [--extract-assets PATH] [--force]

Options:
  --codex-home PATH    Override CODEX_HOME. Default: $CODEX_HOME or $HOME/.codex
  --install-dir PATH   Override the skills root directory. Default: <codex-home>/skills
  --tsplay PATH        Explicit tsplay binary path used with --extract-assets
  --extract-assets DIR After install, run tsplay -action extract-assets -extract-root DIR
  --force              Replace an existing install without creating a backup
  -h, --help           Show this help message
EOF
}

log() {
  printf '%s\n' "$*"
}

die() {
  printf 'Error: %s\n' "$*" >&2
  exit 1
}

resolve_dir() {
  dir=$1
  if [ ! -d "$dir" ]; then
    die "directory does not exist: $dir"
  fi
  (
    cd "$dir"
    pwd -P
  )
}

find_tsplay() {
  if [ -n "${TSPLAY_BIN:-}" ]; then
    if [ ! -x "$TSPLAY_BIN" ]; then
      die "tsplay binary is not executable: $TSPLAY_BIN"
    fi
    printf '%s\n' "$TSPLAY_BIN"
    return 0
  fi
  if command -v tsplay >/dev/null 2>&1; then
    command -v tsplay
    return 0
  fi
  return 1
}

SCRIPT_DIR=$(resolve_dir "$(dirname "$0")")
SKILL_NAME=$(basename "$SCRIPT_DIR")

if [ ! -f "$SCRIPT_DIR/SKILL.md" ]; then
  die "SKILL.md not found next to install.sh"
fi

CODEX_HOME_VALUE=${CODEX_HOME:-}
INSTALL_ROOT=""
TSPLAY_BIN=""
EXTRACT_ROOT=""
FORCE=0

while [ $# -gt 0 ]; do
  case "$1" in
    --codex-home)
      [ $# -ge 2 ] || die "--codex-home requires a path"
      CODEX_HOME_VALUE=$2
      shift 2
      ;;
    --install-dir)
      [ $# -ge 2 ] || die "--install-dir requires a path"
      INSTALL_ROOT=$2
      shift 2
      ;;
    --tsplay)
      [ $# -ge 2 ] || die "--tsplay requires a path"
      TSPLAY_BIN=$2
      shift 2
      ;;
    --extract-assets)
      [ $# -ge 2 ] || die "--extract-assets requires a target directory"
      EXTRACT_ROOT=$2
      shift 2
      ;;
    --force)
      FORCE=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      die "unknown argument: $1"
      ;;
  esac
done

if [ -z "$INSTALL_ROOT" ]; then
  if [ -n "$CODEX_HOME_VALUE" ]; then
    INSTALL_ROOT=$CODEX_HOME_VALUE/skills
  else
    INSTALL_ROOT=$HOME/.codex/skills
  fi
fi

mkdir -p "$INSTALL_ROOT"
INSTALL_ROOT=$(resolve_dir "$INSTALL_ROOT")
TARGET_DIR=$INSTALL_ROOT/$SKILL_NAME

if [ -d "$TARGET_DIR" ]; then
  TARGET_REAL=$(resolve_dir "$TARGET_DIR")
  if [ "$TARGET_REAL" = "$SCRIPT_DIR" ]; then
    log "Skill is already installed at $TARGET_DIR"
  else
    if [ "$FORCE" -eq 1 ]; then
      rm -rf "$TARGET_DIR"
      log "Removed existing install at $TARGET_DIR"
    else
      BACKUP_DIR="${TARGET_DIR}.backup.$(date +%Y%m%d%H%M%S)"
      mv "$TARGET_DIR" "$BACKUP_DIR"
      log "Backed up existing install to $BACKUP_DIR"
    fi
    mkdir -p "$TARGET_DIR"
    cp -R "$SCRIPT_DIR/." "$TARGET_DIR/"
    log "Installed $SKILL_NAME to $TARGET_DIR"
  fi
else
  mkdir -p "$TARGET_DIR"
  cp -R "$SCRIPT_DIR/." "$TARGET_DIR/"
  log "Installed $SKILL_NAME to $TARGET_DIR"
fi

if [ -n "$EXTRACT_ROOT" ]; then
  if TSPLAY_PATH=$(find_tsplay); then
    "$TSPLAY_PATH" -action extract-assets -extract-root "$EXTRACT_ROOT"
    log "Extracted bundled assets to $EXTRACT_ROOT"
  else
    die "tsplay was not found on PATH; rerun with --tsplay /path/to/tsplay or install tsplay first"
  fi
fi

if TSPLAY_PATH=$(find_tsplay); then
  log "Detected tsplay: $TSPLAY_PATH"
  log "Optional next step: $TSPLAY_PATH -action list-assets"
else
  log "tsplay was not found on PATH."
  log "Install a matching tsplay release binary, or use ./tsplay / go run . when working inside the repo."
fi

log "Skill install complete."
