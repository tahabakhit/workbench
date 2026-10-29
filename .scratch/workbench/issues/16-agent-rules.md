# 16 Ship the owner's agent rules as `agent-rules.md`

Status: open

Workbench owns the text of the owner's personal agent rules (Autonomy,
Working practices, Git, Tools). The dotfiles include it: `~/.claude/CLAUDE.md`
includes only this file, and `~/.agents/AGENTS.md` includes it with
`contract.md` for Codex and Pi.

## Done when

- `agent-rules.md` is at the repo root and matches the text the dotfiles
  rendered before the hand-over.
- The dotfiles session has the path.
- `sh scripts/check.sh` passes.

## Out of scope

- The include lines themselves (dotfiles).
- Rewording the rules.
