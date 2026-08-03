#!/usr/bin/env bash
# name=List info
# desc=Show type, size, and permissions for each selected path
# path=Files
set -euo pipefail

if [ "$#" -eq 0 ]; then
	echo "no paths selected"
	exit 1
fi

for p in "$@"; do
	if [ -d "$p" ]; then
		kind="dir "
	elif [ -f "$p" ]; then
		kind="file"
	else
		kind="?   "
	fi
	# ls -ld gives one line per path (the dir itself, not its contents) with perms + size.
	printf '%s  %s\n' "$kind" "$(ls -ldh "$p" 2>/dev/null || echo "$p")"
done
