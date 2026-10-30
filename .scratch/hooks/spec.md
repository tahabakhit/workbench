# hooks

## Goal

The plugin ships no Bash guard. Git policy is text in `agent-rules.md`, and
gitleaks runs on every push through a global git pre-push hook the dotfiles
ship. The tasks below built a `PreToolUse` guard that the last row of
Decisions removes; the rest of this spec describes that guard.

Settled input: `.scratch/bootstrap/decisions.md` (wins on conflict) and the
`## Decisions` below.

## Principles

- Influence, not force. The hook replies "ask" so the owner approves in
  place. It replies "deny" only when gitleaks is missing.
  Any command it doesn't recognise passes with no reply.
- Narrow matchers. Plugin hooks run in every session and can't be turned
  off one by one. `hooks.json` keeps a single `PreToolUse` entry on `Bash`,
  and its `if` matchers start the script only for `git push`, `git commit`,
  `git reset`, `git clean`, `git checkout`, `git restore` and `rm`.
- One script, `hooks/guard-bash.sh`, POSIX `sh` with `jq`, executable. It
  reads the tool input JSON on stdin and checks `tool_input.command`. Each
  row in Decisions is one function in it.
- Matching on the command text is best effort. A miss or a needless "ask"
  is acceptable. A false "deny" is a bug.
- Tests first. `hooks/test-safety.sh` pipes sample tool-input JSON into the
  script and checks the reply, using throwaway git repos under `mktemp -d`.
  `scripts/check.sh` runs it.
- The reply format comes from the installed Claude Code's hook docs, checked
  for `permissionDecision` (`ask`, `deny`) and `permissionDecisionReason`
  under `hookSpecificOutput`, not from memory. `claude plugin validate
  --strict` doesn't check hook fields, so `if` and `ask` are proved in a live
  session with `claude --plugin-dir ~/dev/workbench`.

## Tasks

| # | Task | Blocked by | Status |
|---|---|---|---|
| 01 | Safety hook: ask on force push and remote branch delete | | done |
| 02 | Gitleaks on commit and push | 01 | done |
| 03 | Ask before destructive git on a dirty tree | 01 | done |
| 04 | Ask before `rm` outside the working areas | 01 | done |

## Out of scope

- Permission rules in `~/.claude/settings.json`. No repo manages that file;
  the owner adds rules by hand on each machine.
- `gh` and other outward-facing actions (merge, release, repo delete,
  workflow runs). They stay with the agent rules and the permission mode.
- `git branch -D`, because the reflog recovers it.
- `Stop` or end-of-task checks. The closing check stays an instruction line
  until it is seen being skipped.
- Repo-specific guards such as the `claude.md`/`agents.md` filename rule.
  They go in that repo's own settings.
- New safety lines in `contract.md`. The hook enforces what
  `agent-rules.md` already says.
- Hooks for Codex and Pi.
- Commits or pushes made from inside scripts or aliases.

## Decisions

| Decision | Reason |
|---|---|
| `git push` with `--force`, `-f`, `--force-with-lease` or a `+refspec`, on any branch: ask | Rewrites remote history; the agent rules forbid force push without explicit approval |
| `git push --delete` or `git push <remote> :<branch>`: ask | Deletes a remote branch |
| `git commit`: run `gitleaks git --staged --redact`; ask on a finding, naming file and line with the secret redacted | A committed secret that gets pushed can't be taken back, but test fixtures and false positives must stay approvable in place |
| `git push`: run gitleaks over `@{u}..HEAD`, or the branch's commits not on any remote when it has no upstream; ask on a finding | Same reason |
| Gitleaks not installed when committing or pushing: deny with a "gitleaks missing" message | Matches the agent rule "do not commit or push; tell the owner" |
| `git reset --hard`, `git clean -f*`, `git checkout -- .`, `git restore .`: ask only when `git status --porcelain` is non-empty | Destroys uncommitted work; on a clean tree there is nothing to lose |
| `rm` with `-r` or `-f` and a target that resolves outside the project dir, the session scratchpad and `/tmp`: ask | Deletes outside the repo, one of the irreversible cases in the bootstrap decisions |
| No hooks for `gh` actions, `git branch -D`, `Stop` checks or repo-specific rules, and no new contract lines | Influence, not force; each is reviewable, recoverable or covered elsewhere |
| One hook, one script, a single `PreToolUse` entry in `hooks.json`; tests in `hooks/test-safety.sh`, run by `scripts/check.sh` | Narrow matchers and one place to read the rules |
| The guard hook, its tests and its `hooks.json` entry are removed | The owner keeps git rules as text in `agent-rules.md` and adds permission rules by hand per machine; gitleaks on push is the dotfiles' global pre-push hook |
