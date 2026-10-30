# 11 notebook: build the vault tool in this repo, up to switch-over

Status: open

Step 4 of the Order. The notebook tool lives in `notebook/` in this repo and
ships in this plugin. Its own spec goes in `.scratch/notebook/`, written with
`write-spec` from the notebook section of `.scratch/bootstrap/decisions.md`.
Skills already write the notebook layout in `contract.md` directly; the tool
indexes, searches and files into those same files.

## Done when

- `.scratch/notebook/spec.md` is done.
- The plugin ships the tool, and `record`, `plan-day` and `study-topic` file
  through it into `~/notebook` on this PC.
- Daily use on every machine runs on the notebook, with no other vault.
