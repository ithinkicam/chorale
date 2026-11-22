#!/usr/bin/env bash
# Rewrites golden files by running, with -update, only the packages whose
# tests compare against them. Use on purpose only, and read every change in
# the diff.
set -euo pipefail

pkgs=()
while IFS= read -r dir; do
  pkgs+=("./$dir")
done < <(git grep -l --untracked -e 'RequireGolden(' -- '*_test.go' | xargs -n1 dirname | sort -u)

if [[ ${#pkgs[@]} -eq 0 ]]; then
  echo "No tests use golden files yet."
  exit 0
fi

go test -count=1 "${pkgs[@]}" -update
