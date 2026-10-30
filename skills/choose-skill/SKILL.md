---
name: choose-skill
description: Say which workbench skill or flow fits the user's situation, or what to do with the session when a phase of work ends. Run with the situation in a sentence.
disable-model-invocation: true
---

# Choose skill

The user doesn't remember every skill, so they ask. Answer with one
recommendation: the skill to run next, its argument, and one line on why.
Name a runner-up only when the choice is close.

Before saying what a skill does or that a step can be skipped, read its
`SKILL.md` in this plugin's skills folder, one level above this skill's
folder. The map below is for orientation only.

## The main flow: idea to done

1. **`settle`** sharpens a plan or decision by interview until every
   question has an answer.
2. When a question needs more than talk:
   - **`prototype`** when it is about logic or a state model you have to
     click through;
   - **`research`** when it needs facts from sources;
   - **`draft-questions`** when the answer is in someone else's head.

   Take what comes back into `settle`.
3. Multi-session work: **`write-spec`** turns the settled thread into
   `.scratch/<feature>/` with one issue per task. Work the issues in order of
   their `Blocked by` lines, each to its done-when as the contract says.
   Work that fits this session needs no spec; do it here.
4. **`record`** files the finished work, decisions, lessons and corrections
   into the notebook.

Keep steps 1 to 3 in one session where the context allows, so the spec is
written from the interview itself, not a summary of it.

## On-ramps

- **Start of a day or of a session with no set task**: `plan-day`.
- **Something is broken** (an error, a failing test, a slowdown):
  `diagnose`. It runs on its own when a bug comes up.
- **Two numbers that should match don't**: `explain-gap`. Also runs on its
  own.

## Upkeep

- **The same correction keeps coming back**: `improve`, after `record` has
  filed corrections from a few sessions.
- **New machine or new repo**: `onboard`, machine scope first.
- **Writing or reviewing a skill or an `AGENTS.md`**: `write-skill`.

## Standalone

- **Learning a topic over several sessions**: `study-topic`.
- **The last message didn't land**: `rephrase`.
- **Work moves to another harness, machine or person**: `handoff`.

## Phase boundaries

When the situation is that a chunk of work just ended and the question is
what to do with the session, read `phase-boundaries.md` in this skill's
folder and walk its tree. Mid-phase there is no such decision: continue, or
give the rest to a subagent.

Done when the user has one recommendation they can run as typed.
