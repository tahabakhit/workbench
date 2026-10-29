# 17 Install from GitHub on every machine

Status: done

The README gives the install steps from the `tahabakhit/workbench` GitHub
marketplace. The Fedora PC switches from the `directory` source at
`~/dev/workbench` to that GitHub source, so sessions get only pushed commits.

## Done when

- The README's Install section uses `claude plugin marketplace add
  tahabakhit/workbench`.
- `claude plugin marketplace list` on the PC shows `workbench` with a GitHub
  source, and `claude plugin list` shows `workbench@workbench` enabled.
- `sh scripts/check.sh` passes.

## Out of scope

- The Mac, which waits until the repo is public (issue 15).
- Making the repo public.
