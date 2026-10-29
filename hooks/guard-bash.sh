#!/bin/sh
# PreToolUse hook on Bash: replies "ask" before commands that do irreversible
# damage and stays silent for everything else. Matching on the command text
# is best effort; quoting and commands built inside scripts are not parsed.

reply() {
  jq -n --arg d "$1" --arg r "$2" \
    '{hookSpecificOutput: {hookEventName: "PreToolUse", permissionDecision: $d, permissionDecisionReason: $r}}'
  exit 0
}

# Sets $sub to the git subcommand of the words in "$@" and $subargs to the
# words after it, skipping git's global options; $sub is empty for non-git.
git_words() {
  sub=
  subargs=
  [ "${1-}" = git ] || return 0
  shift
  while [ $# -gt 0 ]; do
    case $1 in
      -C|-c|--git-dir|--work-tree|--namespace) shift 2 || return 0 ;;
      -*) shift ;;
      *) sub=$1; shift; subargs=$*; return 0 ;;
    esac
  done
}

check_force_push() {
  [ "$sub" = push ] || return 0
  for w in $subargs; do
    case $w in
      --force|--force-with-lease|--force-with-lease=*|+*) ;;
      --*) continue ;;
      -*f*) ;;
      *) continue ;;
    esac
    reply ask "Force push rewrites remote history ($w); the agent rules require the owner's approval."
  done
}

check_remote_delete() {
  [ "$sub" = push ] || return 0
  for w in $subargs; do
    case $w in
      --delete|:?*) ;;
      --*) continue ;;
      -*d*) ;;
      *) continue ;;
    esac
    reply ask "This push deletes a remote branch ($w)."
  done
}

command=$(jq -r '.tool_input.command // empty')
[ -n "$command" ] || exit 0

set -f
segments=$(printf '%s\n' "$command" | tr ';|&' '\n\n\n')
IFS='
'
for segment in $segments; do
  IFS=' 	'
  git_words $segment
  check_force_push
  check_remote_delete
  IFS='
'
done
exit 0
