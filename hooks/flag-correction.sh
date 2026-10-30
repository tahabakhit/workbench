#!/bin/sh
# UserPromptSubmit hook: appends prompts that read like a correction of the
# agent to pending-corrections.jsonl for the record skill, which confirms each one. It
# prints nothing, so the prompt and the session context are unchanged.
pending="${XDG_DATA_HOME:-$HOME/.local/share}/workbench/state/pending-corrections.jsonl"

line=$(jq -c '
  select(.prompt | test(
    "^\\s*(no|nope|wrong|stop|undo|revert|don.?t|do not)\\b"
    + "|\\b(that.?s|this is|it.?s) (not right|wrong|incorrect)\\b"
    + "|\\bi (said|told you|asked( you)?)\\b"
    + "|\\bwhy (did|would) you\\b"
    + "|\\byou (didn.?t|did not|forgot|missed|ignored|should have)\\b"; "i"))
  | {ts: (now | todate), session_id, cwd, prompt: (.prompt | .[0:500])}') || exit 0

if [ -n "$line" ]; then
  mkdir -p "$(dirname "$pending")"
  printf '%s\n' "$line" >>"$pending"
fi
exit 0
