#!/usr/bin/env bash
# Fails unless the installed golangci-lint is the version the Makefile names.
set -euo pipefail

want="$1"
have="$(golangci-lint version --short 2>/dev/null || true)"
if [[ "$have" != "$want" ]]; then
  echo "golangci-lint $want is required, but found '${have:-none}'."
  exit 1
fi
