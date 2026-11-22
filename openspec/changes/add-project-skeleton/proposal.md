# Proposal: add the project skeleton

## Why

Chorale is test first, including the interface, so the test harness has to exist before any real feature does. This change sets up the Go project, every test layer in `docs/testing.md`, `make check` and continuous integration. It proves them on the smallest real screen Chorale has: the one shown when there are no sessions yet.

This is milestone 0 in `docs/plan.md`.

## What Changes

- A Go module at `github.com/ithinkicam/chorale` with a `chorale` command.
- `chorale` opens a full-screen terminal interface that shows an empty session list and quits with `q`.
- `chorale --version` prints the version without opening the interface.
- Below 60 columns, Chorale says it needs a wider window instead of drawing a broken layout.
- Test harness for unit tests, interface snapshot tests, interaction tests and terminal-level tests, with golden files under `testdata/`.
- `scripts/check-scenarios`, which fails if any OpenSpec scenario has no matching test.
- A `Makefile` with `test`, `e2e`, `lint`, `golden` and `check` targets.
- GitHub Actions running `make check` on every push.

## Capabilities

### New Capabilities

- `app-shell`: starting and quitting the terminal app, reporting its version, the empty state when there are no sessions, and how it behaves when the window is too narrow.

### Modified Capabilities

None.

## Out of scope

Themes beyond the terminal's own palette, the background service, agents, the prompt editor, sprites and dictation. Those belong to milestones 1 to 4. The end-to-end layer with a scripted agent needs the service, so it arrives in milestone 1. This change only creates the empty `make e2e` target and the terminal-level tests.

## Impact

New repository content only: `go.mod`, `cmd/chorale/`, `internal/`, `scripts/`, `testdata/` directories, `Makefile`, `.golangci.yml` and `.github/workflows/check.yml`. New dependencies: Bubble Tea v2, Lip Gloss v2, Charm's `teatest` v2, a pseudo-terminal library and a terminal emulator library for tests.
