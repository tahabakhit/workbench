#!/bin/sh
# Pipes sample PreToolUse input into guard-bash.sh and checks the
# permissionDecision it replies with; "none" means no reply at all.
set -u

root=$(cd "$(dirname "$0")/.." && pwd)
guard="$root/hooks/guard-bash.sh"
shell=$(command -v sh)
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

# Split so that scanning this file does not flag it.
key="AKIA""QYLPMN5HHHFPZAM2"

# A PATH holding only what the guard needs, without gitleaks.
nogitleaks="$work/nogitleaks-bin"
mkdir "$nogitleaks"
for tool in jq git tr sed cat; do
  ln -s "$(command -v "$tool")" "$nogitleaks/$tool"
done

pass=0
fail=0

# run_guard <dir> <command> [PATH]: sets $got and $reason.
run_guard() {
  input=$(jq -n --arg c "$2" --arg d "$1" \
    '{hook_event_name: "PreToolUse", tool_name: "Bash", cwd: $d, tool_input: {command: $c}}')
  out=$(cd "$1" && printf '%s\n' "$input" | PATH=${3:-$PATH} "$shell" "$guard")
  status=$?
  reason=
  if [ "$status" -ne 0 ]; then
    got="exit $status"
  elif [ -z "$out" ]; then
    got=none
  else
    got=$(printf '%s\n' "$out" | jq -r '.hookSpecificOutput.permissionDecision // "invalid"' 2>/dev/null) || got=invalid
    reason=$(printf '%s\n' "$out" | jq -r '.hookSpecificOutput.permissionDecisionReason // ""' 2>/dev/null)
  fi
}

# check <description> <want> [reason must contain] [reason must not contain]
check() {
  ok=1
  [ "$got" = "$2" ] || ok=0
  if [ -n "${3-}" ]; then
    case $reason in *"$3"*) ;; *) ok=0 ;; esac
  fi
  if [ -n "${4-}" ]; then
    case $reason in *"$4"*) ok=0 ;; esac
  fi
  if [ "$ok" -eq 1 ]; then
    pass=$((pass + 1))
    echo "ok    $2: $1"
  else
    fail=$((fail + 1))
    echo "FAIL  $1: want $2, got $got ($reason)"
  fi
}

# expect <decision|none> <command>, run outside any git repo.
expect() {
  run_guard "$work" "$2"
  check "$2" "$1"
}

# new_repo <name>: a repo with one clean commit pushed to a local bare remote.
new_repo() {
  git init -q --bare "$work/$1.git"
  git init -q -b main "$work/$1"
  git -C "$work/$1" config user.email test@example.invalid
  git -C "$work/$1" config user.name test
  echo clean > "$work/$1/readme"
  git -C "$work/$1" add readme
  git -C "$work/$1" commit -qm clean
  git -C "$work/$1" remote add origin "../$1.git"
  git -C "$work/$1" push -q -u origin main 2>/dev/null
}

add_key() {
  printf 'aws_access_key_id = %s\n' "$key" > "$1/config.ini"
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

new_repo commit
r="$work/commit"
echo more >> "$r/readme"
git -C "$r" add readme
run_guard "$r" 'git commit -m "clean change"'
check 'commit of a clean staged change' none
add_key "$r"
git -C "$r" add config.ini
run_guard "$r" 'git commit -m "add config"'
check 'commit of a staged key' deny 'config.ini:1' "$key"
run_guard "$work" "git -C $r commit -m x"
check 'commit of a staged key via git -C' deny 'config.ini:1'
run_guard "$work" "cd $r && git commit -m x"
check 'commit of a staged key after cd' deny 'config.ini:1'
run_guard "$r" 'git commit -m x' "$nogitleaks"
check 'commit without gitleaks on PATH' deny 'gitleaks missing'

new_repo all
r="$work/all"
echo "other = 1" > "$r/config.ini"
git -C "$r" add config.ini
git -C "$r" commit -qm config
add_key "$r"
run_guard "$r" 'git commit -m x'
check 'commit with nothing staged and a key in the tree' none
run_guard "$r" 'git commit -am x'
check 'commit -a picking up an unstaged key' deny 'config.ini'

new_repo upstream
r="$work/upstream"
echo more >> "$r/readme"
git -C "$r" commit -qam more
run_guard "$r" 'git push'
check 'push of a clean unpushed commit' none
add_key "$r"
git -C "$r" add config.ini
git -C "$r" commit -qm key
run_guard "$r" 'git push'
check 'push of an unpushed commit holding a key' deny 'config.ini:1' "$key"
run_guard "$r" 'git push --force'
check 'force push of a key is denied, not asked' deny 'config.ini:1'
run_guard "$r" 'git push' "$nogitleaks"
check 'push without gitleaks on PATH' deny 'gitleaks missing'

new_repo branch
r="$work/branch"
git -C "$r" checkout -q -b feature
add_key "$r"
git -C "$r" add config.ini
git -C "$r" commit -qm key
run_guard "$r" 'git push -u origin feature'
check 'push of a branch with no upstream holding a key' deny 'config.ini:1'

echo "safety hook tests: $pass passed, $fail failed"
[ "$fail" -eq 0 ]
