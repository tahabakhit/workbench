---
name: write-spec
description: Turn a settled conversation or decisions file into .scratch/<feature>/spec.md and one issues/NN-slug.md per task, in one pass, with no new interview. For code and non-code work alike.
disable-model-invocation: true
---

# Write spec

Write work that spans sessions as a spec and vertical-slice issues that a
fresh session can start from without asking anything. Synthesize what is
already settled; the interview happened before this skill.

## 1. Gather

Read the conversation and every input file the user names. List:

- the goal, in two to four sentences;
- each settled decision, with its reason;
- what was explicitly put out of scope or deferred;
- open questions the slicing depends on;
- questions still too vague to ask, for `## Not yet specified`.

An open question gets no interview here. It becomes its own task ("Decide
X", with the options in its issue) or goes to Out of scope, and the quiz in
step 3 flags it.

Done when every task you are about to propose traces to a settled decision
or to a "Decide X" task.

## 2. Check what already exists

Search the repo, `.scratch/`, `docs/adr/` and any glossary by concept, not
by the request's wording: the same thing may exist under another name. Note
where you looked. Something that already exists becomes a change to it, or
drops out; a prior rejection goes to Out of scope with its reason.

Done when you can name, for each task, the existing file or thing it builds
on, or say that none exists and where you looked.

## 3. Slice and quiz

Cut the work into vertical slices: each one delivers an outcome that can be
checked on its own (a command that runs, a query that returns the right
rows, a document a reader can act on, a host that answers), not one layer
of a later outcome. Put any prefactor (restructuring that makes the rest
easy) first. Size each slice for one fresh session.

Show the slicing as a numbered list, one line per issue: title, `Blocked
by`, what it delivers, its done-when. Then ask, numbered, each with your
recommendation: granularity right, blocking edges right, anything to merge
or split, the feature folder name, and each flagged open question. Word
them so "yes" accepts the recommendation.

Done when the user approves the list. Rework and re-show until then. If the
user asked for a preview, show the spec and issues as they would be written
and stop here.

## 4. Write

Write `.scratch/<feature>/`. If the folder exists, add to it: number new
issues after the highest existing one and append rows to the task table.

`spec.md`, in this order:

- `## Goal`, plus the settled input files by path and which wins on conflict.
- Context sections for settled material that more than one issue leans on,
  each named for what it holds. `## Principles` carries every cross-cutting
  rule from the input that each issue must obey; an inventory or a layout
  gets its own section. No user stories.
- `## Tasks`: a table `| # | Task | Blocked by | Status |`, one row per issue.
- `## Not yet specified`, only when step 1 listed any: one line per question.
  Each stays toward the goal; anything that won't be built is Out of scope.
- `## Out of scope`: each deferred or rejected item, with where it went (a
  later task, another spec, dropped).
- `## Decisions`: empty, except a Decision and Reason table for answers the
  user approved in the step 3 quiz. Decisions from the input stay in the
  input files the Goal names.

`issues/NN-slug.md`, numbered from `01` in dependency order:

```
# NN Title

Status: open
Blocked by: NN, NN

## What

## Done when
```

Omit the `Blocked by` line when nothing blocks the issue. Under What, name
the paths, tables, hosts or documents to touch, the inputs to read, and what
may not change. The done-when is a check a fresh session can run and quote:
a command and its expected output, a query result, a file with named
sections, or a named reviewer's sign-off for a document. Name the check
only if it exists; otherwise the issue creates it first.

Done when `scripts/check-spec.sh .scratch/<feature>` in this skill's folder
prints that spec and issues agree. Fix what it reports and rerun.

## 5. Publish

Read the `tracker` source: the repo's `## Sources`, else the machine's. With
`local` or no line, skip this step and say so once. Otherwise show the issues
and, on approval, run its create command for each in number order, carrying
`Blocked by` across as the line allows, and add `Tracker: <url or key>` under
`Status:`.

Done when each issue has a Tracker line, or the user declined.

## 6. Hand off

Reread issue 01 as a session that has seen nothing but the repo. Fix
anything it would have to ask. Then give the user the prompt to start it:
"Work issue 01 in .scratch/<feature>/issues/01-slug.md. Settled input:
<files>."

Done when issue 01 names every file and fact it needs.
