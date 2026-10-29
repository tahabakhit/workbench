# 02 Gitleaks on commit and push

Status: open
Blocked by: 01

## What

Read `spec.md` first. Add the gitleaks rows to `hooks/guard-bash.sh`, with
the tests written first in `hooks/test-safety.sh`.

- `git commit` runs `gitleaks git --staged --redact` in the command's repo.
- `git push` scans `@{u}..HEAD`. When the branch has no upstream, it scans
  `HEAD --not --remotes`.
- A finding means deny, and the reason names each file and line with the
  secret redacted. Gitleaks missing from `PATH` also means deny, with the
  reason "gitleaks missing".
- Check `gitleaks git --help` for the installed version's flags (8.30.1 on
  the Fedora PC).

## Done when

- `sh scripts/check.sh` passes with these cases:
  - A staged fake AWS key is denied on commit, and the reason shows a file
    and line but not the key.
  - A clean staged change gets no reply.
  - An unpushed commit holding the fake key is denied on push.
  - With `PATH` not containing gitleaks, commit and push are denied with
    "gitleaks missing".
- `jq '.hooks.PreToolUse | length' hooks/hooks.json` still prints `1`.
