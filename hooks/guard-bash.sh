#!/bin/sh
# PreToolUse hook on Bash: replies "deny" when gitleaks flags a commit or push
# or is missing, "ask" before commands that do irreversible damage, and stays
# silent for everything else. Matching on the command text is best effort;
# quoting and commands built inside scripts are not parsed.

reply() {
  jq -n --arg d "$1" --arg r "$2" \
    '{hookSpecificOutput: {hookEventName: "PreToolUse", permissionDecision: $d, permissionDecisionReason: $r}}'
  exit 0
}

# Keeps the first reason; it is sent only after every deny check has run.
ask() {
  [ -n "$ask_reason" ] || ask_reason=$1
}

# Sets $sub to the git subcommand of the words in "$@", $subargs to the words
# after it and $gitdir to a -C directory; $sub is empty for non-git.
git_words() {
  sub=
  subargs=
  gitdir=
  [ "${1-}" = git ] || return 0
  shift
  while [ $# -gt 0 ]; do
    case $1 in
      -C) gitdir=${2-}; shift 2 || return 0 ;;
      -c|--git-dir|--work-tree|--namespace) shift 2 || return 0 ;;
      -*) shift ;;
      *) sub=$1; shift; subargs=$*; return 0 ;;
    esac
  done
}

# Sets $repo to the top of the work tree the git command runs in.
in_repo() {
  case $gitdir in
    '') target=$dir ;;
    /*) target=$gitdir ;;
    *) target=$dir/$gitdir ;;
  esac
  repo=$(git -C "$target" rev-parse --show-toplevel 2>/dev/null) && [ -n "$repo" ]
}

need_gitleaks() {
  command -v gitleaks >/dev/null 2>&1 ||
    reply deny "gitleaks missing: install it; the agent rules forbid committing or pushing without a gitleaks scan."
}

# Prints one line per finding of a gitleaks git scan of $repo; a scan that
# fails prints nothing, so it never causes a deny.
leaks() {
  gitleaks git --no-banner --redact --log-level error --report-format json --report-path - "$@" "$repo" 2>/dev/null |
    jq -r '.[] | "- \(.File):\(.StartLine) \(.RuleID)" + (if (.Commit // "") != "" then " in commit \(.Commit[0:7])" else "" end)' 2>/dev/null
}

deny_leaks() {
  [ -n "$1" ] || return 0
  reply deny "gitleaks found secrets in the $2 (redacted):
$1"
}

check_commit_secrets() {
  [ "$sub" = commit ] && in_repo || return 0
  need_gitleaks
  found=$(leaks --staged)
  for w in $subargs; do
    case $w in
      --all) ;;
      --*) continue ;;
      -*a*) ;;
      *) continue ;;
    esac
    found="$found${found:+
}$(leaks --pre-commit)"
    break
  done
  deny_leaks "$found" "staged changes"
}

check_push_secrets() {
  [ "$sub" = push ] && in_repo || return 0
  need_gitleaks
  if git -C "$repo" rev-parse -q --verify '@{u}' >/dev/null 2>&1; then
    range='@{u}..HEAD'
  else
    range='HEAD --not --remotes'
  fi
  deny_leaks "$(leaks --log-opts="$range")" "commits being pushed"
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
    ask "Force push rewrites remote history ($w); the agent rules require the owner's approval."
    return 0
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
    ask "This push deletes a remote branch ($w)."
    return 0
  done
}

input=$(cat)
command=$(printf '%s\n' "$input" | jq -r '.tool_input.command // empty')
[ -n "$command" ] || exit 0
dir=$(printf '%s\n' "$input" | jq -r '.cwd // empty')
[ -n "$dir" ] || dir=$PWD
ask_reason=

set -f
segments=$(printf '%s\n' "$command" | tr ';|&' '\n\n\n')
IFS='
'
for segment in $segments; do
  IFS=' 	'
  set -- $segment
  if [ "${1-}" = cd ]; then
    case ${2-} in
      '') dir=$HOME ;;
      /*) dir=$2 ;;
      *) dir=$dir/$2 ;;
    esac
  fi
  git_words "$@"
  check_commit_secrets
  check_push_secrets
  check_force_push
  check_remote_delete
  IFS='
'
done

[ -z "$ask_reason" ] || reply ask "$ask_reason"
exit 0
