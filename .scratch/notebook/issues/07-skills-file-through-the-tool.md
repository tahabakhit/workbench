# 07 Skills file through the tool

Status: open
Blocked by: 03, 05

## What

Move filing and notebook searching to the CLI at
`~/.local/share/workbench/bin/notebook`.

- `skills/record/SKILL.md`: file entries with `notebook file` and search with
  `notebook search`, in place of direct writes, commits and reading files.
- `skills/plan-day/SKILL.md` and `skills/study-topic/SKILL.md`: same for plans
  and topic files.
- `skills/improve/SKILL.md`, `skills/settle/SKILL.md`,
  `skills/prototype/SKILL.md` and `skills/onboard/SKILL.md`: change only where
  they write to the notebook. `glossary-adr-format.md` copies stay identical.
- Each skill keeps the fallback: with no binary, write the file directly and
  commit it, as now.
- `contract.md`: drop "will search and file into it; until it does", add
  `bin/` to what lives under `~/.local/share/workbench/`, and name the commands.
- `GLOSSARY.md`: add the tool if a skill needs the term. `AGENTS.md` and
  `README.md`: add what a reader needs to know about `notebook/`.

`study-topic` must still never overwrite a file in a topic folder. Keep the
contract under the hook's cap.

## Done when

- `sh scripts/check.sh` passes.
- `rg -n "notebook file|notebook search" skills` shows `record`, `plan-day` and
  `study-topic`.
- A real run of `record`, `plan-day` and `study-topic` each adds a commit to
  the vault made by `notebook file`; `git -C
  ~/.local/share/workbench/notebook log --oneline -3` is quoted here.
- With the binary renamed away, a `record` run still files directly.

## Notes from issue 03

- `notebook file --update` replaces the whole file: read it, add the new
  text, pipe the result in. That is how `study-topic` adds to `quiz.md`.
- Every note needs frontmatter with `date`, so `study-topic` must write it, and
  the five files under `learning/posix-sh-parameter-expansion/` fail
  `notebook validate` until they gain it.
