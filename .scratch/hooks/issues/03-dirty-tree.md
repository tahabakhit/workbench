# 03 Ask before destructive git on a dirty tree

Status: open
Blocked by: 01

## What

Read `spec.md` first. Add the dirty-tree row to `hooks/guard-bash.sh`, with
the tests written first in `hooks/test-safety.sh`. `git reset --hard`,
`git clean -f` (including `-fd` and `-fdx`), `git checkout -- .` and
`git restore .` get "ask" only when `git status --porcelain` in the
command's repo prints anything. A `cd <dir> &&` prefix sets that repo.

## Done when

- `sh scripts/check.sh` passes with these cases:
  - Each of the four commands gets "ask" in a repo with an untracked or
    modified file.
  - Each gets no reply in a clean repo.
  - `git reset --soft HEAD~1` and `git checkout main` get no reply.
- `jq '.hooks.PreToolUse | length' hooks/hooks.json` still prints `1`.
