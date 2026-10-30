#!/bin/sh
# SessionStart hook: plain stdout becomes session context, capped at 10,000
# characters, so contract.md must stay under that.
root=$(cd "$(dirname "$0")/.." && pwd)
cat "$root/contract.md"
