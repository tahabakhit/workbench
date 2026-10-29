#!/bin/sh
# Validates the plugin and marketplace, then checks that every side file
# present in more than one skill folder is byte-identical across them.
set -eu

root=$(cd "$(dirname "$0")/.." && pwd)
cd "$root"

report=$(claude plugin validate --strict --json .) || true
if ! printf '%s\n' "$report" | jq -e '.manifest.file' >/dev/null 2>&1; then
  echo "plugin validate produced no usable report:" >&2
  printf '%s\n' "$report" >&2
  exit 1
fi

# plugin.json has no version on purpose, so that every commit counts as an
# update; that is the only strict warning allowed.
problems=$(printf '%s\n' "$report" | jq -r '
  [.manifest, .contents[]?]
  | map(.file as $f
      | (.errors[] | "error: \($f): \(.path): \(.message)"),
        (.warnings[]
          | select((.path | test("(^|→ )version$")) and (.message | startswith("No version specified")) | not)
          | "warning: \($f): \(.path): \(.message)"))
  | .[]')

if [ -n "$problems" ]; then
  printf '%s\n' "$problems" >&2
  exit 1
fi
echo "plugin validate --strict: passes apart from the intended missing version"

# SKILL.md and agents/openai.yaml are per-skill by design; every other file
# with the same path relative to its skill folder is a shared copy.
mismatch=$(find skills -mindepth 2 -type f \
  ! -path 'skills/*/SKILL.md' ! -path 'skills/*/agents/openai.yaml' \
  | sed 's|^skills/[^/]*/||' | sort | uniq -d \
  | while IFS= read -r rel; do
      first=
      for f in skills/*/"$rel"; do
        [ -f "$f" ] || continue
        if [ -z "$first" ]; then
          first=$f
        elif ! cmp -s "$first" "$f"; then
          echo "side file differs: $first vs $f"
        fi
      done
    done)

if [ -n "$mismatch" ]; then
  printf '%s\n' "$mismatch" >&2
  exit 1
fi
echo "side files: all shared copies match"
