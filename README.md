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

For Codex, run `codex plugin marketplace add tahabakhit/workbench` and
`codex plugin add workbench@workbench`, then accept the workbench hooks when
`codex` next starts.

For Pi, run `pi install git:github.com/tahabakhit/workbench`. It adds the
contract, the skills and correction flagging. A personal skill with the same
name wins over the package's. It also writes the `researcher` agent to
`~/.pi/agent/agents/`, unless a `researcher.md` of your own is there; that agent
needs the `pi-subagents` and `pi-web-access` packages.

## Notebook

`notebook/` holds a Go tool that searches and files into the notebook. A release
is built on each push that changes it, and the plugin installs the matching
binary to `~/.local/share/workbench/bin/notebook`, so a machine needs no Go. Claude
Code gets it as an MCP server too; Codex and Pi use the command line.

## Use

Run `/workbench:onboard machine` once per machine, and
`/workbench:onboard repo` once per repository to choose where its tasks live:
local markdown, GitHub, GitLab or another tracker. When unsure which skill
fits, run `/workbench:choose-skill <your situation>`.

## Check

```sh
sh scripts/check.sh
```

Validates the plugin with `claude plugin validate --strict`, fails if a side
file copied into several skill folders differs between them or `choose-skill`
misses a skill, and runs the correction flag, session context, quiz shuffle,
spec check, Pi, notebook and notebook shim tests.

## Safety

The plugin ships no git safety: nothing here blocks a force push or scans for
secrets, because plugins cannot carry permission rules. Set those up in each
machine's own Claude Code settings or git hooks.
