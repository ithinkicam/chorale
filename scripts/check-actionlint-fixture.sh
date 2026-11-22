#!/usr/bin/env bash
# Fails unless actionlint finds problems in the known-broken workflow fixture.
set -uo pipefail

fixture=testdata/actionlint/broken.yml
if [[ ! -f "$fixture" ]]; then
  echo "$fixture is missing."
  exit 1
fi
go tool actionlint "$fixture" >/dev/null 2>&1
status=$?
if [[ $status -ne 1 ]]; then
  echo "actionlint exited with status $status on $fixture, want 1, meaning it found problems."
  exit 1
fi
