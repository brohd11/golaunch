#!/usr/bin/env bash
# name=Shell here
# desc=Open an interactive shell listing the selected paths (runs in an external terminal)
# path=Files
# terminal=true
set -euo pipefail

echo "golaunch selection:"
for p in "$@"; do
	echo "  $p"
done
echo
echo "You are in: $(pwd)"
echo "Press enter to close…"
# terminal=true launches this in its own window, so it can block on input the TUI could not give.
read -r _
