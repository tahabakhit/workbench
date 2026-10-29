# write-spec trials

## Preview on the bootstrap input

A fresh subagent followed `skills/write-spec/SKILL.md` in preview mode on a
detached checkout of this repo from before `.scratch/workbench/` existed, with
the input `.scratch/bootstrap/decisions.md` and
`.scratch/bootstrap/mattpocock-analysis.md`, scope Order steps 2 and 3. Quiz
answers were assumed "yes". Nothing was written.

| `.scratch/workbench/spec.md` | First run | After the fix |
|---|---|---|
| Goal | Goal | Goal |
| Principles | missing | Principles |
| Skills | Skill rename table | Skill inventory and names |
| Repo layout | Layout | Repo layout |
| (none) | (none) | Contract contents, Machine setup |
| Tasks | Tasks | Tasks |
| Out of scope | Out of scope | Out of scope |
| Decisions | Decisions | Decisions |

The first run left out the cross-cutting rules. The skill then named `## Principles`
as the home for every rule each issue must obey; the rerun carried them over.
Both runs produced eight issues with plain `Status:` lines, `Blocked by`
lines that match the task table, and done-when checks that run as commands
(`sh scripts/check.sh`, `grep -c`, `claude plugin list`). The rerun's extras
are context sections holding settled facts that several issues lean on,
which the skill allows.

## Scenario runs

Each scenario is a throwaway git folder under `~/dev/ref/` holding a settled
`decisions.md` (brownfield is a copy of an existing `todo` CLI repo). Three
headless sessions per scenario, each with Read, Write, Edit and Bash allowed,
and `ssh`, `scp`, `rsync`, `curl`, `git push` and the vault's filing tools
denied:

1. `/workbench:write-spec decisions.md` with this checkout as the plugin
   (`--plugin-dir`). It stopped at the step 3 quiz every time.
2. `yes`, on the same session. It wrote the folder and ran `check-spec.sh`.
3. A fresh session in the folder with only `start issue 01`.

No run failed on the skill, so it was not changed. The decision-document
session first stopped to ask for WebFetch, which the trial had not allowed
and its issue needs for vendor pricing pages; the rerun with WebFetch
allowed is the one below.

| Scenario | Folder | Issues | Quiz flagged |
|---|---|---|---|
| Greenfield code | `~/dev/ref/ws-greenfield/.scratch/dupfind` | 3 | zero-byte files, symlinks, exit codes, sort order |
| Brownfield repo | `~/dev/ref/ws-brownfield/.scratch/complete-and-due` | 3 | an older `todo-cli` folder already held `complete`; new folder, old issue 02 set `wontfix`, 03 cut down |
| Analytics (DuckDB SQL) | `~/dev/ref/ws-analytics/.scratch/weekly-active-listings` | 1 | `run.sh` runs the checks, category list, ignore the database file |
| Homelab change | `~/dev/ref/ws-homelab/.scratch/uptime-monitoring` | 2 | untrusted `tls internal` cert inside Kuma; stop before deploying to the host |
| Decision document | `~/dev/ref/ws-decision/.scratch/password-manager` | 2 | "two pages" as under 1,000 words; carrying out the move left to a later spec |

`check-spec.sh` on each folder afterwards printed `spec and issues agree`
for all five. Every write-spec run ended with the paste-ready "Work issue
01 in ..." prompt.

First reply of the fresh `start issue 01` session, then its outcome:

| Scenario | First reply | Outcome |
|---|---|---|
| Greenfield | "Reading the spec and issue 01." | Tests first (7 failing, then passing), issue closed and committed |
| Brownfield | "Reading the spec and issue 01." | Picked the new folder's 01 because the old `todo-cli` 01 was done; 10 tests pass, committed |
| Analytics | "I'll find the spec and issue 01 first." | Failing check first; 12 rows and cars = 3 match; negative test caught; committed |
| Homelab | "The config is edited and the local compose check passes (`exit=0`, `uptime-kuma` has no `ports`, and `kuma_data` is declared). Next I'm checking that the Caddyfile parses." | Config done, `docker compose config` passes; stopped before touching the host, as the issue says |
| Decision document | "Reading the spec and issue 01 first." | Fetched both pricing pages, wrote `password-manager.md` (766 words, five sections in order), closed and committed |

The homelab session wrote no text until after its edits, so its first reply is a progress line; the transcript ends "Issue 01 is done up to the deploy step".

Side findings, not skill defects: the homelab write-spec run noticed the
existing Caddyfile proxies to `host.docker.internal` without the
`extra_hosts` entry Linux needs, and reported it separately. The brownfield
run edited issues in the older folder only after the quiz approved it.
