#!/bin/sh
# SessionStart hook: installs the notebook binary in the background, because a
# download can outlast the hook's timeout.
root=$(cd "$(dirname "$0")/.." && pwd)
nohup sh "$root/notebook/run.sh" --install >/dev/null 2>&1 &
