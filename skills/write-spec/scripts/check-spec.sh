#!/bin/sh
# Usage: check-spec.sh .scratch/<feature>
# Checks that spec.md and issues/NN-slug.md agree and follow the plain-line
# format that grep-based tooling relies on. Prints each problem; exits 1 on any.
set -eu

dir=${1:?usage: check-spec.sh .scratch/<feature>}
spec=$dir/spec.md
problems=$(mktemp)
trap 'rm -f "$problems"' EXIT
say() { printf '%s\n' "$*" >>"$problems"; }

[ -f "$spec" ] || { echo "no spec: $spec" >&2; exit 1; }

for h in '## Goal' '## Tasks' '## Out of scope' '## Decisions'; do
  grep -qx "$h" "$spec" || say "$spec: missing section '$h'"
done

trim() { sed 's/^[[:space:]]*//; s/[[:space:]]*$//'; }

# blockers <NN>: the issue numbers on that issue's Blocked by line.
blockers() {
  for f in "$dir"/issues/"$1"-*.md; do
    if [ -f "$f" ]; then sed -n 's/^Blocked by: //p' "$f" | tr ',' ' '; fi
  done
}

# Task table rows look like: | NN | title | blockers | status |
rows=$(grep -E '^\| *[0-9]{2} *\|' "$spec" || true)
[ -n "$rows" ] || say "$spec: task table has no rows"

printf '%s\n' "$rows" | while IFS= read -r row; do
  [ -n "$row" ] || continue
  nn=$(printf '%s\n' "$row" | awk -F'|' '{print $2}' | trim)
  blocked=$(printf '%s\n' "$row" | awk -F'|' '{print $(NF-2)}' | trim)
  status=$(printf '%s\n' "$row" | awk -F'|' '{print $(NF-1)}' | trim)
  set -- "$dir"/issues/"$nn"-*.md
  if [ ! -f "$1" ]; then
    say "$spec: row $nn has no issues/$nn-*.md"
    continue
  fi
  issue=$1
  istatus=$(sed -n 's/^Status: //p' "$issue")
  [ "$status" = "$istatus" ] || say "$spec: row $nn status '$status' but $issue says '$istatus'"
  iblocked=$(sed -n 's/^Blocked by: //p' "$issue")
  [ "$blocked" = "$iblocked" ] || say "$spec: row $nn blocked by '$blocked' but $issue says '$iblocked'"
done

for issue in "$dir"/issues/*.md; do
  [ -f "$issue" ] || { say "$dir/issues: no issue files"; break; }
  nn=$(basename "$issue" | cut -c1-2)
  printf '%s\n' "$rows" | grep -qE "^\| *$nn *\|" || say "$issue: no row $nn in the task table"
  n=$(grep -cE '^Status: (open|done|wontfix)$' "$issue" || true)
  [ "$n" = 1 ] || say "$issue: needs exactly one plain 'Status: open|done|wontfix' line, found $n"
  ! grep -q '\*\*Status' "$issue" || say "$issue: bold **Status** breaks grep; use a plain line"
  grep -qx '## Done when' "$issue" || say "$issue: missing '## Done when'"
  status=$(sed -n 's/^Status: //p' "$issue")
  for b in $(blockers "$nn"); do
    set -- "$dir"/issues/"$b"-*.md
    if [ ! -f "$1" ]; then
      say "$issue: blocked by $b, which has no issue file"
    elif [ "$status" = open ] && grep -qx 'Status: wontfix' "$1"; then
      say "$issue: blocked by $b, which is wontfix; wontfix this issue too or name another blocker"
    fi
  done
  # Walk the blockers breadth-first; reaching this issue again is a cycle.
  frontier=$(blockers "$nn")
  seen=
  while [ -n "$frontier" ]; do
    next=
    for b in $frontier; do
      if [ "$b" = "$nn" ]; then
        say "$issue: Blocked by forms a cycle back to $nn"
        next=
        break
      fi
      case " $seen " in *" $b "*) continue ;; esac
      seen="$seen $b"
      next="$next $(blockers "$b")"
    done
    frontier=$(printf '%s\n' $next)
  done
done

if [ -s "$problems" ]; then
  cat "$problems" >&2
  exit 1
fi
echo "spec and issues agree: $dir"
