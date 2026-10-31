#!/bin/sh
# Tests hooks/session-context.sh with a stand-in for the notebook binary.
set -eu

hook=$(cd "$(dirname "$0")" && pwd)/session-context.sh
contract=$(cd "$(dirname "$0")/.." && pwd)/contract.md
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
export XDG_DATA_HOME="$work/data"
binary=$XDG_DATA_HOME/workbench/bin/notebook

fail() {
  echo "FAIL  $1" >&2
  exit 1
}
ok() { echo "ok    $1"; }

stub() { # stub <body>
  mkdir -p "$(dirname "$binary")"
  printf '#!/bin/sh\n%s\n' "$1" > "$binary"
  chmod +x "$binary"
}

export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null
export GIT_AUTHOR_NAME=T GIT_AUTHOR_EMAIL=t@example.invalid GIT_COMMITTER_NAME=T GIT_COMMITTER_EMAIL=t@example.invalid
repo=$work/myrepo
git init -q "$repo"
git -C "$repo" commit -q --allow-empty -m first
git -C "$repo" worktree add -q "$work/elsewhere/area" -b area
cd "$repo"

only_contract() { # only_contract <label>
  sh "$hook" > "$work/out" || fail "$1: exit $?"
  cmp -s "$work/out" "$contract" || fail "$1: output differs from the contract"
  ok "$1"
}

only_contract "missing binary prints the contract alone"

stub 'echo "## Notebook args: $*"'
{ cat "$contract"; printf '\n## Notebook args: push --repo myrepo\n'; } > "$work/want"
sh "$hook" > "$work/out"
cmp -s "$work/out" "$work/want" || fail "push is appended"
ok "push follows the contract, with the repo name"

(cd "$work/elsewhere/area" && sh "$hook" | grep -q -- '--repo myrepo') || fail "worktree repo name"
ok "a linked worktree gives the repo's name"

(cd "$repo" && WORKBENCH_NOTEBOOK_PUSH=off sh "$hook" > "$work/out")
cmp -s "$work/out" "$contract" || fail "off switch"
ok "off switch prints the contract alone"

stub 'exit 0'
only_contract "empty result prints the contract alone"

stub 'echo partial; exit 1'
only_contract "failing binary prints the contract alone"

stub 'echo pushed'
cd "$work"
only_contract "outside a repository prints the contract alone"

size=$(wc -m < "$contract")
[ $((size + 2000 + 2)) -lt 10000 ] || fail "contract plus push: $size + 2000"
ok "contract plus the push fits the cap ($size + 2000)"

echo "session context tests: passed"
