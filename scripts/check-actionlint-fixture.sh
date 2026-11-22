#!/usr/bin/env bash
# Fails if actionlint stops rejecting the known-broken workflow fixture.
set -uo pipefail

if go tool actionlint testdata/actionlint/broken.yml >/dev/null 2>&1; then
  echo "actionlint accepted testdata/actionlint/broken.yml, which it must reject."
  exit 1
fi
