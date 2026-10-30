# Handoff

## Goal

Prove the two untested paths of the per-repo tracker and Codex work: a live Codex session loading the plugin (issue 20), and `onboard` plus `write-spec` publishing to a real tracker.

## State

- Branch `main`, last commit before this handoff `568284e`, level with `origin/main`. No other uncommitted changes.
- Done: per-repo tracker choice, shared-repo rule, `## Not yet specified`, self-noticed struggles in `record`, router check in `check.sh`, stale text fixes.
- Open issues: 11 (notebook tool), 18 (connect the Mac), 20 (contract for Codex and Pi).

## Claims

- verified: `sh scripts/check.sh` exits 0 at `568284e`.
- verified: the router check and the `tracker-and-fog` test each passed, failed on a deliberately broken case, and passed after the revert.
- verified: `.git/info/exclude` hides `.scratch/` in a temporary repo (0 status lines).
- verified: `onboard repo preview` in a throwaway repo with a GitHub remote and no issues asked where tasks live and recommended local. It ran without read access to `skills/onboard/tracker-lines.md`, so the seed lines were not exercised.
- verified: `codex plugin add workbench@workbench` in a temporary `CODEX_HOME` installed 16 skills and `hooks/hooks.json` (`.scratch/workbench/issues/20-contract-for-codex-pi.md`).
- verified: this machine's `~/.local/share/workbench/sources.md` is keyed by kind (`host`, `tracker`).
- assumed: a live Codex session loads the installed skills and runs the `SessionStart` hook. Upstream reports only model-invoked skills loading, and 13 of the 16 skills here are user-invoked.
- assumed: the `write-spec` publish step works against a real tracker. It has never run.
- assumed: the GitLab line in `tracker-lines.md` works. Its commands come from upstream's template and `glab` is not installed here.

## Next

1. Start a new session so the pushed skills load. Done when `/workbench:onboard repo preview` shows the tracker question.
2. Ask the owner before installing the plugin into the real Codex home, or get a way to authenticate a temporary one. Then ask a Codex session "where does a decision made outside a feature folder go?". Done when issue 20's done-when passes or the issue says why the hook is not enough.
3. Trial publishing on a repo the owner names, with a GitHub tracker line and a throwaway issue. Creating an issue needs the owner's approval. Done when an issue file holds a `Tracker:` line whose URL resolves.
4. Run `/workbench:onboard machine` on each other machine. Done when each `sources.md` is keyed by kind.

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
