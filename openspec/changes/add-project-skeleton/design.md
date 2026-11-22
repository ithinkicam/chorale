# Design

## Context

Nothing exists yet but the plan, the testing standard and OpenSpec itself. See `proposal.md` for why the skeleton comes first. The stack was chosen in `docs/plan.md`: Go, with Bubble Tea v2 and Lip Gloss v2 for the interface.

## Goals / Non-Goals

**Goals:**

- Every test layer in `docs/testing.md` that doesn't need the service exists and has at least one real test.
- `make check` is the single gate, the same locally and in continuous integration.
- The empty session list is built test first, as the worked example for every later change.

**Non-Goals:**

- No service, agents, themes beyond the terminal palette, or prompt editor.
- No release packaging or Homebrew formula.

## Decisions

### Go version and module

The module is `github.com/ithinkicam/chorale`. It targets the Go version Bubble Tea v2 needs. Petal builds with Go 1.27.3, so the floor is 1.27, pinned in `go.mod` with a `toolchain` line, so every machine and continuous integration use the same one.

### Layout

```
cmd/chorale/          the command: flag parsing, then hands off to internal/app
internal/app/         the Bubble Tea model for the whole app
internal/screens/     one package per screen; this change adds sessionlist
internal/termtest/    the pseudo-terminal test harness
scripts/              check-scenarios and other development scripts
testdata/             shared golden files and fixtures
```

Each screen package has its own `testdata/` for its golden files.

### Interface tests use teatest v2

Snapshot and interaction tests use `github.com/charmbracelet/x/exp/teatest/v2`. It runs the real model, sends it messages, including window sizes, keys and mouse events, and compares the final output to a golden file. It is experimental and published as a pseudo-version, so the exact version is pinned in `go.mod`. The alternative, calling the model's `View` directly, would skip the program loop and miss bugs in how messages are handled.

### Terminal-level tests use a real pseudo-terminal

`internal/termtest` starts the built `chorale` binary in a pseudo-terminal with `creack/pty`, which Petal also uses. It sends raw bytes and feeds the output to a terminal emulator library to read back the screen and the cursor position. The emulator library is chosen by the first task in group 2: candidates are `github.com/charmbracelet/x/vt` and `github.com/hinshun/vt10x`. The test that picks it is the "Quitting restores the terminal" scenario, because that needs alternate screen, cursor visibility and mouse mode to be tracked correctly.

### The scenario check

`scripts/check-scenarios` is a small Go program, run with `go run`. It reads every `#### Scenario:` heading in `openspec/specs/` and `openspec/changes/`, skipping `openspec/changes/archive/`, turns each name into a slug, and looks for a Go test whose name or golden file matches. It prints each scenario with no test and exits non-zero if there are any. A scenario marked "Checked by hand" in its spec is listed but doesn't fail.

### Archived changes are not tracked

`openspec/changes/archive/` is ignored in `.gitignore`, except for its `.gitkeep`, so archived changes don't pile up in the repo. Nothing is lost: a change's full text is in git history up to the commit that archives it, and OpenSpec merges its requirements into `openspec/specs/` at that point, including any it removes. An archived folder still exists on the machine that did the archiving, which is why the scenario check skips the folder rather than relying on git.

### OpenSpec in continuous integration

Continuous integration runs `npx @fission-ai/openspec@2.1.1 validate --strict` with telemetry turned off through `OPENSPEC_TELEMETRY=0`. The version is pinned so a new release can't break the gate without a deliberate update.

### Continuous integration runs on macOS and Linux

`make check` runs on `macos-latest`, because that is where Chorale runs, and on `ubuntu-latest`, because it is faster and catches anything that accidentally depends on macOS. The repo is public, so minutes on both are free.

### Continuous integration setup

The workflow in `.github/workflows/check.yml` is set up as follows:

1. It runs on pushes to any branch and on pull requests, so a branch is checked before it becomes a pull request.
2. Every third-party action is pinned to a full commit hash, with the version in a comment beside it, so a moved tag can't change what runs.
3. The workflow token has `contents: read` and nothing more.
4. Go comes from `go.mod` through `actions/setup-go`, with its module and build cache on. golangci-lint and OpenSpec are installed at exact versions.
5. A newer push to the same branch cancels the older run, and every job has a time limit of 15 minutes.
6. The job names are stable, `check (macos)` and `check (linux)`, because the `harden-repository` change makes them required checks on `main`.

## Test strategy

Each requirement in `specs/app-shell/spec.md` and the layer that covers it:

1. **Empty session list on first start:** interface snapshot tests at 60 by 20 and 120 by 30.
2. **Readable without colour:** an interface snapshot with colour off, plus a unit test that the output contains no colour escape codes.
3. **No motion:** a terminal-level test that reads the screen twice, one second apart.
4. **Quit with q:** an interaction test that `q` ends the program, plus a terminal-level test of the exit status and terminal modes.
5. **Report the version:** a terminal-level test that runs `chorale --version`.
6. **Window too narrow:** interface snapshots at 40 by 20 and at 80 by 15, and an interaction test that sends a resize message.

## Risks / Trade-offs

- **teatest v2 is experimental and may change.** The version is pinned, and only `internal/` test helpers call it directly, so an update touches one place.
- **Golden files can be regenerated carelessly.** `make golden` is the only way to write them, and `AGENTS.md` forbids regenerating them to make a test pass.
- **Snapshot tests on two platforms can differ in small ways, such as line endings or the width of `♪`.** Tests set a fixed environment, including `TERM`, the locale and `NO_COLOR` where needed, and the first failing difference is investigated, never ignored.

## Open Questions

- Whether to adopt Charm's VHS tool for recorded demo tapes. It isn't a gate, so it can be decided later.
