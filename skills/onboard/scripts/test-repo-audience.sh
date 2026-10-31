#!/bin/sh
# Builds small git repos with a stubbed gh and checks what repo-audience.sh
# reports for each.
set -u

script="$(cd "$(dirname "$0")" && pwd)/repo-audience.sh"
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
fail=0

# The stub answers from $work/visibility/<owner>-<repo>; no file means the
# call failed, as it does when gh is logged out or the repo is hidden.
mkdir -p "$work/bin" "$work/visibility"
cat >"$work/bin/gh" <<'EOF'
#!/bin/sh
f="$STUB_DIR/$(printf '%s' "$3" | tr / -)"
[ -f "$f" ] || exit 1
cat "$f"
EOF
chmod +x "$work/bin/gh"

# repo <name>: a repo with one commit by the owner.
repo() {
  d=$work/$1
  git init -q "$d"
  git -C "$d" config user.name Owner
  git -C "$d" config user.email owner@example.test
  echo x >"$d/README.md"
  git -C "$d" add README.md
  git -C "$d" commit -qm first
}

# by <name> <author>: adds a commit by that author.
by() {
  git -C "$work/$1" commit -q --allow-empty -m other \
    --author="$2 <$2@example.test>"
}

# expect <name> <line the report must hold>...
expect() {
  name=$1
  shift
  out=$(cd "$work/$name" && PATH="$work/bin:$PATH" STUB_DIR="$work/visibility" sh "$script" 2>&1)
  for line in "$@"; do
    if ! printf '%s\n' "$out" | grep -qxF "$line"; then
      echo "FAIL  $name: no line '$line' in: $out"
      fail=1
    fi
  done
}

repo no-remote
expect no-remote 'visibility: none' 'other_authors: 0' 'placement: ask'

repo private-solo
git -C "$work/private-solo" remote add origin git@github.com:o/private-solo.git
echo PRIVATE >"$work/visibility/o-private-solo"
expect private-solo 'visibility: private' 'placement: ask'

repo public-solo
git -C "$work/public-solo" remote add origin https://github.com/o/public-solo.git
echo PUBLIC >"$work/visibility/o-public-solo"
expect public-solo 'visibility: public' 'other_authors: 0' 'placement: public'

repo internal-solo
git -C "$work/internal-solo" remote add origin https://github.com/o/internal-solo
echo INTERNAL >"$work/visibility/o-internal-solo"
expect internal-solo 'visibility: internal' 'placement: public'

repo private-team
git -C "$work/private-team" remote add origin git@github.com:o/private-team.git
echo PRIVATE >"$work/visibility/o-private-team"
by private-team teammate
expect private-team 'other_authors: 1' 'placement: shared'

# Other authors win over public: their instruction files must stay untouched.
repo public-team
git -C "$work/public-team" remote add origin git@github.com:o/public-team.git
echo PUBLIC >"$work/visibility/o-public-team"
by public-team teammate
by public-team second
expect public-team 'visibility: public' 'other_authors: 2' 'placement: shared'

repo bot-only
by bot-only 'dependabot[bot]'
expect bot-only 'other_authors: 0' 'placement: ask'

# The owner under a second email is still the owner.
repo two-emails
git -C "$work/two-emails" commit -q --allow-empty -m other \
  --author='Owner <owner@work.example.test>'
expect two-emails 'other_authors: 0'

# A private fork of a public project is read by that project's people.
repo private-fork
git -C "$work/private-fork" remote add origin git@github.com:o/fork.git
git -C "$work/private-fork" remote add upstream https://github.com/up/project.git
echo PRIVATE >"$work/visibility/o-fork"
echo PUBLIC >"$work/visibility/up-project"
expect private-fork 'visibility: public' 'placement: public'

repo gh-fails
git -C "$work/gh-fails" remote add origin git@github.com:o/hidden.git
expect gh-fails 'visibility: unknown' 'placement: ask'

repo gitlab
git -C "$work/gitlab" remote add origin git@gitlab.com:o/gitlab.git
expect gitlab 'visibility: unknown' 'placement: ask'

repo tracked
mkdir -p "$work/tracked/.scratch/feature"
echo spec >"$work/tracked/.scratch/feature/spec.md"
printf '# Repo\n\n## Sources\n\n- tracker: local, .scratch/ only\n' >"$work/tracked/AGENTS.md"
git -C "$work/tracked" add .
git -C "$work/tracked" commit -qm notes
expect tracked 'tracked: .scratch/ (1 files)' 'tracked: AGENTS.md ## Sources'

repo clean
expect clean 'tracked: none'

if [ "$fail" -ne 0 ]; then
  echo "repo-audience tests: failed"
  exit 1
fi
echo "repo-audience tests: passed"
