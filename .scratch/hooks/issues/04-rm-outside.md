# 04 Ask before `rm` outside the working areas

Status: open
Blocked by: 01

## What

Read `spec.md` first. Add the `rm` row to `hooks/guard-bash.sh`, with the
tests written first in `hooks/test-safety.sh`. `rm` with `-r`, `-R` or `-f`
in any flag group gets "ask" when any target, resolved against the hook
input's `cwd`, falls outside:

- `$CLAUDE_PROJECT_DIR`;
- the session scratchpad, meaning any path under `${TMPDIR:-/tmp}`;
- `/tmp`.

Resolve `..` and `~` without requiring the target to exist. A glob or a
variable the hook can't resolve counts as outside.

## Done when

- `sh scripts/check.sh` passes with these cases:
  - `rm -rf build`, `rm -rf ./a/../b` and `rm -rf /tmp/x` inside the project
    get no reply.
  - `rm -rf ~/docs`, `rm -rf ../other`, `rm -rf /` and `rm -rf "$HOME"` get
    "ask".
  - `rm file.txt` with no `-r`/`-f` gets no reply.
- `jq '.hooks.PreToolUse | length' hooks/hooks.json` still prints `1`.
