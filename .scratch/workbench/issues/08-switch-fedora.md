# 08 Switch the Fedora PC to the plugin

Status: done
Blocked by: 07

## What

Run on the Fedora PC. Clone this repo, add it as a marketplace, install the
plugin with auto-update, include `contract.md` in that machine's global
instructions (the machine's own policy states that a repository's
`AGENTS.md` adds to the global rules and wins on conflict), and link `~/.agents/skills` as on the Mac. Find the
`rhel-fedora` skill and move it into the Fedora machine's own setup (not this
plugin). Run `onboard` machine scope to write the Fedora `## Knowledge`
section.

## Done when

- On Fedora: `claude plugin details workbench` lists all ten skills, the
  global instructions contain the contract's first heading, and
  `rhel-fedora` is reachable from the Fedora setup and absent from this repo.
