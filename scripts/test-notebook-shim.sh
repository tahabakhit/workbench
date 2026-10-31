#!/bin/sh
# Tests notebook/source-hash.sh and notebook/run.sh. A file:// release stands in
# for GitHub, and a shell script for the binary, so this needs no network or Go.
set -eu

root=$(cd "$(dirname "$0")/.." && pwd)
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

fail() {
  echo "FAIL  $1" >&2
  exit 1
}
ok() { echo "ok    $1"; }

# source-hash.sh, on a known tree. The expected value is from:
#   printf '%s  ./a.txt\n%s  ./sub/b.txt\n' <sha256 of "a\n"> <sha256 of "b\n"> | sha256sum
tree=$work/tree
mkdir -p "$tree/sub"
printf 'a\n' > "$tree/a.txt"
printf 'b\n' > "$tree/sub/b.txt"
mkdir "$work/hasher"
{ printf '#!/bin/sh\ncd "%s"\n' "$tree"; sed -n '/^if command/,$p' "$root/notebook/source-hash.sh"; } > "$work/hasher/hash.sh"
[ "$(sh "$work/hasher/hash.sh")" = 226010271c020b22 ] || fail "known tree hash: $(sh "$work/hasher/hash.sh")"
ok "known tree hash"
printf 'c\n' > "$tree/a.txt"
[ "$(sh "$work/hasher/hash.sh")" != 226010271c020b22 ] || fail "hash ignores content"
ok "hash follows content"

# run.sh, from a copy of notebook/ with a release built to match it.
plugin=$work/plugin
mkdir "$plugin"
cp "$root/notebook/run.sh" "$root/notebook/source-hash.sh" "$plugin/"
printf 'v1\n' > "$plugin/source.txt"

case $(uname -m) in x86_64 | amd64) arch=amd64 ;; *) arch=arm64 ;; esac
file=notebook-$(uname -s | tr '[:upper:]' '[:lower:]')-$arch

release() { # release <version text>
  hash=$(sh "$plugin/source-hash.sh")
  dir=$work/releases/notebook-$hash
  mkdir -p "$dir"
  printf '#!/bin/sh\necho "%s $*"\n' "$1" > "$dir/$file"
  (cd "$dir" && { sha256sum "$file" 2>/dev/null || shasum -a 256 "$file"; } > checksums.txt)
}

export XDG_DATA_HOME="$work/data"
export NOTEBOOK_RELEASE_URL="file://$work/releases"
binary=$XDG_DATA_HOME/workbench/bin/notebook

if sh "$plugin/run.sh" search x 2>"$work/err"; then fail "ran with no binary and no release"; fi
grep -q 'not available' "$work/err" || fail "no message without a binary"
ok "no binary and no release fails with a message"

release one
[ "$(sh "$plugin/run.sh" search x)" = "one search x" ] || fail "first install"
ok "installs the matching release and starts it"

rm -rf "$work/releases"
[ "$(sh "$plugin/run.sh" recall y)" = "one recall y" ] || fail "second run"
ok "a matching binary needs no download"

printf 'v2\n' > "$plugin/source.txt"
[ "$(sh "$plugin/run.sh" search x)" = "one search x" ] || fail "release not up"
ok "keeps the installed binary while the release is not up"

release two
dir=$work/releases/notebook-$(sh "$plugin/source-hash.sh")
sed 's/^./0/' "$dir/checksums.txt" > "$dir/bad" && mv "$dir/bad" "$dir/checksums.txt"
got=$(sh "$plugin/run.sh" search x 2>"$work/err")
[ "$got" = "one search x" ] || fail "bad checksum replaced the binary: $got"
grep -q 'checksum mismatch' "$work/err" || fail "bad checksum was silent"
ok "a bad checksum refuses the install"

: > "$dir/checksums.txt"
[ "$(sh "$plugin/run.sh" search x 2>/dev/null)" = "one search x" ] || fail "empty checksums"
ok "a checksums file without the binary refuses the install"

release two
sh "$plugin/run.sh" --install
[ "$("$binary" x)" = "two x" ] || fail "--install"
ok "--install updates the binary and starts nothing"
rm "$binary"
sh "$plugin/run.sh" --install
[ -x "$binary" ] || fail "deleted binary not reinstalled"
ok "a deleted binary is reinstalled though its stamp matches"
[ -z "$(find "$(dirname "$binary")" -name '.install-*')" ] || fail "work folder left"
ok "no work folder left"

rm -rf "$XDG_DATA_HOME"
pids=
for n in 1 2 3 4 5 6 7 8; do
  sh "$plugin/run.sh" search x > "$work/out$n" 2>&1 &
  pids="$pids $!"
done
for pid in $pids; do wait "$pid" || fail "concurrent run failed"; done
for n in 1 2 3 4 5 6 7 8; do
  [ "$(cat "$work/out$n")" = "two search x" ] || fail "concurrent run $n: $(cat "$work/out$n")"
done
ok "concurrent installs all start the binary"

[ -x "$root/notebook/run.sh" ] && [ -x "$root/hooks/install-notebook.sh" ] || fail "scripts not executable"
ok "scripts are executable"

echo "notebook shim tests: passed"
