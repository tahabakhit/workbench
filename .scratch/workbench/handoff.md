# Handoff

## Goal

Prove the untested tracker path: `onboard` plus `write-spec` publishing to a real tracker.

## State

- Branch `main`, level with `origin/main`.
- Done: per-repo tracker choice, shared-repo rule, `## Not yet specified`, self-noticed struggles in `record`, router check in `check.sh`, stale text fixes.
- Open issues: 11 (notebook tool), 18 (connect the Mac).

## Claims

- verified: `sh scripts/check.sh` exits 0 at `568284e`.
- verified: the router check and the `tracker-and-fog` test each passed, failed on a deliberately broken case, and passed after the revert.
- verified: `.git/info/exclude` hides `.scratch/` in a temporary repo (0 status lines).
- verified: a fresh headless session ran `/workbench:onboard repo preview` from the installed `568284e` in a throwaway repo. It read `skills/onboard/tracker-lines.md`, asked where tasks live, recommended local and wrote nothing. Its answer named only the recommendation, not GitHub, GitLab or other as options.
- verified: `tahabakhit/workbench` is public, so a push to `main` needs the owner's approval.
- verified: `codex plugin add workbench@workbench` in a temporary `CODEX_HOME` installed 16 skills and `hooks/hooks.json` (`.scratch/workbench/issues/20-contract-for-codex-pi.md`).
- verified: this machine's `~/.local/share/workbench/sources.md` is keyed by kind (`host`, `tracker`).
- verified: the plugin is installed in the real `~/.codex`. A live `codex exec` session loads the 3 model-invoked skills, not the 13 user-invoked ones, and does not get the contract (issue 20).
- verified: after the owner trusted the hooks, a new `codex exec` session answers the issue 20 question from the contract.
- verified: Pi has the package from git; a plain `pi -p` session answers the same question from the contract, and all 16 skills load. `researcher` runs in Pi; it needs `pi-subagents` and `pi-web-access`.
- assumed: the `write-spec` publish step works against a real tracker. It has never run.
- assumed: the GitLab line in `tracker-lines.md` works. Its commands come from upstream's template and `glab` is not installed here.

## Next

1. Trial publishing on a repo the owner names, with a GitHub tracker line and a throwaway issue. Creating an issue needs the owner's approval. Done when an issue file holds a `Tracker:` line whose URL resolves.
2. Run `/workbench:onboard machine` on each other machine. Done when each `sources.md` is keyed by kind.

## Read first

- `.scratch/workbench/issues/20-contract-for-codex-pi.md`
- `.scratch/workbench/spec.md`, the last rows of `## Decisions`
- `contract.md`, `## Sources`
- `skills/onboard/SKILL.md`, `skills/onboard/tracker-lines.md`
- `skills/write-spec/SKILL.md`, step 5
- commits `820dff2` and `568284e`

## Skills

- `/workbench:onboard repo preview`
- `/workbench:onboard machine`
- `/workbench:write-spec`
