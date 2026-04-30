#!/bin/sh

set -eu

REPO="tensafe/tsplay"
BASE_URL="https://github.com/$REPO/releases/latest/download"
INSTALL_DIR="."
BINARY_NAME="tsplay"
SKIP_RUN=0

usage() {
  cat <<'EOF'
Usage:
  sh quickstart-install.sh [--install-dir PATH] [--binary-name NAME] [--skip-run]

What it does:
  - Detects macOS or Linux automatically
  - Downloads the matching TSPlay release binary
  - Makes it executable
  - Runs: ./tsplay -action quickstart-demo

Examples:
  sh quickstart-install.sh
  sh quickstart-install.sh --install-dir ./bin
  sh quickstart-install.sh --skip-run
EOF
}

die() {
  printf 'error: %s\n' "$*" >&2
  exit 1
}

log() {
  printf '%s\n' "$*"
}

download_file() {
  url=$1
  dest=$2
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$dest"
    return
  fi
  if command -v wget >/dev/null 2>&1; then
    wget -qO "$dest" "$url"
    return
  fi
  die "curl or wget is required"
}

detect_os() {
  case "$(uname -s)" in
    Darwin) printf 'darwin' ;;
    Linux) printf 'linux' ;;
    *)
      die "unsupported operating system: $(uname -s). Use the Windows PowerShell installer on Windows."
      ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) printf 'amd64' ;;
    arm64|aarch64) printf 'arm64' ;;
    *) die "unsupported architecture: $(uname -m)" ;;
  esac
}

while [ $# -gt 0 ]; do
  case "$1" in
    --install-dir)
      [ $# -ge 2 ] || die "--install-dir requires a value"
      INSTALL_DIR=$2
      shift 2
      ;;
    --binary-name)
      [ $# -ge 2 ] || die "--binary-name requires a value"
      BINARY_NAME=$2
      shift 2
      ;;
    --skip-run)
      SKIP_RUN=1
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

OS_TAG=$(detect_os)
ARCH_TAG=$(detect_arch)
DOWNLOAD_URL="$BASE_URL/tsplay-${OS_TAG}-${ARCH_TAG}"

mkdir -p "$INSTALL_DIR"
INSTALL_DIR_ABS=$(cd "$INSTALL_DIR" && pwd)
TARGET_PATH="$INSTALL_DIR_ABS/$BINARY_NAME"
TMP_FILE=$(mktemp "${TMPDIR:-/tmp}/tsplay-quickstart.XXXXXX")
trap 'rm -f "$TMP_FILE"' EXIT INT TERM

log "Downloading TSPlay for ${OS_TAG}/${ARCH_TAG}"
download_file "$DOWNLOAD_URL" "$TMP_FILE"
chmod +x "$TMP_FILE"
mv "$TMP_FILE" "$TARGET_PATH"
trap - EXIT INT TERM

log "Installed: $TARGET_PATH"

if [ "$SKIP_RUN" -eq 1 ]; then
  log "Skipped quickstart run."
  exit 0
fi

log "Running quickstart demo..."
"$TARGET_PATH" -action quickstart-demo

log ""
log "Next steps:"
log "  $TARGET_PATH -action file-srv -addr :8000"
log "  $TARGET_PATH -flow script/tutorials/10_assert_page_state.flow.yaml"
