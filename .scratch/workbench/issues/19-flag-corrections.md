# 19 Flag corrections in every session

Status: done

## What

A `UserPromptSubmit` hook in the plugin, `hooks/flag-correction.sh`, appends
prompts that read like a correction of the agent to
`${XDG_STATE_HOME:-~/.local/state}/agent-corrections/pending.jsonl`, one
JSON line with `ts`, `session_id`, `cwd` and the first 500 characters of the
prompt. It prints nothing. `record` confirms each line against its session
and files real corrections into the notebook's `corrections.md`.
`hooks/test-flag-correction.sh`, run by `scripts/check.sh`, covers what is
flagged and what is not.

## Done when

- `sh scripts/check.sh` passes, and the flag tests failed before the hook
  existed.
- In a live session with the plugin loaded, a correcting prompt adds one line
  to `pending.jsonl` and a plain request adds none.

## Out of scope

- Removing the Mac's own copy of the hook (dotfiles).
