# 04 Move and rename the seven one-to-one skills

Status: open
Blocked by: 01

## What

Copy seven skills from `~/dev/agents/home/dot_agents/skills/` into `skills/`
under their new names. Mechanical change only; issue 06 changes content.

| From | To |
|---|---|
| diagnose | diagnose |
| an-metric-diff | explain-gap |
| grill | settle |
| improve | improve |
| questionnaire | draft-questions |
| triage | plan-day |
| writing-skills | write-skill |

`rephrase` came in 01. For each:

- rename the folder, `name:`, heading and `agents/openai.yaml`;
- chezmoi file prefixes become real names with the executable bit
  (`scripts/executable_x.sh` to `scripts/x.sh`, `chmod +x`);
- replace every old skill name in all skill texts (`rg -n
  'an-metric-diff|grill|questionnaire|repitch|triage|writing-skills|setup-repo|retro|harvest' skills/`);
- remove skill-to-skill loading: `improve` reads the "cheapest form" table
  from `writing-skills`; copy that table into a side file in both folders.

Leave `~/dev/agents` unchanged.

## Done when

- `sh scripts/check.sh` exits 0.
- `claude plugin details workbench` lists eight skills: these seven plus `rephrase`.
- The `rg` command above prints only intended mentions (quote it), and
  `find skills -name 'executable_*'` prints nothing.
