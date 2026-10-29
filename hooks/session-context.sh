#!/bin/sh
# SessionStart hook: plain stdout becomes session context, capped at 10,000
# characters, so the contract plus the knowledge file must stay under that.
root=$(cd "$(dirname "$0")/.." && pwd)
knowledge="${XDG_CONFIG_HOME:-$HOME/.config}/workbench/knowledge.md"

cat "$root/contract.md"
echo
if [ -r "$knowledge" ]; then
  echo "# This machine's Knowledge destinations ($knowledge)"
  echo
  cat "$knowledge"
else
  echo "# This machine's Knowledge destinations"
  echo
  echo "None: $knowledge does not exist. Run the onboard skill, machine scope, to write it."
fi
