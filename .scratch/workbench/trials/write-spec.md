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

Open: greenfield code, brownfield repo, analytics, homelab change, non-code
decision document.
