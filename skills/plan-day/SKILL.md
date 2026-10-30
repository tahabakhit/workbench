---
name: plan-day
description: Sweep the machine's work sources (tracker, notes, calendar, chat, open specs), rank what needs doing now, save a short daily plan, and propose tracker updates for work that is already done. Run at the start of a session with no set task.
disable-model-invocation: true
---

# Plan day

A short brief, not a report. The user reads at most five ranked items and a
handful of proposed updates, then picks what to work on.

## 1. Date and window

Get the date, weekday and time from the shell (`date`); never infer them.
Open the most recent plan in `~/notebook/plans/` from the last seven days.
Its date is the window start; with no plan, use 14 days ago. Read
`~/notebook/goals.md` when it exists.

Done when today's date, the window start and the previous plan (or its
absence) are known.

## 2. Gather

Sources are the lines in `~/.config/workbench/sources.md` that apply to
`plan-day`, plus three every machine has: the previous plan's open items, open
tasks in `.scratch/*/spec.md` of the current repo, and the user's commits in
the window. Follow any limit written beside a source.

For trackers, gather both open items and items updated in the window. For
chat, gather mentions, direct messages and threads the user is in, and check
whether the user already replied.

When the harness has subagents and there are three or more sources, give each
source to a read-only subagent with the window, the exact query and this
return shape: item, link, status, owner, deadline, who is waiting, last
activity, and whether the user already acted. Otherwise query them in turn.

Done when each source returned items or has a one-line note such as
"Chat: skipped (not connected)".

## 3. Reconcile

Merge items across sources: tracker key first, then a link, branch or PR,
then a close title. Then sort each into one case:

| Case | Action |
|---|---|
| In a source, not in the tracker or plan | Offer to add to the plan, or to create a ticket if it is real work |
| In the plan or tracker, nothing seen in the window | Flag as possibly stale |
| Done elsewhere, still open in the tracker | Propose a transition or comment, with the evidence linked |
| Same item twice | Keep one |

Evidence of "done" depends on the task: a merged PR, a filed work note, a
delivered file, or the user saying so in a meeting note or chat. With no
linkable evidence, propose nothing.

Done when every item has a case and every proposed update has a link.

## 4. Rank

- **Now**: someone is blocked or waiting on the user, overdue, due today, or
  a meeting today needs preparation.
- **Today**: due within three days, or a carry-over from the previous plan.
- **Later**: everything else, shown as a count with the top three named.

Within a tier, order by who is waiting, then deadline, then fit with the
goals the instructions name (read that goals note if one is declared), then
how long it has waited. Chat already answered drops out.

Done when Now plus Today hold at most five items; the rest move to Later.

## 5. Present

1. Ranked table: tier, item, why now, next action, link.
2. Proposed updates, numbered: ticket, the exact change or comment text,
   evidence link.
3. Possibly stale items and new items without a ticket.
4. Skipped sources.

Ask which updates to apply and whether the ranking needs changes.

## 6. Apply and save

Apply only the approved updates, with the user's edits, and report each
result. Then write the plan to `~/notebook/plans/YYYY-MM-DD.md` under a
`## Plan` heading, replacing that section if today's file already has one,
and commit it. Open the
section with a Mermaid `mindmap`: the root is the date, branches are Now,
Today and Later, leaves are each item in a few words. Keep brackets out of
labels, since they change node shapes. Below it list the ranked items with
links, carried-over items, and the updates applied. Items from a
source whose limit forbids storing its content appear as a link and a
few-word action, or not at all if the limit says so.

Done when approved updates are applied and the plan is committed in the notebook.
