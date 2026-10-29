# 01 Safety hook: ask on force push and remote branch delete

Status: open

## What

Read `spec.md` (Principles and Decisions) first.

- Check the installed Claude Code's hook docs for the `PreToolUse` reply
  format and the `if` field before writing anything.
- Write `hooks/test-safety.sh` first, with cases that fail: `git push
  --force`, `git push -f origin main`, `git push --force-with-lease`,
  `git push origin +main`, `git push --delete origin b` and
  `git push origin :b` each get "ask". `git push`, `git push origin main`
  and `ls` get no reply.
- Add `hooks/guard-bash.sh` and one `PreToolUse` entry to `hooks/hooks.json`
  next to the existing `SessionStart` entry. Its `if` matchers cover every
  command the spec lists, so 02 to 04 only add functions and tests.
- Make `scripts/check.sh` run `hooks/test-safety.sh`.
- Don't change `hooks/session-context.sh`.

## Done when

- `sh scripts/check.sh` passes and prints the test results.
- `jq '.hooks.PreToolUse | length' hooks/hooks.json` prints `1`.
- In a session started with `claude --plugin-dir ~/dev/workbench`, inside a
  throwaway repo with a local bare remote, `git push --force` and `git push
  origin :b` show an approval prompt naming the reason, and a plain `git
  push` runs with no prompt. Quote what the session showed.
