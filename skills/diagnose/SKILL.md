---
name: diagnose
description: Use this before investigating any bug, error, crash, failing command, failing test, regression or slowdown in code, scripts, builds, tools or setup. Enforces reproduce-first debugging. Not for two numbers that disagree; that is explain-gap.
---

# Diagnose

Write `<REDACTED>` in place of every secret in the commands, output and notes
you show. Keep credentials in environment variables, not in the command.

## 1. Reproduce

Find one command that shows this exact failure, and run it. Narrow it until it
is fast and fails every time (or at a known rate for flaky failures).

Reading code to build a theory before this command exists is guessing. If no
such command can be built, stop and say what you tried. Ask for access to an
environment that fails, a redacted log or artifact, or approval for temporary
logging.

Then minimise: cut inputs, config, data and steps one at a time, rerunning
after each cut, and keep only what the failure needs.

Done when the command has been run, its output shows the failure, and
removing any remaining part makes it pass.

## 2. Hypothesise

List three to five causes, ranked. State each as a prediction that can be
disproved: "if it is the cache, clearing it makes the command pass". Show the
ranked list to the user before testing; they may know one is already ruled
out. Keep testing while you wait for a reply.

Done when each hypothesis has a test that could prove it wrong.

## 3. Test one at a time

Run the cheapest test of the top hypothesis. Change one thing per test. Tag
temporary logging with `[diag]` so it can be found and removed.

Done when one hypothesis survives and the others are ruled out by output.

## 4. Fix and confirm

Fix the cause, not the symptom. Rerun the command from step 1; it now passes.
Remove every `[diag]` line. If a test or check could catch this next time, say
so. The commit message names the cause that the surviving hypothesis found.

Done when the step 1 command passes and its output is quoted.
