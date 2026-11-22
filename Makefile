# Chorale's development commands. See docs/testing.md.

SHELL := /bin/bash

# go.mod's toolchain line is only a minimum, so pin it for real here.
GO_TOOLCHAIN := $(shell sed -n 's/^toolchain //p' go.mod)
export GOTOOLCHAIN := $(GO_TOOLCHAIN)

.PHONY: go-version
go-version:
	@go version

.PHONY: build
build:
	go build -o bin/chorale ./cmd/chorale
