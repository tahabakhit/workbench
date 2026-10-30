#!/bin/sh
# Builds small feature folders and checks which ones check-spec.sh accepts.
set -u

check="$(cd "$(dirname "$0")" && pwd)/check-spec.sh"
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
fail=0

# feature <name> <row>...: each row is "NN|blocked by|status".
feature() {
  d=$work/$1
  shift
  mkdir -p "$d/issues"
  {
    printf '# Spec\n\n## Goal\n\n## Tasks\n\n| # | Task | Blocked by | Status |\n|---|---|---|---|\n'
    for r in "$@"; do
      IFS='|' read -r nn by st <<EOF
$r
EOF
      printf '| %s | Task %s | %s | %s |\n' "$nn" "$nn" "$by" "$st"
      {
        printf '# %s\n\nStatus: %s\n' "$nn" "$st"
        [ -z "$by" ] || printf 'Blocked by: %s\n' "$by"
        printf '\n## Done when\n\n- It works.\n'
      } >"$d/issues/$nn-task.md"
    done
    printf '\n## Out of scope\n\n## Decisions\n'
  } >"$d/spec.md"
}

# expect <name> pass|fail [text the output must contain]
expect() {
  out=$(sh "$check" "$work/$1" 2>&1)
  rc=$?
  if [ "$2" = pass ] && [ "$rc" -ne 0 ]; then
    echo "FAIL  $1: expected pass, got: $out"
    fail=1
  elif [ "$2" = fail ] && [ "$rc" -eq 0 ]; then
    echo "FAIL  $1: expected a failure"
    fail=1
  elif [ -n "${3:-}" ] && ! printf '%s\n' "$out" | grep -q "$3"; then
    echo "FAIL  $1: output lacks '$3': $out"
    fail=1
  else
    echo "ok    $1"
  fi
}

feature chain '01||done' '02|01|open' '03|02|open'
expect chain pass

feature both-wontfix '01||wontfix' '02|01|wontfix'
expect both-wontfix pass

feature missing-blocker '01|09|open'
expect missing-blocker fail 'which has no issue file'

feature open-behind-wontfix '01||wontfix' '02|01|open'
expect open-behind-wontfix fail 'which is wontfix'

feature two-cycle '01|02|open' '02|01|open'
expect two-cycle fail 'cycle'

feature three-cycle '01|03|open' '02|01|open' '03|02|open'
expect three-cycle fail 'cycle'

feature tracker-and-fog '01||open' '02|01|open'
awk '{ print } /^Status: open$/ { print "Tracker: https://example.test/1" }' \
  "$work/tracker-and-fog/issues/01-task.md" >"$work/issue"
mv "$work/issue" "$work/tracker-and-fog/issues/01-task.md"
awk '/^## Out of scope$/ { print "## Not yet specified\n\n- A vague question.\n" } { print }' \
  "$work/tracker-and-fog/spec.md" >"$work/spec"
mv "$work/spec" "$work/tracker-and-fog/spec.md"
expect tracker-and-fog pass

if [ "$fail" -ne 0 ]; then
  echo "check-spec tests: failed"
  exit 1
fi
echo "check-spec tests: passed"
