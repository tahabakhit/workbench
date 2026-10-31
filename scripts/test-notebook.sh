#!/bin/sh
# Runs the notebook tests with cgo off, the way the released binary is built.
set -eu

root=$(cd "$(dirname "$0")/.." && pwd)

if ! command -v go >/dev/null 2>&1; then
  echo "notebook tests: skipped, go not found"
  exit 0
fi

cd "$root/notebook"
CGO_ENABLED=0 go test ./...
