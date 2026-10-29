#!/bin/sh
# Runs one Claude Code session in plan mode (no edits) and reports whether the
# named skill was invoked. Exit 0 when the result matches the expectation.
set -eu
[ $# -eq 3 ] || { echo "usage: trigger-check.sh <skill> <yes|no> <prompt>" >&2; exit 2; }
skill=$1 expect=$2 prompt=$3
fired=$(claude -p --permission-mode plan --max-budget-usd 1 \
  --output-format stream-json --verbose "$prompt" 2>/dev/null |
  jq -rR --arg s "$skill" 'fromjson? | select(.type == "assistant")
    | .message.content[]? | select(.type == "tool_use" and .name == "Skill")
    | .input.skill | select(. == $s or endswith(":" + $s))' | head -1)
if [ -n "$fired" ]; then got=yes; else got=no; fi
if [ "$got" = "$expect" ]; then
  echo "PASS $skill expected=$expect: $prompt"
else
  echo "FAIL $skill expected=$expect got=$got: $prompt"; exit 1
fi
