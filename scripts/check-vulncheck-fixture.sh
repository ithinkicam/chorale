#!/usr/bin/env bash
# Fails unless govulncheck finds the planted vulnerability in its fixture.
set -uo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
fixture="$root/testdata/vulncheck"
db="$fixture/db"
if [[ ! -f "$fixture/go.mod" || ! -f "$db/index/db.json" ]]; then
  echo "The govulncheck fixture in $fixture is missing."
  exit 1
fi
cd "$root" || exit 1
go tool govulncheck -C "$fixture" -db "file://$db" ./... >/dev/null 2>&1
status=$?
if [[ $status -ne 3 ]]; then
  echo "govulncheck exited with status $status on $fixture, want 3, meaning it found a vulnerability."
  exit 1
fi
