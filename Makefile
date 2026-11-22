# Chorale's development commands. See docs/testing.md.

SHELL := /bin/bash

# go.mod's toolchain line is only a minimum, so pin it for real here.
GO_TOOLCHAIN := $(shell sed -n 's/^toolchain //p' go.mod)
export GOTOOLCHAIN := $(GO_TOOLCHAIN)
export OPENSPEC_TELEMETRY := 0

# The one golangci-lint version Chorale is checked with, here and in CI.
GOLANGCI_LINT_VERSION := 2.15.0

# Every package except the terminal-level tests, which make e2e runs.
UNIT_PKGS := $(shell go list ./... | grep -v '/e2e$$')

GOTESTSUM := go tool gotestsum --format testname

.PHONY: go-version build test e2e lint golden scenarios specs check

go-version:
	@go version

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

build:
	go build -ldflags "-X github.com/ithinkicam/chorale/internal/version.Version=$(VERSION)" \
	  -o bin/chorale ./cmd/chorale

# Unit, interface snapshot and interaction tests. The JSON is read by the
# scenario check. -count=1 because some tests build programs from source,
# and Go's test cache can't see changes to what they build.
test:
	@mkdir -p bin
	$(GOTESTSUM) --jsonfile bin/test.json -- -count=1 $(UNIT_PKGS)

# Terminal-level tests. -count=1 because they build chorale themselves, and
# Go's test cache can't see changes to the code they build.
e2e:
	@mkdir -p bin
	$(GOTESTSUM) --jsonfile bin/e2e.json -- -count=1 ./e2e/...

lint:
	@./scripts/check-lint-version.sh $(GOLANGCI_LINT_VERSION)
	@./scripts/check-gofmt.sh
	go vet ./...
	golangci-lint run ./...
	go tool actionlint

# Rewrites golden files. Use on purpose only, and read every change in the diff.
golden:
	@./scripts/update-golden.sh

# Every OpenSpec scenario must have a matching test that passed in the
# results make test and make e2e just wrote.
scenarios: test e2e
	go run ./scripts/check-scenarios bin/test.json bin/e2e.json

# OpenSpec comes from package-lock.json, so it is the pinned version.
# --all matters: without it, validate fails when there is no terminal.
node_modules/.package-lock.json: package-lock.json
	npm ci --no-fund --no-audit

specs: node_modules/.package-lock.json
	npx --no-install openspec validate --strict --all

check: lint test e2e scenarios specs
