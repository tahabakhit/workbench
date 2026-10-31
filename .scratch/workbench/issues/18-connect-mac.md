# 18 Connect the Mac to the public repo

Status: wontfix

## What

The repo is public. On the Mac, replace the `directory` marketplace at
`~/dev/workbench` with the `tahabakhit/workbench` `github` source, turn on
auto-update, and remove the old checkout once nothing reads it. Check first
that the Codex `contract.md` include in `~/.codex/AGENTS.md` no longer needs
it. Nothing else waits on this.

How plugin updates authenticate (from the Claude Code docs on hosting a
marketplace): Claude Code has no token of its own and runs the machine's
`git` with prompts off, trying SSH first, then HTTPS through a credential
helper. A public repo needs no credentials, so the Mac's work-account
default does not matter. An update that would need a prompt fails quietly
and keeps the old copy. Auto-update is off by default for a third-party
marketplace; turn it on in `/plugin` under Marketplaces.

## Done when

- On the Mac, `claude plugin marketplace list` shows `workbench` with a
  `github` source, and `~/dev/workbench` does not exist.
- A commit pushed from Fedora reaches a new Mac session without a manual
  step: `claude plugin details workbench` shows the change.

## Out of scope

- The settings ownership itself (mac-dotfiles issue 04).
