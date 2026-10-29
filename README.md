# workbench

A Claude Code plugin, and its own marketplace, with skills for working with
AI coding agents across repos and machines.

## Install

```sh
claude plugin marketplace add tahabakhit/workbench
claude plugin install workbench@workbench
```

Turn on auto-update for the marketplace: run `/plugin` in Claude Code, open
Marketplaces, select `workbench` and enable auto-update. New sessions then
pick up each pushed commit.

To pull a change in by hand, run `claude plugin marketplace update workbench`
and start a new session.

## Agent rules

`agent-rules.md` holds the owner's personal rules for every harness. The
plugin does not load it. Each machine's `~/.claude/CLAUDE.md` includes it,
so it keeps the standing of global instructions in Claude Code, and
`~/.agents/AGENTS.md` includes it with `contract.md` for Codex and Pi.

## Check

```sh
sh scripts/check.sh
```

Validates the plugin with `claude plugin validate --strict`, fails if a side
file copied into several skill folders differs between them, and runs the
safety hook and quiz shuffle tests.
