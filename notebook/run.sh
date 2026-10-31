#!/bin/sh
# Starts the notebook binary, first installing the release that matches this
# source. If that release isn't up yet, the installed binary keeps working.
# With --install it only installs.
set -u

here=$(cd "$(dirname "$0")" && pwd)
bin=${XDG_DATA_HOME:-$HOME/.local/share}/workbench/bin
binary=$bin/notebook
stamp=$bin/notebook.source
releases=${NOTEBOOK_RELEASE_URL:-https://github.com/tahabakhit/workbench/releases/download}

sha() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1"
  else
    shasum -a 256 "$1"
  fi
}

platform() {
  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  case $(uname -m) in
    x86_64 | amd64) echo "$os-amd64" ;;
    aarch64 | arm64) echo "$os-arm64" ;;
    *) return 1 ;;
  esac
}

fetch() {
  curl -fsL --connect-timeout 5 --max-time 120 -o "$2" "$1"
}

# Downloads into a private folder and renames over the binary, so a concurrent
# run sees the old binary or the new one, never part of either.
install_release() {
  file=notebook-$(platform) || return 1
  mkdir -p "$bin" || return 1
  work=$(mktemp -d "$bin/.install-XXXXXX") || return 1
  trap 'rm -rf "$work"' EXIT
  url=$releases/notebook-$1
  fetch "$url/checksums.txt" "$work/checksums.txt" || return 1
  fetch "$url/$file" "$work/$file" || return 1
  want=$(awk -v f="$file" '$2 == f { print $1 }' "$work/checksums.txt")
  got=$(sha "$work/$file" | cut -d' ' -f1)
  if [ -z "$want" ] || [ "$want" != "$got" ]; then
    echo "notebook: checksum mismatch for $file, not installed" >&2
    return 1
  fi
  chmod 755 "$work/$file" || return 1
  mv -f "$work/$file" "$binary" || return 1
  printf '%s\n' "$1" > "$work/stamp" && mv -f "$work/stamp" "$stamp"
}

hash=$(sh "$here/source-hash.sh")
if [ "$(cat "$stamp" 2>/dev/null)" != "$hash" ]; then
  (install_release "$hash")
fi

if [ "${1:-}" = --install ]; then
  exit 0
fi
if [ ! -x "$binary" ]; then
  echo "notebook: no binary at $binary and release notebook-$hash is not available" >&2
  exit 1
fi
exec "$binary" "$@"
