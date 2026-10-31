#!/bin/sh
# Prints the hash that names this source's release: the same files give the
# same value on Linux and macOS, with or without a git checkout.
set -eu

cd "$(dirname "$0")"

if command -v sha256sum >/dev/null 2>&1; then
  sha() { sha256sum "$@"; }
else
  sha() { shasum -a 256 "$@"; }
fi

find . -type f | LC_ALL=C sort | while IFS= read -r file; do
  sha "$file"
done | sha | cut -c1-16
