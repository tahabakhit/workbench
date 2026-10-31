#!/bin/sh
# SessionStart hook: plain stdout becomes session context, capped at 10,000
# characters, so contract.md must stay under 8,000 to leave room for the
# notebook push. The push is optional and never fails the hook.
root=$(cd "$(dirname "$0")/.." && pwd)
cat "$root/contract.md"

[ "${WORKBENCH_NOTEBOOK_PUSH:-on}" = off ] && exit 0
binary=${XDG_DATA_HOME:-$HOME/.local/share}/workbench/bin/notebook
[ -x "$binary" ] || exit 0
# The common git folder names the repo from inside a linked worktree too.
common=$(git rev-parse --path-format=absolute --git-common-dir 2>/dev/null) || exit 0
repo=$(basename "$(dirname "$common")")
push=$("$binary" push --repo "$repo" 2>/dev/null) || exit 0
[ -n "$push" ] && printf '\n%s\n' "$push"
exit 0
