---
name: prototype
description: Use when a question about business logic, a state model or a data shape is better answered by clicking through something than by discussing it. Builds a throwaway one-file HTML demo for anyone to try.
---

# Prototype

Throwaway code that answers one question. The people clicking it may not
read code, so every label is in the domain's words.

## 1. Question

Write one paragraph: the question, the states, and the actions that move
between them, in domain language. Ask the user to confirm it; if they are
away, state it as an assumption at the top of the demo.

Done when the question is written and confirmed or marked as assumed.

## 2. Build

One HTML file with everything inline, no framework and no server, so it
opens by double-click and survives being emailed. Save it as
`.scratch/<feature>/prototype-<slug>.html` when the work has a feature
folder, otherwise as `${TMPDIR:-/tmp}/prototype-<slug>.html`.

- The logic is a pure module in its own `<script>`: a reducer, a state
  machine or a few pure functions, whichever fits the question. It never
  touches the DOM, so it can be lifted into real code later.
- Top to bottom: the title and the question; the current state as labelled
  fields, re-rendered after every action with what changed marked; one
  button per action; then one tab per scenario. Each scenario resets to a
  fixed start and lists its steps as buttons in order.
- Scenarios cover the happy path, the awkward edge case, and an attempt at
  something that should be impossible.
- State lives in memory. No tests, no persistence, one accent colour, no
  animation.

When `node` is on `PATH`, load the module script into it and run each
scenario's actions, printing the end state.

Done when the file exists, and each scenario's end state is printed or,
without `node`, listed with the state you expect.

## 3. Hand over and record

Give the path. The useful reaction is "that shouldn't be possible" or "I
expected X": a bug in the idea. Add actions or scenarios when asked.

When the question is answered, propose the question and verdict as a
decision: to `## Decisions` in `spec.md` when a feature folder exists,
otherwise to `decisions/` in the notebook. Write it only on approval. Lift
the module into real code only when asked; the HTML page stays a prototype.

Done when the verdict is written on approval, or the user has said the
prototype answered nothing worth keeping.
