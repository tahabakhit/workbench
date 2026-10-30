# workbench

A Claude Code plugin holding one person's skills and the rules every agent
session shares.

## Language

**Contract**:
The rules in `contract.md` that the plugin adds to every session.
_Avoid_: Global rules, agent rules

**Notebook**:
The markdown vault in a git repository at `~/notebook` where skills file
logs, decisions, lessons, corrections, plans, goals and studied topics, and
the tool in this plugin that searches and files into it.
_Avoid_: Knowledge destination, brain

**Sources file**:
`~/.config/workbench/sources.md`, the machine's list of where work comes
from: tracker, chat, calendar, meeting notes, code host.
_Avoid_: Knowledge file

**Side file**:
Reference material a skill reads, copied byte for byte into each skill folder
that needs it.
_Avoid_: Shared file, include

**Done-when**:
The check in an issue that decides whether the task is finished.
_Avoid_: Acceptance criteria
