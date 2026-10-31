#!/bin/sh
# Reports who can read the current repository and where personal workbench
# configuration belongs. Read-only. placement is one of:
#   shared  other authors work here
#   public  no other authors, but people without workbench can read it
#   ask     the repo cannot say; only the owner knows whether it will go public
set -u

git rev-parse --git-dir >/dev/null 2>&1 || {
  echo "repo-audience: not a git repository" >&2
  exit 1
}

# GitHub only: another host's visibility stays unknown, which leads to asking.
visibility_of() {
  case $1 in
    *github.com[:/]*) ;;
    *) echo unknown; return ;;
  esac
  slug=$(printf '%s\n' "$1" | sed 's#^.*github\.com[:/]##; s#/$##; s#\.git$##')
  v=$(gh repo view "$slug" --json visibility --jq .visibility 2>/dev/null) || v=
  if [ -n "$v" ]; then
    printf '%s\n' "$v" | tr '[:upper:]' '[:lower:]'
  else
    echo unknown
  fi
}

# The most open remote decides, so a private fork of a public project is public.
seen=$(for r in $(git remote); do
    visibility_of "$(git remote get-url "$r")"
  done)
visibility=none
for v in private unknown internal public; do
  if printf '%s\n' "$seen" | grep -qx "$v"; then
    visibility=$v
  fi
done

me_name=$(git config user.name || true)
me_email=$(git config user.email || true)
others=$(git log -30 --format='%an%x09%ae' 2>/dev/null \
  | awk -F '\t' -v n="$me_name" -v e="$me_email" \
      '$1 != n && $2 != e && $1 !~ /\[bot\]$/ { print $2 }' \
  | sort -u | wc -l | tr -d ' ')

if [ "$others" -gt 0 ]; then
  placement=shared
elif [ "$visibility" = public ] || [ "$visibility" = internal ]; then
  placement=public
else
  placement=ask
fi

echo "visibility: $visibility"
echo "other_authors: $others"
echo "placement: $placement"

tracked=0
count=$(git ls-files .scratch | wc -l | tr -d ' ')
if [ "$count" -gt 0 ]; then
  echo "tracked: .scratch/ ($count files)"
  tracked=1
fi
for f in AGENTS.md CLAUDE.md; do
  if git ls-files --error-unmatch "$f" >/dev/null 2>&1 \
    && grep -q '^## Sources' "$f"; then
    echo "tracked: $f ## Sources"
    tracked=1
  fi
done
[ "$tracked" -eq 1 ] || echo "tracked: none"
