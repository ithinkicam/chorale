# Proposal: add the project skeleton

## Why

Chorale is test first, including the interface, so the test harness has to exist before any real feature does. This change sets up the Go project, test layers 1 to 4 in `docs/testing.md`, `make check` and continuous integration. It proves them on the smallest real screen Chorale has: the one shown when there are no sessions yet.

This is milestone 0 in `docs/plan.md`.

## What Changes

- A Go module at `github.com/ithinkicam/chorale` with a `chorale` command.
- `chorale` opens a full-screen terminal interface that shows an empty session list and quits with `q`.
- `chorale --version` prints the version without opening the interface.
- In a window smaller than 60 by 20, Chorale says the size it needs instead of drawing a broken layout.
- Test harness for unit tests, interface snapshot tests, interaction tests and terminal-level tests. Tests assert each scenario directly, and plain-text golden files record each screen.
- `scripts/check-scenarios`, which fails if any OpenSpec scenario has no matching test that passed.
- A `Makefile` with `build`, `test`, `e2e`, `lint`, `golden` and `check` targets, and OpenSpec pinned through `package.json` and its lockfile.
- GitHub Actions running `make check` on pushes to `main` and on every pull request, with actions pinned to commit hashes and a read-only token.
- `actionlint` in `make lint`, so mistakes in workflow files fail the gate.
- `docs/style.md`, the code style every agent follows. Code explains itself, and the reasons live in specs and design documents rather than in comments.

## Capabilities

### New Capabilities

- `app-shell`: starting and quitting the terminal app, reporting its version, the empty state when there are no sessions, and how it behaves when the window is too small.

### Modified Capabilities

None.

## Out of scope

Themes beyond the default one, the background service, agents, the prompt editor, sprites and dictation. Those belong to milestones 1 to 4. The end-to-end layer with a scripted agent needs the service, so it arrives in milestone 1. In this change, `make e2e` runs the terminal-level tests only. The theme contrast check arrives with the first fixed-colour theme in milestone 1, and the speech and phone test layer with `chorale-listen` in milestone 2.

Repository hardening is a separate change, `harden-repository`, which follows this one because it needs continuous integration to exist. It covers security scanning, Dependabot updates and auto-merge, the required checks on `main`, `CODEOWNERS` and `SECURITY.md`.

## Impact

New repository content only: `go.mod`, the `cmd/chorale/`, `internal/`, `e2e/`, `scripts/` and `testdata/` directories, `Makefile`, `.golangci.yml`, `package.json`, `package-lock.json`, `.nvmrc` and `.github/workflows/check.yml`, plus one rule in `.gitignore`. New dependencies: Bubble Tea v2, Lip Gloss v2, Charm's `teatest` v2, a pseudo-terminal library and a terminal emulator library for tests. The development tools are OpenSpec, through npm, and gotestsum and actionlint, as Go tools pinned in `go.mod`.
