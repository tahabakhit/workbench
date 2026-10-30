# 15 Develop on the Fedora PC only; decide how the Mac consumes the plugin

Status: open
Blocked by: 14

The decisions are in `## Decisions` in `spec.md`. Connecting the Mac is
issue 18, which blocks nothing.

## What

Workbench is developed, committed and pushed only from the Fedora PC. The
Mac's `~/dev/workbench` gets no new commits. Decide and set up how the Mac
receives updates, and whether the repo goes public.

Facts gathered on the Mac:

- The Mac installs the plugin from a `directory` marketplace at
  `~/dev/workbench` with `autoUpdate`. Its chezmoi setup (`~/dev/agents`)
  also uses the checkout for the `~/.agents/skills` links and for the
  `contract.md` include in `~/.claude/CLAUDE.md` and `~/.codex/AGENTS.md`.
  Issue 14 removes the links and moves the contract into a plugin hook, which
  leaves only the Codex include line needing a local file.
- The Mac's default `gh` and git identity is a work
  account; only this clone's local config pushes as the personal account. A
  `github` marketplace source on the Mac would clone with the default
  credentials and probably fail to fetch or auto-update. Untested.
- Settings: `workbench@workbench` and the `workbench` marketplace are live in
  the Mac's `~/.claude/settings.json` but not owned by its chezmoi source.
  Its modify script keeps unowned keys, so nothing drops them. The owned
  version is planned in mac-dotfiles issue 04 for both OSes.

Options for the Mac:

1. Keep the checkout, pull only. A Mac `SessionStart` hook could run
   `git -C ~/dev/workbench pull --ff-only` so updates arrive without a manual
   step. Lives in the Mac's chezmoi source.
2. After issue 14, switch the Mac to a `github` marketplace source and drop
   the checkout. Needs the private-repo auth above solved, or a public repo.

Publishing: the old versions of the onboard and record trials held the
employer's name, its GitHub Enterprise host, a work account and details of a
work repo. History was rewritten so every version of those files is the
cleaned one, then force-pushed. A grep of all history for those terms and
gitleaks over all history find nothing. GitHub keeps the old commits
reachable by hash until it garbage-collects them, so the clean history goes
to a fresh public repo; the current one is renamed and stays private.

## Done when

- The Mac option and the visibility decision are in `## Decisions` in
  `spec.md`.
- `gh repo view tahabakhit/workbench --json visibility` shows `PUBLIC`, and
  a grep of all history for the employer terms finds nothing.

## Out of scope

- Connecting the Mac (issue 18).
- The settings ownership itself (mac-dotfiles issue 04).
