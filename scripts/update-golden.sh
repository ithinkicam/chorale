#!/usr/bin/env bash
# Rewrites golden files by running only the packages that have them with
# -update. Use on purpose only, and read every change in the diff.
set -euo pipefail

pkgs=()
while IFS= read -r dir; do
  pkgs+=("./$dir")
done < <(git ls-files --cached --others --exclude-standard '*.golden' | sed -E 's|/testdata/.*$||' | sort -u)

if [[ ${#pkgs[@]} -eq 0 ]]; then
  echo "No golden files yet."
  exit 0
fi

go test -count=1 "${pkgs[@]}" -update
