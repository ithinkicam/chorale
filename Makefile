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

.PHONY: go-version build test e2e lint golden check

go-version:
	@go version

build:
	go build -o bin/chorale ./cmd/chorale

# Unit, interface snapshot and interaction tests. The JSON is read by the
# scenario check.
test:
	@mkdir -p bin
	$(GOTESTSUM) --jsonfile bin/test.json -- $(UNIT_PKGS)

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

# Rewrites golden files. Use on purpose only, and read every change in the diff.
golden:
	@./scripts/update-golden.sh

check: lint test e2e
