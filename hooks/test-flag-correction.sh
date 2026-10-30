#!/bin/sh
# Pipes sample UserPromptSubmit input into flag-correction.sh and checks
# which prompts land in pending.jsonl, and that the hook prints nothing.
set -u

hook="$(cd "$(dirname "$0")" && pwd)/flag-correction.sh"
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
pending=$work/state/agent-corrections/pending.jsonl
fail=0

# run <prompt>: feeds one prompt to the hook; fails the test on any output.
run() {
  out=$(jq -n --arg p "$1" \
    '{hook_event_name: "UserPromptSubmit", session_id: "s1", cwd: "/repo", prompt: $p}' |
    XDG_STATE_HOME=$work/state sh "$hook")
  if [ -n "$out" ]; then
    echo "FAIL  printed output for: $1"
    fail=1
  fi
}

# expect <flag|pass> <prompt>
expect() {
  before=$(cat "$pending" 2>/dev/null | wc -l)
  run "$2"
  after=$(cat "$pending" 2>/dev/null | wc -l)
  if [ "$1" = flag ] && [ "$after" -ne $((before + 1)) ]; then
    echo "FAIL  not flagged: $2"
    fail=1
  elif [ "$1" = pass ] && [ "$after" -ne "$before" ]; then
    echo "FAIL  flagged: $2"
    fail=1
  else
    echo "ok    $1: $2"
  fi
}

expect flag "no, I said the staging table, not prod"
expect flag "That's wrong, the column is user_id"
expect flag "don't push yet"
expect flag "Why did you delete the test?"
expect flag "you didn't run the check"
expect flag "stop, revert that"
expect flag "I told you to use jq"
expect pass "add a test for the parser"
expect pass "now push it"
expect pass "nothing else, thanks"
expect pass "what does the stop hook do?"

if ! jq -e 'select(.session_id == "s1" and .cwd == "/repo" and (.prompt | length) > 0 and (.ts | test("^[0-9]{4}-")))' "$pending" >/dev/null; then
  echo "FAIL  a pending line lacks ts, session_id, cwd or prompt"
  fail=1
fi

if [ "$fail" -ne 0 ]; then
  echo "correction flag tests: failed"
  exit 1
fi
echo "correction flag tests: passed"
