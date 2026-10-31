# Tracker lines

Seed lines for the `tracker` source. Fill in the project, then replace the
line's commands with what the user gives for anything else. These are the only
place a tracker's commands are written down; other skills run what the line
says.

```md
- tracker: local, .scratch/ only
- tracker: GitHub <owner/repo> | list `gh issue list --repo <owner/repo> --state open --limit 50` | create `gh issue create --repo <owner/repo> --title <title> --body-file <file>` | close `gh issue close <n> --repo <owner/repo> --comment <evidence>`
- tracker: GitHub <owner/repo>, project <n> | list `gh issue list --repo <owner/repo> --state open --limit 50` | create `gh issue create --repo <owner/repo> --title <title> --body-file <file>`, then `gh project item-add <n> --owner <owner> --url <issue url>` | close `gh issue close <n> --repo <owner/repo> --comment <evidence>`
- tracker: GitLab <group/repo> | list `glab issue list -O json` | create `glab issue create --title <title> --description <text>` | close `glab issue note <n> --message <evidence>`, then `glab issue close <n>`
- tracker: <name> | list <command or tool call> | create <command or tool call> | close <command or tool call>
```

A GitHub Project needs `gh auth refresh --hostname github.com --scopes project`
once. Its board sets new items to Todo and closed ones to Done by itself.

`glab` works out the project from the clone's remote, so run its commands
inside the repo.

A line with a create command names one repo, so it belongs to that repo: in
its `## Sources`, or keyed `tracker (<repo>)` when the repo is shared. The
machine's plain `tracker` line gathers work across repos and can stay
search-only; `write-spec` cannot publish from it.
