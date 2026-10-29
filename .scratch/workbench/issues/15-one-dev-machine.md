# 15 Develop on the Fedora PC only; decide how the Mac consumes the plugin

Status: open
Blocked by: 14

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
- The repo is private. The Mac's default `gh` and git identity is a work
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

Visibility audit (all history, all files):

- gitleaks over all history: no leaks. All commit authors use the GitHub
  noreply address.
- Skills mention only generic tool names.
- `.scratch/workbench/trials/onboard.md` and `record.md` describe how the
  work machine's tools are wired (tracker CLI, chat through a company MCP
  gateway, a pending data-handling check, a work tracker project key).
- History before the commit that put placeholders into the onboard trial
  holds the employer's name and its internal GitHub Enterprise host.
  Publishing needs those files scrubbed and that history rewritten (squash or
  `git filter-repo`), then a force push and fresh clones on both machines.
  Irreversible; needs explicit approval.

Also noticed:

- The Fedora PC enables the `mattpocock-skills` plugin, which overlaps with
  several workbench skills.
- Issue 08's done-when says ten skills; the plugin ships eleven.

## Done when

- The chosen Mac option is written to `## Decisions` in `spec.md`, with the
  visibility decision, after approval.
- On the Mac, a commit pushed from Fedora reaches a new Claude Code session
  without a manual step: `claude plugin details workbench` shows the new
  skill or change.
- `git -C ~/dev/workbench log origin/main..main` on the Mac prints nothing.

## Out of scope

- Doing the history rewrite or publishing without explicit approval.
- The settings ownership itself (mac-dotfiles issue 04).
