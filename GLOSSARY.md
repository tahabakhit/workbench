# workbench

A Claude Code plugin holding one person's skills and the rules every agent
session shares.

## Language

**Contract**:
The rules in `contract.md` that the plugin adds to every session.
_Avoid_: Global rules, agent rules

**Knowledge destination**:
One `- <kind>: <where>, <how>` line in a `## Knowledge` section, naming where
a kind of record is filed or gathered from.
_Avoid_: Sink, target

**Kind**:
The category a Knowledge destination serves, such as `log`, `plan` or
`sources`.

**Side file**:
Reference material a skill reads, copied byte for byte into each skill folder
that needs it.
_Avoid_: Shared file, include

**Done-when**:
The check in an issue that decides whether the task is finished.
_Avoid_: Acceptance criteria
