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

# Prints why the git command throws away uncommitted work, or nothing.
destroys_work() {
  case $sub in
    reset)
      for w in $subargs; do
        [ "$w" = --hard ] && { echo "reset --hard"; return; }
      done ;;
    clean)
      force=
      for w in $subargs; do
        case $w in
          --dry-run) return ;;
          --force) force=$w ;;
          --*) ;;
          -*n*) return ;;
          -*f*) force=$w ;;
        esac
      done
      [ -z "$force" ] || echo "clean $force" ;;
    checkout)
      for w in $subargs; do
        case $w in --|.) echo "checkout $w"; return ;; esac
      done ;;
    restore)
      staged=
      for w in $subargs; do
        case $w in
          --worktree|-W) staged= ; break ;;
          --staged|-S) staged=1 ;;
        esac
      done
      [ -n "$staged" ] || echo "restore" ;;
  esac
}

check_dirty_tree() {
  what=$(destroys_work)
  [ -n "$what" ] && in_repo || return 0
  [ -n "$(git -C "$repo" status --porcelain 2>/dev/null)" ] || return 0
  ask "git $what discards uncommitted changes in $repo."
}

# Prints path $1 made absolute against $dir with . and .. removed, without
# touching the file system.
normalize() {
  case $1 in
    /*) path=$1 ;;
    *) path=$dir/$1 ;;
  esac
  out=
  IFS=/
  for part in $path; do
    case $part in
      ''|.) ;;
      ..) out=${out%/*} ;;
      *) out=$out/$part ;;
    esac
  done
  printf '%s\n' "${out:-/}"
}

# True when $target is strictly inside directory $1.
under() {
  [ -n "$1" ] || return 1
  root=$(normalize "$1")
  case $target in "${root%/}"/?*) return 0 ;; esac
  return 1
}

# True when $1 is strictly inside the project, the scratchpad or /tmp. A
# glob, brace, variable or ~user the hook can't resolve counts as outside.
allowed_target() {
  case $1 in
    *'$'*|*'*'*|*'?'*|*'['*|*'{'*|*'`'*|'~'[!/]*) return 1 ;;
    '~') target=$HOME ;;
    '~/'*) target=$HOME/${1#'~/'} ;;
    *) target=$1 ;;
  esac
  target=$(normalize "$target")
  under "${CLAUDE_PROJECT_DIR:-}" || under "${TMPDIR:-/tmp}" || under /tmp
}

check_rm() {
  [ "${1-}" = rm ] || return 0
  shift
  recursive_or_force=
  targets=
  options=1
  for w in "$@"; do
    w=$(printf '%s' "$w" | tr -d "\"'")
    if [ -n "$options" ]; then
      case $w in
        --) options= ; continue ;;
        --recursive|--force) recursive_or_force=1; continue ;;
        --*) continue ;;
        -*[rRf]*) recursive_or_force=1; continue ;;
        -*) continue ;;
      esac
    fi
    allowed_target "$w" || targets="$targets $w"
  done
  [ -n "$recursive_or_force" ] && [ -n "$targets" ] || return 0
  ask "rm deletes outside the project, the scratchpad and /tmp:$targets"
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
    target=$(printf '%s' "${2-}" | tr -d "\"'")
    case $target in
      '') dir=$HOME ;;
      /*) dir=$target ;;
      *) dir=$dir/$target ;;
    esac
  fi
  git_words "$@"
  check_commit_secrets
  check_push_secrets
  check_force_push
  check_remote_delete
  check_dirty_tree
  check_rm "$@"
  IFS='
'
done

[ -z "$ask_reason" ] || reply ask "$ask_reason"
exit 0
