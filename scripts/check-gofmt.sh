#!/usr/bin/env bash
# Fails if any tracked or new Go file is not gofmt-formatted.
set -euo pipefail

files=()
while IFS= read -r f; do files+=("$f"); done < <(git ls-files --cached --others --exclude-standard '*.go')
[[ ${#files[@]} -eq 0 ]] && exit 0

unformatted="$(gofmt -l "${files[@]}")"
if [[ -n "$unformatted" ]]; then
  echo "These files are not gofmt-formatted:"
  echo "$unformatted"
  exit 1
fi
