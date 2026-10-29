# workbench

A Claude Code plugin, and its own marketplace, with skills for working with
AI coding agents across repos and machines.

## Install

From a local clone:

```sh
git clone https://github.com/tahabakhit/workbench.git ~/dev/workbench
claude plugin marketplace add ~/dev/workbench
claude plugin install workbench@workbench
```

Turn on auto-update for the marketplace: run `/plugin` in Claude Code, open
Marketplaces, select `workbench` and enable auto-update. New sessions then
pick up each commit.

To pull a change in by hand, run `claude plugin marketplace update workbench`
and start a new session.

## Skills

| Skill | Invoked by | Does |
|---|---|---|
| `rephrase` | user | Re-explains the previous message in plain English |

## Check

```sh
sh scripts/check.sh
```

Validates the plugin with `claude plugin validate --strict` and fails if a
side file copied into several skill folders differs between them.
