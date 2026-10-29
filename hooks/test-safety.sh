#!/bin/sh
# Pipes sample PreToolUse input into guard-bash.sh and checks the
# permissionDecision it replies with; "none" means no reply at all.
set -u

root=$(cd "$(dirname "$0")/.." && pwd)
guard="$root/hooks/guard-bash.sh"
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

pass=0
fail=0

# expect <decision|none> <command>
expect() {
  want=$1
  cmd=$2
  input=$(jq -n --arg c "$cmd" --arg d "$work" \
    '{hook_event_name: "PreToolUse", tool_name: "Bash", cwd: $d, tool_input: {command: $c}}')
  out=$(cd "$work" && printf '%s\n' "$input" | sh "$guard")
  status=$?
  if [ -z "$out" ]; then
    got=none
  else
    got=$(printf '%s\n' "$out" | jq -r '.hookSpecificOutput.permissionDecision // "invalid"' 2>/dev/null) || got=invalid
  fi
  if [ "$status" -eq 0 ] && [ "$got" = "$want" ]; then
    pass=$((pass + 1))
    echo "ok    $want: $cmd"
  else
    fail=$((fail + 1))
    echo "FAIL  $cmd: want $want, got $got (exit $status)"
  fi
}

expect ask 'git push --force'
expect ask 'git push -f origin main'
expect ask 'git push --force-with-lease'
expect ask 'git push origin +main'
expect ask 'git push --delete origin b'
expect ask 'git push origin :b'
expect ask 'git push -d origin b'
expect ask 'git fetch && git push --force origin main'
expect none 'git push'
expect none 'git push origin main'
expect none 'git push -u origin feature'
expect none 'git push origin main:main'
expect none 'ls'

echo "safety hook tests: $pass passed, $fail failed"
[ "$fail" -eq 0 ]
