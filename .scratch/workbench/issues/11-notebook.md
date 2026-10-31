# 11 notebook: build the vault tool in this repo, up to switch-over

Status: open
Tracker: https://github.com/tahabakhit/workbench/issues/2

Step 4 of the Order. The notebook tool lives in `notebook/` in this repo and
ships in this plugin. Its own spec goes in `.scratch/notebook/`, written with
`write-spec` from the notebook section of `.scratch/bootstrap/decisions.md`.
Skills already write the notebook layout in `contract.md` directly; the tool
indexes, searches and files into those same files.

## Done when

- `.scratch/notebook/spec.md` is done.
- The plugin ships the tool, and `record`, `plan-day` and `study-topic` file
  through it into `~/.local/share/workbench/notebook` on this PC.
- Daily use on this PC runs on the notebook, with no other vault. Other
  machines are dropped; see `## Decisions` in `.scratch/notebook/spec.md`.
