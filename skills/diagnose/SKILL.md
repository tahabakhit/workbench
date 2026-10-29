---
name: diagnose
description: Use this before investigating any bug, error, crash, failing command, failing test, regression or slowdown in code, scripts, builds, tools or setup. Enforces reproduce-first debugging. Not for two numbers that disagree; that is explain-gap.
---

# Diagnose

## 1. Reproduce

Find one command that shows this exact failure, and run it. Narrow it until it
is fast and fails every time (or at a known rate for flaky failures).

Reading code to build a theory before this command exists is guessing.

Done when the command has been run and its output shows the failure.

## 2. Hypothesise

List three to five causes, ranked. State each as a prediction that can be
disproved: "if it is the cache, clearing it makes the command pass".

Done when each hypothesis has a test that could prove it wrong.

## 3. Test one at a time

Run the cheapest test of the top hypothesis. Change one thing per test. Tag
temporary logging with `[diag]` so it can be found and removed.

Done when one hypothesis survives and the others are ruled out by output.

## 4. Fix and confirm

Fix the cause, not the symptom. Rerun the command from step 1; it now passes.
Remove every `[diag]` line. If a test or check could catch this next time, say
so.

Done when the step 1 command passes and its output is quoted.
